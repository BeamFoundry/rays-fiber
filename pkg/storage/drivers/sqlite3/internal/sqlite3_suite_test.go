package internal3_test

import (
	"reflect"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	_ "github.com/BeamFoundry/rays-fiber"
	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/sqlite3"

	. "github.com/gofiber/storage"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/testing"

	"github.com/gofiber/storage/sqlite3/v2"
)

func TestSQLite3(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rays-Fiber/Storage/SQLite3 Suite")
}

type TestSQLite3Config struct {
	Component
}

// SQLite3ConfigOverride implements [sqlite3.SQLite3ConfigOverride].
func (t *TestSQLite3Config) SQLite3ConfigOverride(_ string, _ *sqlite3.Config) {
}

var _ SQLite3ConfigOverride = &TestSQLite3Config{}

var _ = Describe("SQLite3", func() {
	It("can register a sqlite3 Storage", func() {

		config := map[string]any{
			"primary": "sqlite3",
			"sqlite3": map[string]any{
				"primary": map[string]any{
					"database": "file::memory:",
				},
			},
		}

		GetContainer().LoadConfigAt(config, "rays.fiber.storage")
		AddBeamFor(&TestSQLite3Config{})
		GetContainer().EnableScanning()

		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())

		storage := GetBeamScopedInstance[Storage]()
		Expect(reflect.TypeOf(storage).Elem()).To(Equal(reflect.TypeFor[sqlite3.Storage]()))

		provider := GetBeamScopedInstance[StorageProvider]()
		Expect(provider.Provide("sqlite3")).NotTo(BeNil())
		Expect(provider.Provide("sqlite3:primary")).NotTo(BeNil())
		Expect(provider.Provide("sqlite3:default")).NotTo(BeNil())

	})
})
