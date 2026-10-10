package internal_test

import (
	"context"
	"reflect"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/testing"

	_ "github.com/BeamFoundry/rays-fiber"
	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/postgres"

	. "github.com/gofiber/storage"

	"github.com/gofiber/storage/postgres/v3"

	"github.com/testcontainers/testcontainers-go"
	container "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestPostgres(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rays-Fiber/Storage/Postgres Suite")
}

type TestPostgresConfig struct {
	Component
}

// PostgresConfigOverride implements [postgres.PostgresConfigOverride].
func (t *TestPostgresConfig) PostgresConfigOverride(_ string, _ *postgres.Config) {}

var _ PostgresConfigOverride = &TestPostgresConfig{}

var _ = Describe("Postgres", func() {
	It("can register a postgres Storage", func() {

		ctx := context.Background()
		postgresContainer, err := container.Run(ctx, "postgres:16-alpine",
			testcontainers.WithLogger(GetLogger()),
			container.WithDatabase("fiber"),
			container.WithUsername("fiber"),
			container.WithPassword("fiber"),
			container.BasicWaitStrategies(),
		)

		Expect(err).NotTo(HaveOccurred())

		ip, _ := postgresContainer.Host(ctx)
		port, _ := postgresContainer.MappedPort(ctx, "5432")

		config := map[string]any{
			"primary": "postgres",
			"postgres": map[string]any{
				"primary": map[string]any{
					"host":     ip,
					"port":     port.Port(),
					"username": "fiber",
					"password": "fiber",
				},
			},
		}

		GetContainer().LoadConfigAt(config, "rays.fiber.storage")
		AddBeamFor(&TestPostgresConfig{})
		GetContainer().EnableScanning()

		err = GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())

		storage := GetBeamScopedInstance[Storage]()
		Expect(reflect.TypeOf(storage).Elem()).To(Equal(reflect.TypeFor[postgres.Storage]()))

		provider := GetBeamScopedInstance[StorageProvider]()
		Expect(provider.Provide("postgres")).NotTo(BeNil())
		Expect(provider.Provide("postgres:primary")).NotTo(BeNil())
		Expect(provider.Provide("postgres:default")).NotTo(BeNil())

		testcontainers.CleanupContainer(GinkgoTB(), postgresContainer)
		Expect(err).NotTo(HaveOccurred())
	})
})
