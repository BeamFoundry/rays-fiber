package internal_test

import (
	"context"
	"crypto/tls"
	"reflect"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	_ "github.com/BeamFoundry/rays-fiber"
	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/rueidis"

	. "github.com/gofiber/storage"

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

type TestRueidisConfig struct {
	Component
	capturedTLS *tls.Config
}

// RueidisConfigOverride implements [rueidis.RueidisConfigOverride].
func (t *TestRueidisConfig) RueidisConfigOverride(_ string, config *rueidis.Config) {
	t.capturedTLS = config.TLSConfig
	config.TLSConfig = nil
}

var _ RueidisConfigOverride = &TestRueidisConfig{}

type TestRueidisTLSProvider struct {
	Component
}

func (t *TestRueidisTLSProvider) ProvideClientTLS() (*tls.Config, Error) {
	return &tls.Config{ServerName: "rueidis-tls"}, nil
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
}

func (t *TestRueidisConfigWithTLS) ProvideClientTLS() (*tls.Config, Error) {
	return &tls.Config{ServerName: "rueidis-from-override"}, nil
}

func (t *TestRueidisConfigWithTLS) RueidisClientTLSProvider(instance string) certs.ClientTLSProvider {
	if instance == "primary" {
		return t
	}
	return nil
}

func (t *TestRueidisConfigWithTLS) RueidisConfigOverride(_ string, config *rueidis.Config) {
	t.capturedTLS = config.TLSConfig
	config.TLSConfig = nil
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
		config := map[string]any{
			"primary": "rueidis",
			"rueidis": map[string]any{
				"primary": map[string]any{
					"initAddress": []any{address},
				},
			},
		}

		Expect(GetContainer().LoadConfigAt(config, "rays.fiber.storage")).To(Succeed())
		cfgBeam := &TestRueidisConfig{}
		AddBeamFor(cfgBeam)
		AddBeamFor(&TestRueidisTLSProvider{})
		GetContainer().EnableScanning()

		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())

		provider := GetBeamScopedInstance[StorageProvider]()
		st := provider.Provide("rueidis")
		Expect(st).NotTo(BeNil())
		Expect(cfgBeam.capturedTLS).NotTo(BeNil())
		Expect(cfgBeam.capturedTLS.ServerName).To(Equal("rueidis-tls"))
	})

	It("can configure TLS via RueidisConfigOverride implementing ClientTLSProvider", func() {
		config := map[string]any{
			"primary": "rueidis",
			"rueidis": map[string]any{
				"primary": map[string]any{
					"initAddress": []any{address},
				},
			},
		}

		Expect(GetContainer().LoadConfigAt(config, "rays.fiber.storage")).To(Succeed())
		cfgBeam := &TestRueidisConfigWithTLS{}
		AddBeamFor(cfgBeam)
		GetContainer().EnableScanning()

		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())

		provider := GetBeamScopedInstance[StorageProvider]()
		st := provider.Provide("rueidis")
		Expect(st).NotTo(BeNil())
		Expect(cfgBeam.capturedTLS).NotTo(BeNil())
		Expect(cfgBeam.capturedTLS.ServerName).To(Equal("rueidis-from-override"))
	})
})
