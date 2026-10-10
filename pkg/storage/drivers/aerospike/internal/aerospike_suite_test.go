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
	. "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/aerospike"

	. "github.com/gofiber/storage"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/testing"

	"github.com/gofiber/storage/aerospike"

	"github.com/testcontainers/testcontainers-go"
	container "github.com/testcontainers/testcontainers-go/modules/aerospike"
)

func TestAerospike(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rays-Fiber/Storage/Aerospike Suite")
}

type TestAerospikeConfig struct {
	Component
}

// AerospikeConfigOverride implements [aerospike.AerospikeConfigOverride].
func (t *TestAerospikeConfig) AerospikeConfigOverride(_ string, _ *aerospike.Config) {}

var _ AerospikeConfigOverride = &TestAerospikeConfig{}

var _ = Describe("Aerospike", func() {
	It("can register a aerospike Storage", func() {

		ctx := context.Background()
		aerospikeContainer, err := container.Run(ctx, "aerospike/aerospike-server:latest", testcontainers.WithLogger(GetLogger()))

		Expect(err).NotTo(HaveOccurred())

		ip, _ := aerospikeContainer.Host(ctx)
		port, _ := aerospikeContainer.MappedPort(ctx, "3000")

		config := map[string]any{
			"primary": "aerospike",
			"aerospike": map[string]any{
				"primary": map[string]any{
					"hosts": []any{fmt.Sprintf("%s:%s", ip, port.Port())},
				},
			},
		}

		GetContainer().LoadConfigAt(config, "rays.fiber.storage")
		AddBeamFor(&TestAerospikeConfig{})
		GetContainer().EnableScanning()

		err = GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())

		storage := GetBeamScopedInstance[Storage]()
		Expect(reflect.TypeOf(storage).Elem()).To(Equal(reflect.TypeFor[aerospike.Storage]()))

		provider := GetBeamScopedInstance[StorageProvider]()
		Expect(provider.Provide("aerospike")).NotTo(BeNil())
		Expect(provider.Provide("aerospike:primary")).NotTo(BeNil())
		Expect(provider.Provide("aerospike:default")).NotTo(BeNil())

		testcontainers.CleanupContainer(GinkgoTB(), aerospikeContainer)
		Expect(err).NotTo(HaveOccurred())
	})
})
