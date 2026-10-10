package internal_test

import (
	"testing"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"
	. "github.com/BeamFoundry/rays/pkg/testing"

	. "github.com/BeamFoundry/rays-fiber/pkg/middleware/internal"

	"github.com/gofiber/fiber/v3/middleware/compress"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestInternal(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pkg/Middleware/Internal Suite")
}

type compressConfigurationProperties struct {
	level string
}

type compressConfig struct {
	compress compressConfigurationProperties
}

func (this *compressConfig) Config() compress.Config {
	config := compress.ConfigDefault
	switch String(this.compress.level).ToUpperCase().String() {
	case "DISABLED":
		config.Level = compress.LevelDisabled
	case "BESTSPEED", "SPEED":
		config.Level = compress.LevelBestSpeed
	case "BESTCOMPRESSION", "COMPRESSION":
		config.Level = compress.LevelBestCompression
	default:
		config.Level = compress.LevelDefault
	}
	return config
}

type TestConfigComponent struct {
	Component
	Server `Name:"primary"`

	config *FiberMwConfig[compressConfig] `@:"Inject"`
}

func (this *TestConfigComponent) PostConstruct() Error {
	Expect(this.config.GetConfig(this)).NotTo(BeNil())
	return nil
}

var _ = Describe("FiberMwConfig", func() {
	It("can be created", func() {
		config := map[string]any{
			"level": "default",
		}
		Expect(GetContainer().LoadConfigAt(config, "rays.fiber.servers.primary.compress")).NotTo(HaveOccurred())
		AddBeamFor(&TestConfigComponent{})
		AddBeamFor(&FiberMwConfig[compressConfig]{})
		Expect(GetContainer().Initialize()).NotTo(HaveOccurred())
	})
})
