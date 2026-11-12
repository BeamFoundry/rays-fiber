package internal_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"reflect"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	_ "github.com/BeamFoundry/rays-fiber"
	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/valkey"

	. "github.com/gofiber/storage"

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

type TestValkeyConfig struct {
	Component
	capturedTLS *tls.Config
}

// ValkeyConfigOverride implements [valkey.ValkeyConfigOverride].
func (t *TestValkeyConfig) ValkeyConfigOverride(_ string, config *valkey.Config) {
	t.capturedTLS = config.TLSConfig
	config.TLSConfig = nil
}

var _ ValkeyConfigOverride = &TestValkeyConfig{}

type TestValkeyTLSProvider struct {
	Component
}

func (t *TestValkeyTLSProvider) ProvideClientTLS() (*tls.Config, Error) {
	return &tls.Config{ServerName: "valkey-tls"}, nil
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
}

func (t *TestValkeyConfigWithTLS) ProvideClientTLS() (*tls.Config, Error) {
	return &tls.Config{ServerName: "valkey-from-override"}, nil
}

func (t *TestValkeyConfigWithTLS) ValkeyClientTLSProvider(instance string) certs.ClientTLSProvider {
	if instance == "primary" {
		return t
	}
	return nil
}

func (t *TestValkeyConfigWithTLS) ValkeyConfigOverride(_ string, config *valkey.Config) {
	t.capturedTLS = config.TLSConfig
	config.TLSConfig = nil
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
		config := map[string]any{
			"primary": "valkey",
			"valkey": map[string]any{
				"primary": map[string]any{
					"initAddress": []any{initAddress},
				},
			},
		}

		Expect(GetContainer().LoadConfigAt(config, "rays.fiber.storage")).To(Succeed())
		cfgBeam := &TestValkeyConfig{}
		AddBeamFor(cfgBeam)
		AddBeamFor(&TestValkeyTLSProvider{})
		GetContainer().EnableScanning()

		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())

		provider := GetBeamScopedInstance[StorageProvider]()
		st := provider.Provide("valkey")
		Expect(st).NotTo(BeNil())
		Expect(cfgBeam.capturedTLS).NotTo(BeNil())
		Expect(cfgBeam.capturedTLS.ServerName).To(Equal("valkey-tls"))
	})

	It("can configure TLS via ValkeyConfigOverride implementing ClientTLSProvider", func() {
		config := map[string]any{
			"primary": "valkey",
			"valkey": map[string]any{
				"primary": map[string]any{
					"initAddress": []any{initAddress},
				},
			},
		}

		Expect(GetContainer().LoadConfigAt(config, "rays.fiber.storage")).To(Succeed())
		cfgBeam := &TestValkeyConfigWithTLS{}
		AddBeamFor(cfgBeam)
		GetContainer().EnableScanning()

		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())

		provider := GetBeamScopedInstance[StorageProvider]()
		st := provider.Provide("valkey")
		Expect(st).NotTo(BeNil())
		Expect(cfgBeam.capturedTLS).NotTo(BeNil())
		Expect(cfgBeam.capturedTLS.ServerName).To(Equal("valkey-from-override"))
	})
})
