package internal_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	_ "github.com/BeamFoundry/rays-fiber"
	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/mongodb"

	. "github.com/gofiber/storage"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/testing"

	"github.com/gofiber/storage/mongodb/v2"

	"github.com/testcontainers/testcontainers-go"
	container "github.com/testcontainers/testcontainers-go/modules/mongodb"
)

func TestMongoDB(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rays-Fiber/Storage/MongoDB Suite")
}

type TestMongoDBConfig struct {
	Component
}

// MongodbConfigOverride implements [mongodb.MongoDBConfigOverride].
func (t *TestMongoDBConfig) MongoDBConfigOverride(_ string, _ *mongodb.Config) {}

var _ MongoDBConfigOverride = &TestMongoDBConfig{}

var _ = Describe("MongoDB", func() {
	It("can register a mongodb Storage", func() {

		ctx := context.Background()
		mongodbContainer, err := container.Run(ctx, "mongo:latest", testcontainers.WithLogger(GetLogger()))

		Expect(err).NotTo(HaveOccurred())

		ip, _ := mongodbContainer.Host(ctx)
		port, _ := mongodbContainer.MappedPort(ctx, "27017")

		fmt.Fprintf(GinkgoWriter, "Mongo Container host/port : %s:%d", ip, port.Num())

		config := map[string]any{
			"primary": "mongodb",
			"mongodb": map[string]any{
				"primary": map[string]any{
					"host": ip,
					"port": port.Port(),
				},
			},
		}

		GetContainer().LoadConfigAt(config, "rays.fiber.storage")
		AddBeamFor(&TestMongoDBConfig{})
		GetContainer().EnableScanning()

		err = GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())

		storage := GetBeamScopedInstance[Storage]()
		Expect(reflect.TypeOf(storage).Elem()).To(Equal(reflect.TypeFor[mongodb.Storage]()))

		provider := GetBeamScopedInstance[StorageProvider]()
		Expect(provider.Provide("mongodb")).NotTo(BeNil())
		Expect(provider.Provide("mongodb:primary")).NotTo(BeNil())
		Expect(provider.Provide("mongodb:default")).NotTo(BeNil())

		testcontainers.CleanupContainer(GinkgoTB(), mongodbContainer)
		Expect(err).NotTo(HaveOccurred())
	})
})
