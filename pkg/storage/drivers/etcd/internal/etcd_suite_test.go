package internal_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"reflect"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	_ "github.com/BeamFoundry/rays-fiber"
	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/etcd"

	. "github.com/gofiber/storage"

	"github.com/BeamFoundry/rays-web-security/pkg/certs"
	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"
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
	capturedTLS *tls.Config
}

// EtcdConfigOverride implements [etcd.EtcdConfigOverride].
func (t *TestEtcdConfig) EtcdConfigOverride(_ string, config *etcd.Config) {
	t.capturedTLS = config.TLS
}

var _ EtcdConfigOverride = &TestEtcdConfig{}

type TestEtcdTLSProvider struct {
	Component
}

func (t *TestEtcdTLSProvider) ProvideClientTLS() (*tls.Config, Error) {
	return &tls.Config{ServerName: "etcd-tls"}, nil
}

func (t *TestEtcdTLSProvider) EtcdClientTLSProvider(instance string) certs.ClientTLSProvider {
	if instance == "primary" {
		return t
	}
	return nil
}

var _ EtcdClientTLSProvider = &TestEtcdTLSProvider{}

type TestEtcdConfigWithTLS struct {
	Component
	capturedTLS *tls.Config
}

func (t *TestEtcdConfigWithTLS) ProvideClientTLS() (*tls.Config, Error) {
	return &tls.Config{ServerName: "etcd-from-override"}, nil
}

func (t *TestEtcdConfigWithTLS) EtcdClientTLSProvider(instance string) certs.ClientTLSProvider {
	if instance == "primary" {
		return t
	}
	return nil
}

func (t *TestEtcdConfigWithTLS) EtcdConfigOverride(_ string, config *etcd.Config) {
	t.capturedTLS = config.TLS
}

var (
	_ EtcdConfigOverride    = &TestEtcdConfigWithTLS{}
	_ EtcdClientTLSProvider = &TestEtcdConfigWithTLS{}
)

var _ = Describe("Etcd", func() {
	It("can register a etcd Storage", func() {
		ctx := context.Background()
		etcdContainer, err := container.Run(ctx, "gcr.io/etcd-development/etcd:v3.5.14", testcontainers.WithLogger(GetLogger()))

		Expect(err).NotTo(HaveOccurred())

		ip, err := etcdContainer.Host(ctx)
		Expect(err).NotTo(HaveOccurred())
		port, err := etcdContainer.MappedPort(ctx, "2379")
		Expect(err).NotTo(HaveOccurred())

		endpoints := fmt.Sprintf("%s:%s", ip, port.Port())
		config := map[string]any{
			"primary": "etcd",
			"etcd": map[string]any{
				"primary": map[string]any{
					"endpoints": []any{endpoints},
				},
			},
		}

		Expect(GetContainer().LoadConfigAt(config, "rays.fiber.storage")).To(Succeed())
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

	It("can configure TLS via EtcdClientTLSProvider", func() {
		config := map[string]any{
			"primary": "etcd",
			"etcd": map[string]any{
				"primary": map[string]any{
					"endpoints": []any{"127.0.0.1:2379"},
				},
			},
		}

		Expect(GetContainer().LoadConfigAt(config, "rays.fiber.storage")).To(Succeed())
		cfgBeam := &TestEtcdConfig{}
		AddBeamFor(cfgBeam)
		AddBeamFor(&TestEtcdTLSProvider{})
		GetContainer().EnableScanning()

		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())

		provider := GetBeamScopedInstance[StorageProvider]()
		st := provider.Provide("etcd")
		Expect(st).NotTo(BeNil())
		Expect(cfgBeam.capturedTLS).NotTo(BeNil())
		Expect(cfgBeam.capturedTLS.ServerName).To(Equal("etcd-tls"))
	})

	It("can configure TLS via EtcdConfigOverride implementing ClientTLSProvider", func() {
		config := map[string]any{
			"primary": "etcd",
			"etcd": map[string]any{
				"primary": map[string]any{
					"endpoints": []any{"127.0.0.1:2379"},
				},
			},
		}

		Expect(GetContainer().LoadConfigAt(config, "rays.fiber.storage")).To(Succeed())
		cfgBeam := &TestEtcdConfigWithTLS{}
		AddBeamFor(cfgBeam)
		GetContainer().EnableScanning()

		err := GetContainer().Initialize()
		Expect(err).NotTo(HaveOccurred())

		provider := GetBeamScopedInstance[StorageProvider]()
		st := provider.Provide("etcd")
		Expect(st).NotTo(BeNil())
		Expect(cfgBeam.capturedTLS).NotTo(BeNil())
		Expect(cfgBeam.capturedTLS.ServerName).To(Equal("etcd-from-override"))
	})
})
