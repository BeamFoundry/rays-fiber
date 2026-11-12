package internal_test

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"net"
	"reflect"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	_ "github.com/BeamFoundry/rays-fiber"
	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/valkey"

	. "github.com/gofiber/storage"
	"github.com/mdelapenya/tlscert"

	"github.com/BeamFoundry/rays-web-security/pkg/certs"
	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"
	. "github.com/BeamFoundry/rays/pkg/testing"

	"github.com/gofiber/storage/valkey"

	"github.com/testcontainers/testcontainers-go"
	container "github.com/testcontainers/testcontainers-go/modules/valkey"
)

func TestValkey(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rays-Fiber/Storage/Valkey Suite")
}

func startTLSTestServer(targetAddr string) (tlsAddr string, certPool *x509.CertPool, closeFn func()) {
	cert := tlscert.SelfSigned("127.0.0.1,valkey-tls,valkey-from-override,localhost")
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

type TestValkeyConfig struct {
	Component
	capturedTLS *tls.Config
}

// ValkeyConfigOverride implements [valkey.ValkeyConfigOverride].
func (t *TestValkeyConfig) ValkeyConfigOverride(_ string, config *valkey.Config) {
	t.capturedTLS = config.TLSConfig
}

var _ ValkeyConfigOverride = &TestValkeyConfig{}

type TestValkeyTLSProvider struct {
	Component
	certPool *x509.CertPool
}

func (t *TestValkeyTLSProvider) ProvideClientTLS() (*tls.Config, Error) {
	return &tls.Config{
		ServerName: "valkey-tls",
		RootCAs:    t.certPool,
	}, nil
}

func (t *TestValkeyTLSProvider) ValkeyClientTLSProvider(instance string) certs.ClientTLSProvider {
	if instance == "primary" {
		return t
	}
	return nil
}

var _ ValkeyClientTLSProvider = &TestValkeyTLSProvider{}

type TestValkeyConfigWithTLS struct {
	Component
	capturedTLS *tls.Config
	certPool    *x509.CertPool
}

func (t *TestValkeyConfigWithTLS) ProvideClientTLS() (*tls.Config, Error) {
	return &tls.Config{
		ServerName: "valkey-from-override",
		RootCAs:    t.certPool,
	}, nil
}

func (t *TestValkeyConfigWithTLS) ValkeyClientTLSProvider(instance string) certs.ClientTLSProvider {
	if instance == "primary" {
		return t
	}
	return nil
}

func (t *TestValkeyConfigWithTLS) ValkeyConfigOverride(_ string, config *valkey.Config) {
	t.capturedTLS = config.TLSConfig
}

var (
	_ ValkeyConfigOverride    = &TestValkeyConfigWithTLS{}
	_ ValkeyClientTLSProvider = &TestValkeyConfigWithTLS{}
)

var _ = Describe("Valkey", func() {
	var (
		ctx             context.Context
		valkeyContainer testcontainers.Container
		initAddress     string
	)

	BeforeEach(func() {
		ctx = context.Background()
		var err error
		valkeyContainer, err = container.Run(ctx, "valkey/valkey:7.2.5", testcontainers.WithLogger(GetLogger()))
		Expect(err).NotTo(HaveOccurred())

		ip, err := valkeyContainer.Host(ctx)
		Expect(err).NotTo(HaveOccurred())
		port, err := valkeyContainer.MappedPort(ctx, "6379")
		Expect(err).NotTo(HaveOccurred())
		initAddress = fmt.Sprintf("%s:%s", ip, port.Port())
	})

	AfterEach(func() {
		if valkeyContainer != nil {
			testcontainers.CleanupContainer(GinkgoTB(), valkeyContainer)
		}
	})

	It("can register a valkey Storage", func() {
		config := map[string]any{
			"primary": "valkey",
			"valkey": map[string]any{
				"primary": map[string]any{
					"initAddress": []any{initAddress},
				},
			},
		}

		fmt.Fprintf(GinkgoWriter, "ValKey Container initAddress '%s'\n", initAddress)

		Expect(GetContainer().LoadConfigAt(config, "rays.fiber.storage")).To(Succeed())
		AddBeamFor(&TestValkeyConfig{})
		GetContainer().EnableScanning()

		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())

		storage := GetBeamScopedInstance[Storage]()
		Expect(reflect.TypeOf(storage).Elem()).To(Equal(reflect.TypeFor[valkey.Storage]()))

		provider := GetBeamScopedInstance[StorageProvider]()
		Expect(provider.Provide("valkey")).NotTo(BeNil())
		Expect(provider.Provide("valkey:primary")).NotTo(BeNil())
		Expect(provider.Provide("valkey:default")).NotTo(BeNil())
	})

	It("can configure TLS via ValkeyClientTLSProvider", func() {
		tlsAddr, certPool, stopServer := startTLSTestServer(initAddress)
		defer stopServer()

		config := map[string]any{
			"primary": "valkey",
			"valkey": map[string]any{
				"primary": map[string]any{
					"initAddress": []any{tlsAddr},
				},
			},
		}

		Expect(GetContainer().LoadConfigAt(config, "rays.fiber.storage")).To(Succeed())
		cfgBeam := &TestValkeyConfig{}
		AddBeamFor(cfgBeam)
		AddBeamFor(&TestValkeyTLSProvider{certPool: certPool})
		GetContainer().EnableScanning()

		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())

		provider := GetBeamScopedInstance[StorageProvider]()
		st := provider.Provide("valkey")
		Expect(st).NotTo(BeNil())
		Expect(cfgBeam.capturedTLS).NotTo(BeNil())
		Expect(cfgBeam.capturedTLS.ServerName).To(Equal("valkey-tls"))

		// Perform a storage operation over the TLS connection to verify the handshake
		Expect(st.Set("tls-key", []byte("tls-val"), time.Minute)).To(Succeed())
		val, getErr := st.Get("tls-key")
		Expect(getErr).NotTo(HaveOccurred())
		Expect(val).To(Equal([]byte("tls-val")))
	})

	It("can configure TLS via ValkeyConfigOverride implementing ClientTLSProvider", func() {
		tlsAddr, certPool, stopServer := startTLSTestServer(initAddress)
		defer stopServer()

		config := map[string]any{
			"primary": "valkey",
			"valkey": map[string]any{
				"primary": map[string]any{
					"initAddress": []any{tlsAddr},
				},
			},
		}

		Expect(GetContainer().LoadConfigAt(config, "rays.fiber.storage")).To(Succeed())
		cfgBeam := &TestValkeyConfigWithTLS{certPool: certPool}
		AddBeamFor(cfgBeam)
		GetContainer().EnableScanning()

		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())

		provider := GetBeamScopedInstance[StorageProvider]()
		st := provider.Provide("valkey")
		Expect(st).NotTo(BeNil())
		Expect(cfgBeam.capturedTLS).NotTo(BeNil())
		Expect(cfgBeam.capturedTLS.ServerName).To(Equal("valkey-from-override"))

		// Perform a storage operation over the TLS connection to verify the handshake
		Expect(st.Set("override-key", []byte("override-val"), time.Minute)).To(Succeed())
		val, getErr := st.Get("override-key")
		Expect(getErr).NotTo(HaveOccurred())
		Expect(val).To(Equal([]byte("override-val")))
	})
})
