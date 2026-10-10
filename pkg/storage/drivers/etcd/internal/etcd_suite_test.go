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
	. "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/etcd"

	. "github.com/gofiber/storage"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/testing"

	"github.com/gofiber/storage/etcd/v2"

	"github.com/testcontainers/testcontainers-go"
	container "github.com/testcontainers/testcontainers-go/modules/etcd"
)

func TestEtcd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rays-Fiber/Storage/Etcd Suite")
}

type TestEtcdConfig struct {
	Component
}

// EtcdConfigOverride implements [etcd.EtcdConfigOverride].
func (t *TestEtcdConfig) EtcdConfigOverride(_ string, _ *etcd.Config) {
}

var _ EtcdConfigOverride = &TestEtcdConfig{}

var _ = Describe("Etcd", func() {
	It("can register a etcd Storage", func() {

		ctx := context.Background()
		etcdContainer, err := container.Run(ctx, "gcr.io/etcd-development/etcd:v3.5.14", testcontainers.WithLogger(GetLogger()))

		Expect(err).NotTo(HaveOccurred())

		ip, _ := etcdContainer.Host(ctx)
		port, _ := etcdContainer.MappedPort(ctx, "2379")

		endpoints := fmt.Sprintf("%s:%s", ip, port.Port())
		config := map[string]any{
			"primary": "etcd",
			"etcd": map[string]any{
				"primary": map[string]any{
					"endpoints": []any{endpoints},
				},
			},
		}

		GetContainer().LoadConfigAt(config, "rays.fiber.storage")
		AddBeamFor(&TestEtcdConfig{})
		GetContainer().EnableScanning()

		err = GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())

		storage := GetBeamScopedInstance[Storage]()
		Expect(reflect.TypeOf(storage).Elem()).To(Equal(reflect.TypeFor[etcd.Storage]()))

		provider := GetBeamScopedInstance[StorageProvider]()
		Expect(provider.Provide("etcd")).NotTo(BeNil())
		Expect(provider.Provide("etcd:primary")).NotTo(BeNil())
		Expect(provider.Provide("etcd:default")).NotTo(BeNil())

		testcontainers.CleanupContainer(GinkgoTB(), etcdContainer)
		Expect(err).NotTo(HaveOccurred())
	})
})
