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
	. "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/redis"

	. "github.com/gofiber/storage"

	"github.com/BeamFoundry/rays-web-security/pkg/certs"
	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"
	. "github.com/BeamFoundry/rays/pkg/testing"

	"github.com/gofiber/storage/redis/v3"

	"github.com/testcontainers/testcontainers-go"
	container "github.com/testcontainers/testcontainers-go/modules/redis"
)

func TestRedis(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rays-Fiber/Storage/Redis Suite")
}

type TestRedisConfig struct {
	Component
	capturedTLS *tls.Config
}

var _ RedisConfigOverride = &TestRedisConfig{}

// RedisConfig implements [redis.RedisConfigOverride].
func (t *TestRedisConfig) RedisConfigOverride(_ string, config *redis.Config) {
	Expect(config).NotTo(BeNil())
	t.capturedTLS = config.TLSConfig
}

type TestRedisTLSProvider struct {
	Component
}

func (t *TestRedisTLSProvider) ProvideClientTLS() (*tls.Config, Error) {
	return &tls.Config{ServerName: "redis-tls"}, nil
}

func (t *TestRedisTLSProvider) RedisClientTLSProvider(instance string) certs.ClientTLSProvider {
	if instance == "primary" {
		return t
	}
	return nil
}

var _ RedisClientTLSProvider = &TestRedisTLSProvider{}

type TestRedisConfigWithTLS struct {
	Component
	capturedTLS *tls.Config
}

func (t *TestRedisConfigWithTLS) ProvideClientTLS() (*tls.Config, Error) {
	return &tls.Config{ServerName: "redis-from-override"}, nil
}

func (t *TestRedisConfigWithTLS) RedisClientTLSProvider(instance string) certs.ClientTLSProvider {
	if instance == "primary" {
		return t
	}
	return nil
}

func (t *TestRedisConfigWithTLS) RedisConfigOverride(_ string, config *redis.Config) {
	t.capturedTLS = config.TLSConfig
}

var (
	_ RedisConfigOverride    = &TestRedisConfigWithTLS{}
	_ RedisClientTLSProvider = &TestRedisConfigWithTLS{}
)

var _ = Describe("Redis", func() {
	It("can register a redis Storage", func() {
		ctx := context.Background()

		redisContainer, err := container.Run(ctx, "redis:8.2",
			testcontainers.WithLogger(GetLogger()),
		)

		Expect(err).NotTo(HaveOccurred())

		ip, err := redisContainer.Host(ctx)
		Expect(err).NotTo(HaveOccurred())
		port, err := redisContainer.MappedPort(ctx, "6379")
		Expect(err).NotTo(HaveOccurred())

		config := map[string]any{
			"primary": "redis",
			"redis": map[string]any{
				"primary": map[string]any{
					"host": ip,
					"port": port.Port(),
				},
			},
		}

		Expect(GetContainer().LoadConfigAt(config, "rays.fiber.storage")).To(Succeed())
		AddBeamFor(&TestRedisConfig{})
		GetContainer().EnableScanning()

		err = GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())

		storage := GetBeamScopedInstance[Storage]()
		Expect(reflect.TypeOf(storage).Elem()).To(Equal(reflect.TypeFor[redis.Storage]()))

		provider := GetBeamScopedInstance[StorageProvider]()
		Expect(provider.Provide("redis")).NotTo(BeNil())
		Expect(provider.Provide("redis:primary")).NotTo(BeNil())
		Expect(provider.Provide("redis:default")).NotTo(BeNil())

		testcontainers.CleanupContainer(GinkgoTB(), redisContainer)
		Expect(err).NotTo(HaveOccurred())
	})

	It("can configure TLS via RedisClientTLSProvider", func() {
		config := map[string]any{
			"primary": "redis",
			"redis": map[string]any{
				"primary": map[string]any{
					"skipConnectionCheck": true,
				},
			},
		}

		Expect(GetContainer().LoadConfigAt(config, "rays.fiber.storage")).To(Succeed())
		cfgBeam := &TestRedisConfig{}
		AddBeamFor(cfgBeam)
		AddBeamFor(&TestRedisTLSProvider{})
		GetContainer().EnableScanning()

		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())

		provider := GetBeamScopedInstance[StorageProvider]()
		st := provider.Provide("redis")
		Expect(st).NotTo(BeNil())
		Expect(cfgBeam.capturedTLS).NotTo(BeNil())
		Expect(cfgBeam.capturedTLS.ServerName).To(Equal("redis-tls"))
	})

	It("can configure TLS via RedisConfigOverride implementing ClientTLSProvider", func() {
		config := map[string]any{
			"primary": "redis",
			"redis": map[string]any{
				"primary": map[string]any{
					"skipConnectionCheck": true,
				},
			},
		}

		Expect(GetContainer().LoadConfigAt(config, "rays.fiber.storage")).To(Succeed())
		cfgBeam := &TestRedisConfigWithTLS{}
		AddBeamFor(cfgBeam)
		GetContainer().EnableScanning()

		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())

		provider := GetBeamScopedInstance[StorageProvider]()
		st := provider.Provide("redis")
		Expect(st).NotTo(BeNil())
		Expect(cfgBeam.capturedTLS).NotTo(BeNil())
		Expect(cfgBeam.capturedTLS.ServerName).To(Equal("redis-from-override"))
	})
})
