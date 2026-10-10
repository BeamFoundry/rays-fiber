//nolint:unchecked,errcheck
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
	. "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/memcache"

	. "github.com/gofiber/storage"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/testing"

	"github.com/gofiber/storage/memcache/v2"

	"github.com/testcontainers/testcontainers-go"
	container "github.com/testcontainers/testcontainers-go/modules/memcached"
)

func TestMemcache(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rays-Fiber/Storage/Memcache Suite")
}

type TestMemcacheConfig struct {
	Component
}

// MemcacheConfigOverride implements [memcache.MemcacheConfigOverride].
func (t *TestMemcacheConfig) MemcacheConfigOverride(_ string, _ *memcache.Config) {}

var _ MemcacheConfigOverride = &TestMemcacheConfig{}

var _ = Describe("Memcache", func() {
	It("can register a memcache Storage", func() {

		ctx := context.Background()
		memcachedContainer, err := container.Run(ctx, "docker.io/library/memcached:latest", testcontainers.WithLogger(GetLogger()))

		Expect(err).NotTo(HaveOccurred())

		ip, _ := memcachedContainer.Host(ctx)
		port, _ := memcachedContainer.MappedPort(ctx, "11211")

		config := map[string]any{
			"primary": "memcache",
			"memcache": map[string]any{
				"primary": map[string]any{
					"servers": fmt.Sprintf("%s:%s", ip, port.Port()),
				},
			},
		}

		GetContainer().LoadConfigAt(config, "rays.fiber.storage")
		AddBeamFor(&TestMemcacheConfig{})
		GetContainer().EnableScanning()

		err = GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())

		storage := GetBeamScopedInstance[Storage]()
		Expect(reflect.TypeOf(storage).Elem()).To(Equal(reflect.TypeFor[memcache.Storage]()))

		provider := GetBeamScopedInstance[StorageProvider]()
		Expect(provider.Provide("memcache")).NotTo(BeNil())
		Expect(provider.Provide("memcache:primary")).NotTo(BeNil())
		Expect(provider.Provide("memcache:default")).NotTo(BeNil())

		testcontainers.CleanupContainer(GinkgoTB(), memcachedContainer)
		Expect(err).NotTo(HaveOccurred())
	})
})
