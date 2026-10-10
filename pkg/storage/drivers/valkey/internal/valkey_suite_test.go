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
	. "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/valkey"

	. "github.com/gofiber/storage"

	. "github.com/BeamFoundry/rays/pkg/core"
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
}

// ValkeyConfigOverride implements [valkey.ValkeyConfigOverride].
func (t *TestValkeyConfig) ValkeyConfigOverride(_ string, _ *valkey.Config) {}

var _ ValkeyConfigOverride = &TestValkeyConfig{}

var _ = Describe("Valkey", func() {
	It("can register a valkey Storage", func() {

		ctx := context.Background()
		valkeyContainer, err := container.Run(ctx, "valkey/valkey:7.2.5", testcontainers.WithLogger(GetLogger()))

		Expect(err).NotTo(HaveOccurred())

		ip, _ := valkeyContainer.Host(ctx)
		port, _ := valkeyContainer.MappedPort(ctx, "6379")

		initAddress := fmt.Sprintf("%s:%s", ip, port.Port())

		config := map[string]any{
			"primary": "valkey",
			"valkey": map[string]any{
				"primary": map[string]any{
					"initAddress": []any{initAddress},
				},
			},
		}

		fmt.Fprintf(GinkgoWriter, "ValKey Container initAddress '%s'\n", initAddress)

		GetContainer().LoadConfigAt(config, "rays.fiber.storage")
		AddBeamFor(&TestValkeyConfig{})
		GetContainer().EnableScanning()

		err = GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())

		storage := GetBeamScopedInstance[Storage]()
		Expect(reflect.TypeOf(storage).Elem()).To(Equal(reflect.TypeFor[valkey.Storage]()))

		provider := GetBeamScopedInstance[StorageProvider]()
		Expect(provider.Provide("valkey")).NotTo(BeNil())
		Expect(provider.Provide("valkey:primary")).NotTo(BeNil())
		Expect(provider.Provide("valkey:default")).NotTo(BeNil())

		testcontainers.CleanupContainer(GinkgoTB(), valkeyContainer)
		Expect(err).NotTo(HaveOccurred())
	})
})
