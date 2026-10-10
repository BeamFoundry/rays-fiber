package internal_test

import (
	"context"
	"reflect"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	_ "github.com/BeamFoundry/rays-fiber"
	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/rueidis"

	. "github.com/gofiber/storage"

	. "github.com/BeamFoundry/rays/pkg/core"
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
}

// RueidisConfigOverride implements [rueidis.RueidisConfigOverride].
func (t *TestRueidisConfig) RueidisConfigOverride(_ string, _ *rueidis.Config) {}

var _ RueidisConfigOverride = &TestRueidisConfig{}

var _ = Describe("Rueidis", func() {
	It("can register a rueidis Storage", func() {

		ctx := context.Background()

		redisContainer, err := container.Run(ctx, "redis:8.2", testcontainers.WithLogger(GetLogger()))

		Expect(err).NotTo(HaveOccurred())

		ip, _ := redisContainer.Host(ctx)
		port, _ := redisContainer.MappedPort(ctx, "6379")

		config := map[string]any{
			"primary": "rueidis",
			"rueidis": map[string]any{
				"primary": map[string]any{
					"initAddress": []any{ip + ":" + port.Port()},
				},
			},
		}

		GetContainer().LoadConfigAt(config, "rays.fiber.storage")
		AddBeamFor(&TestRueidisConfig{})
		GetContainer().EnableScanning()

		err = GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())

		storage := GetBeamScopedInstance[Storage]()
		Expect(reflect.TypeOf(storage).Elem()).To(Equal(reflect.TypeFor[rueidis.Storage]()))

		provider := GetBeamScopedInstance[StorageProvider]()
		Expect(provider.Provide("rueidis")).NotTo(BeNil())
		Expect(provider.Provide("rueidis:primary")).NotTo(BeNil())
		Expect(provider.Provide("rueidis:default")).NotTo(BeNil())

		testcontainers.CleanupContainer(GinkgoTB(), redisContainer)
		Expect(err).NotTo(HaveOccurred())
	})
})
