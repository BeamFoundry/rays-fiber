package internal_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"

	_ "github.com/BeamFoundry/rays-fiber"
	. "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/couchbase"

	. "github.com/gofiber/storage"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/testing"

	"github.com/gofiber/storage/couchbase/v2"

	"github.com/testcontainers/testcontainers-go"
	container "github.com/testcontainers/testcontainers-go/modules/couchbase"
)

func TestCouchbase(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rays-Fiber/Storage/Couchbase Suite")
}

type TestCouchbaseConfig struct {
	Component
}

// CouchbaseConfigOverride implements [couchbase.CouchbaseConfigOverride].
func (t *TestCouchbaseConfig) CouchbaseConfigOverride(_ string, _ *couchbase.Config) {}

var _ CouchbaseConfigOverride = &TestCouchbaseConfig{}

var _ = Describe("Couchbase", func() {
	It("can register a couchbase Storage", func() {

		ctx := context.Background()
		bucket := container.NewBucket("fiber_storage")

		bucket = bucket.WithQuota(100).
			WithReplicas(0).
			WithFlushEnabled(false).
			WithPrimaryIndex(true)

		couchbaseContainer, err := container.Run(ctx,
			"couchbase:community-8.0.0",
			testcontainers.WithLogger(GetLogger()),
			container.WithBuckets(bucket),
		)

		Expect(err).NotTo(HaveOccurred())

		ip, _ := couchbaseContainer.Host(ctx)
		port, _ := couchbaseContainer.MappedPort(ctx, "8091")

		config := map[string]any{
			"primary": "couchbase",
			"couchbase": map[string]any{
				"primary": map[string]any{
					"host": fmt.Sprintf("couchbase://%s:%s", ip, port.Port()),
				},
			},
		}

		GetContainer().LoadConfigAt(config, "rays.fiber.storage")
		AddBeamFor(&TestCouchbaseConfig{})
		GetContainer().EnableScanning()
		err = GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())

		storage := GetBeamScopedInstance[Storage]()
		Expect(reflect.TypeOf(storage).Elem()).To(Equal(reflect.TypeFor[couchbase.Storage]()))

		provider := GetBeamScopedInstance[StorageProvider]()
		Expect(provider.Provide("couchbase")).NotTo(BeNil())
		Expect(provider.Provide("couchbase:primary")).NotTo(BeNil())
		Expect(provider.Provide("couchbase:default")).NotTo(BeNil())

		testcontainers.CleanupContainer(GinkgoTB(), couchbaseContainer)
		Expect(err).NotTo(HaveOccurred())
	})
})
