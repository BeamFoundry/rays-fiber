package internal_test

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"net"
	"reflect"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	_ "github.com/BeamFoundry/rays-fiber"
	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/rueidis"

	. "github.com/gofiber/storage"
	"github.com/mdelapenya/tlscert"

	"github.com/BeamFoundry/rays-web-security/pkg/certs"
	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"
	. "github.com/BeamFoundry/rays/pkg/testing"

	"github.com/gofiber/storage/rueidis"

	"github.com/testcontainers/testcontainers-go"
	container "github.com/testcontainers/testcontainers-go/modules/redis"
)

func TestRueidis(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rays-Fiber/Storage/Rueidis Suite")
}

func startTLSTestServer(targetAddr string) (tlsAddr string, certPool *x509.CertPool, closeFn func()) {
	cert := tlscert.SelfSigned("127.0.0.1,rueidis-tls,rueidis-from-override,localhost")
	tlsConfig := cert.TLSConfig()

	certPool = x509.NewCertPool()
	certPool.AddCert(cert.Cert)

	listener, err := tls.Listen("tcp", "127.0.0.1:0", tlsConfig)
	Expect(err).NotTo(HaveOccurred())

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			clientConn, acceptErr := listener.Accept()
			if acceptErr != nil {
				return
			}
			go func(c net.Conn) {
				defer GinkgoRecover()
				defer func() {
					if closeErr := c.Close(); closeErr != nil && !errors.Is(closeErr, net.ErrClosed) {
						Expect(closeErr).NotTo(HaveOccurred())
					}
				}()
				var dialer net.Dialer
				backendConn, dialErr := dialer.DialContext(ctx, "tcp", targetAddr)
				if dialErr != nil {
					return
				}
				defer func() {
					if closeErr := backendConn.Close(); closeErr != nil && !errors.Is(closeErr, net.ErrClosed) {
						Expect(closeErr).NotTo(HaveOccurred())
					}
				}()

				done := make(chan struct{}, 2)
				go func() {
					defer GinkgoRecover()
					defer func() {
						done <- struct{}{}
					}()
					if _, copyErr := io.Copy(backendConn, c); copyErr != nil && !errors.Is(copyErr, net.ErrClosed) {
						Expect(copyErr).NotTo(HaveOccurred())
					}
				}()
				go func() {
					defer GinkgoRecover()
					defer func() {
						done <- struct{}{}
					}()
					if _, copyErr := io.Copy(c, backendConn); copyErr != nil && !errors.Is(copyErr, net.ErrClosed) {
						Expect(copyErr).NotTo(HaveOccurred())
					}
				}()
				select {
				case <-done:
				case <-ctx.Done():
				}
			}(clientConn)
		}
	}()

	return listener.Addr().String(), certPool, func() {
		cancel()
		if closeErr := listener.Close(); closeErr != nil && !errors.Is(closeErr, net.ErrClosed) {
			Expect(closeErr).NotTo(HaveOccurred())
		}
	}
}

type TestRueidisConfig struct {
	Component
	capturedTLS *tls.Config
}

// RueidisConfigOverride implements [rueidis.RueidisConfigOverride].
func (t *TestRueidisConfig) RueidisConfigOverride(_ string, config *rueidis.Config) {
	t.capturedTLS = config.TLSConfig
}

var _ RueidisConfigOverride = &TestRueidisConfig{}

type TestRueidisTLSProvider struct {
	Component
	certPool *x509.CertPool
}

func (t *TestRueidisTLSProvider) ProvideClientTLS() (*tls.Config, Error) {
	return &tls.Config{
		ServerName: "rueidis-tls",
		RootCAs:    t.certPool,
	}, nil
}

func (t *TestRueidisTLSProvider) RueidisClientTLSProvider(instance string) certs.ClientTLSProvider {
	if instance == "primary" {
		return t
	}
	return nil
}

var _ RueidisClientTLSProvider = &TestRueidisTLSProvider{}

type TestRueidisConfigWithTLS struct {
	Component
	capturedTLS *tls.Config
	certPool    *x509.CertPool
}

func (t *TestRueidisConfigWithTLS) ProvideClientTLS() (*tls.Config, Error) {
	return &tls.Config{
		ServerName: "rueidis-from-override",
		RootCAs:    t.certPool,
	}, nil
}

func (t *TestRueidisConfigWithTLS) RueidisClientTLSProvider(instance string) certs.ClientTLSProvider {
	if instance == "primary" {
		return t
	}
	return nil
}

func (t *TestRueidisConfigWithTLS) RueidisConfigOverride(_ string, config *rueidis.Config) {
	t.capturedTLS = config.TLSConfig
}

var (
	_ RueidisConfigOverride    = &TestRueidisConfigWithTLS{}
	_ RueidisClientTLSProvider = &TestRueidisConfigWithTLS{}
)

var _ = Describe("Rueidis", func() {
	var (
		ctx            context.Context
		redisContainer testcontainers.Container
		address        string
	)

	BeforeEach(func() {
		ctx = context.Background()
		var err error
		redisContainer, err = container.Run(ctx, "redis:8.2", testcontainers.WithLogger(GetLogger()))
		Expect(err).NotTo(HaveOccurred())

		ip, err := redisContainer.Host(ctx)
		Expect(err).NotTo(HaveOccurred())
		port, err := redisContainer.MappedPort(ctx, "6379")
		Expect(err).NotTo(HaveOccurred())
		address = ip + ":" + port.Port()
	})

	AfterEach(func() {
		if redisContainer != nil {
			testcontainers.CleanupContainer(GinkgoTB(), redisContainer)
		}
	})

	It("can register a rueidis Storage", func() {
		config := map[string]any{
			"primary": "rueidis",
			"rueidis": map[string]any{
				"primary": map[string]any{
					"initAddress": []any{address},
				},
			},
		}

		Expect(GetContainer().LoadConfigAt(config, "rays.fiber.storage")).To(Succeed())
		AddBeamFor(&TestRueidisConfig{})
		GetContainer().EnableScanning()

		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())

		storage := GetBeamScopedInstance[Storage]()
		Expect(reflect.TypeOf(storage).Elem()).To(Equal(reflect.TypeFor[rueidis.Storage]()))

		provider := GetBeamScopedInstance[StorageProvider]()
		Expect(provider.Provide("rueidis")).NotTo(BeNil())
		Expect(provider.Provide("rueidis:primary")).NotTo(BeNil())
		Expect(provider.Provide("rueidis:default")).NotTo(BeNil())
	})

	It("can configure TLS via RueidisClientTLSProvider", func() {
		tlsAddr, certPool, stopServer := startTLSTestServer(address)
		defer stopServer()

		config := map[string]any{
			"primary": "rueidis",
			"rueidis": map[string]any{
				"primary": map[string]any{
					"initAddress": []any{tlsAddr},
				},
			},
		}

		Expect(GetContainer().LoadConfigAt(config, "rays.fiber.storage")).To(Succeed())
		cfgBeam := &TestRueidisConfig{}
		AddBeamFor(cfgBeam)
		AddBeamFor(&TestRueidisTLSProvider{certPool: certPool})
		GetContainer().EnableScanning()

		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())

		provider := GetBeamScopedInstance[StorageProvider]()
		st := provider.Provide("rueidis")
		Expect(st).NotTo(BeNil())
		Expect(cfgBeam.capturedTLS).NotTo(BeNil())
		Expect(cfgBeam.capturedTLS.ServerName).To(Equal("rueidis-tls"))

		// Perform a storage operation over the TLS connection to verify the handshake
		Expect(st.Set("rueidis-tls-key", []byte("rueidis-tls-val"), time.Minute)).To(Succeed())
		val, getErr := st.Get("rueidis-tls-key")
		Expect(getErr).NotTo(HaveOccurred())
		Expect(val).To(Equal([]byte("rueidis-tls-val")))
	})

	It("can configure TLS via RueidisConfigOverride implementing ClientTLSProvider", func() {
		tlsAddr, certPool, stopServer := startTLSTestServer(address)
		defer stopServer()

		config := map[string]any{
			"primary": "rueidis",
			"rueidis": map[string]any{
				"primary": map[string]any{
					"initAddress": []any{tlsAddr},
				},
			},
		}

		Expect(GetContainer().LoadConfigAt(config, "rays.fiber.storage")).To(Succeed())
		cfgBeam := &TestRueidisConfigWithTLS{certPool: certPool}
		AddBeamFor(cfgBeam)
		GetContainer().EnableScanning()

		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())

		provider := GetBeamScopedInstance[StorageProvider]()
		st := provider.Provide("rueidis")
		Expect(st).NotTo(BeNil())
		Expect(cfgBeam.capturedTLS).NotTo(BeNil())
		Expect(cfgBeam.capturedTLS.ServerName).To(Equal("rueidis-from-override"))

		// Perform a storage operation over the TLS connection to verify the handshake
		Expect(st.Set("rueidis-override-key", []byte("rueidis-override-val"), time.Minute)).To(Succeed())
		val, getErr := st.Get("rueidis-override-key")
		Expect(getErr).NotTo(HaveOccurred())
		Expect(val).To(Equal([]byte("rueidis-override-val")))
	})
})
