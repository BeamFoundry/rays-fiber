package internal_test

import (
	"context"
	"reflect"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	_ "github.com/BeamFoundry/rays-fiber/pkg/storage"
	. "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/mysql"

	. "github.com/gofiber/storage"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/testing"

	"github.com/gofiber/storage/mysql/v2"

	"github.com/testcontainers/testcontainers-go"
	container "github.com/testcontainers/testcontainers-go/modules/mysql"
)

type TestMySQLConfig struct {
	Component
}

// MySQLConfigOverride implements [mysql.MySQLConfigOverride].
func (t *TestMySQLConfig) MySQLConfigOverride(_ string, _ *mysql.Config) {}

var _ MySQLConfigOverride = &TestMySQLConfig{}

func TestMySQL(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rays-Fiber/Storage/MySQL Suite")
}

var _ = Describe("MySQL", func() {
	It("can register a mysql Storage", func() {

		ctx := context.Background()
		mysqlContainer, err := container.Run(ctx,
			"mysql:latest",
			testcontainers.WithLogger(GetLogger()),
			container.WithDatabase("fiber"),
			container.WithUsername("root"),
			container.WithPassword("password"),
		)

		Expect(err).NotTo(HaveOccurred())

		ip, _ := mysqlContainer.Host(ctx)
		port, _ := mysqlContainer.MappedPort(ctx, "3306")

		config := map[string]any{
			"primary": "mysql",
			"mysql": map[string]any{
				"primary": map[string]any{
					"host":     ip,
					"port":     port.Port(),
					"username": "root",
					"password": "password",
				},
			},
		}

		GetContainer().LoadConfigAt(config, "rays.fiber.storage")
		AddBeamFor(&TestMySQLConfig{})
		GetContainer().EnableScanning()

		err = GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())

		storage := GetBeamScopedInstance[Storage]()
		Expect(reflect.TypeOf(storage).Elem()).To(Equal(reflect.TypeFor[mysql.Storage]()))

		provider := GetBeamScopedInstance[StorageProvider]()
		Expect(provider.Provide("mysql")).NotTo(BeNil())
		Expect(provider.Provide("mysql:primary")).NotTo(BeNil())
		Expect(provider.Provide("mysql:default")).NotTo(BeNil())

		testcontainers.CleanupContainer(GinkgoTB(), mysqlContainer)
		Expect(err).NotTo(HaveOccurred())
	})
})
