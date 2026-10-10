package internal_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	_ "github.com/BeamFoundry/rays-fiber"
	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	_ "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/memory"

	. "github.com/BeamFoundry/rays/pkg/testing"
)

func TestStorageInternal(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rays-Fiber/Storage/Internal Suite")
}

var _ = Describe("Provider", func() {
	It("can get the default storage", func() {
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
