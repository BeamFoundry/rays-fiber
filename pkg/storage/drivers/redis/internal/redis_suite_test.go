package internal_test

import (
	"context"
	"reflect"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	_ "github.com/BeamFoundry/rays-fiber"
	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/redis"

	. "github.com/gofiber/storage"

	. "github.com/BeamFoundry/rays/pkg/core"
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
}

var _ RedisConfigOverride = &TestRedisConfig{}

// RedisConfig implements [redis.RedisConfigOverride].
func (t *TestRedisConfig) RedisConfigOverride(name string, config *redis.Config) {
	Expect(config).NotTo(BeNil())
}

var _ = Describe("Redis", func() {
	It("can register a redis Storage", func() {

		ctx := context.Background()

		redisContainer, err := container.Run(ctx, "redis:8.2",
			testcontainers.WithLogger(GetLogger()),
		)

		Expect(err).NotTo(HaveOccurred())

		ip, _ := redisContainer.Host(ctx)
		port, _ := redisContainer.MappedPort(ctx, "6379")

		config := map[string]any{
			"primary": "redis",
			"redis": map[string]any{
				"primary": map[string]any{
					"host": ip,
					"port": port.Port(),
				},
			},
		}

		GetContainer().LoadConfigAt(config, "rays.fiber.storage")
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
})
