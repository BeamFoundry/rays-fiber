package memory_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/BeamFoundry/rays/pkg/core"

	_ "github.com/BeamFoundry/rays-fiber"
	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/memory"

	. "github.com/BeamFoundry/rays/pkg/testing"

	"github.com/gofiber/storage/memory/v2"
)

func TestStorageMemory(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rays-Fiber/Storage/Memory Suite")
}

type TestMemoryConfig struct {
	Component
}

var _ MemoryConfigOverride = &TestMemoryConfig{}

// MemoryConfig implements [MemoryConfigOverride].
func (t *TestMemoryConfig) MemoryConfigOverride(name string, config *memory.Config) {
	Expect(name).NotTo(BeEmpty())
	Expect(config).NotTo(BeNil())
}

var _ = Describe("Provider", func() {
	It("can get the default storage", func() {
		AddBeamFor(&TestMemoryConfig{})
		GetContainer().EnableScanning()
		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())

		provider := GetBeamScopedInstance[StorageProvider]()

		Expect(provider.Provide("default")).NotTo(BeNil())
		Expect(provider.Provide("memory:primary")).NotTo(BeNil())
		Expect(provider.Provide("memory")).NotTo(BeNil())
		Expect(provider.Provide("undefined")).NotTo(BeNil())
		Expect(provider.Provide("memory:undefined")).NotTo(BeNil())
	})
})
