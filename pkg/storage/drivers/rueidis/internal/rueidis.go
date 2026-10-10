package internal

import (
	"time"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"

	"github.com/gofiber/storage/rueidis"

	"github.com/gofiber/fiber/v3"
)

type RueidisConfigOverride interface {
	RueidisConfigOverride(instance string, config *rueidis.Config)
}

type rueidisConfig struct {
	username                string
	password                string
	clientName              string
	url                     string
	selectDB                int
	initAddress             []string
	cacheSizeEachConn       int
	ringScaleEachConn       int
	readBufferEachConn      int
	writeBufferEachConn     int
	blockingPoolSize        int
	pipelineMultiplex       int
	disableRetry            bool
	disableCache            bool
	alwaysPipelining        bool
	disableAlwaysPipelining bool
	reset                   bool
	cacheTTL                time.Duration
}

func (this rueidisConfig) Config() rueidis.Config {
	cfg := rueidis.ConfigDefault
	cfg.Username = this.username
	cfg.Password = this.password
	cfg.ClientName = this.clientName
	cfg.URL = this.url
	cfg.SelectDB = this.selectDB
	cfg.InitAddress = this.initAddress
	cfg.CacheSizeEachConn = this.cacheSizeEachConn
	cfg.RingScaleEachConn = this.ringScaleEachConn
	cfg.ReadBufferEachConn = this.readBufferEachConn
	cfg.WriteBufferEachConn = this.writeBufferEachConn
	cfg.BlockingPoolSize = this.blockingPoolSize
	cfg.PipelineMultiplex = this.pipelineMultiplex
	cfg.DisableRetry = this.disableRetry
	cfg.DisableCache = this.disableCache
	cfg.AlwaysPipelining = this.alwaysPipelining
	cfg.DisableAlwaysPipelining = this.disableAlwaysPipelining
	cfg.Reset = this.reset
	cfg.CacheTTL = this.cacheTTL
	return cfg
}

// ---

type RueidisConfigurationProperties struct {
	ConfigurationProperties `Name:"rays.fiber.storage"`
	//
	rueidis map[string]rueidisConfig
}

func (this *RueidisConfigurationProperties) Config(name string) rueidis.Config {
	if c, ok := this.rueidis[name]; ok {
		return c.Config()
	} else {
		return this.rueidis["primary"].Config()
	}
}

var _, _ = any(&RueidisConfigurationProperties{}).(Stereotype)

type rueidisStorageProvider struct {
	cfg      *RueidisConfigurationProperties
	override Optional[RueidisConfigOverride]
}

func (this *rueidisStorageProvider) Storage(name string) fiber.Storage {
	cfg := this.cfg.Config(name)
	if config, ok := this.override.Get(); ok {
		config.RueidisConfigOverride(name, &cfg)
	}
	return rueidis.New(cfg)
}

type RueidisAutoConfig struct {
	Configuration

	config   *RueidisConfigurationProperties `@:"Inject"`
	override Optional[RueidisConfigOverride] `@:"Inject"`
}

func (this *RueidisAutoConfig) Rueidis(_ struct {
	Qualifier `Name:"rueidis"`
}) StorageEngine {
	return &rueidisStorageProvider{this.config, this.override}
}

func (this *RueidisAutoConfig) RueidisDefault(
	that struct {
		Qualifier             `Name:"rueidis"`
		ConditionalOnProperty `Name:"rays.fiber.storage.primary" HavingValue:"rueidis"`
		provider              StorageEngine `@:"Inject,Qualifier{rueidis}"`
	},
) fiber.Storage {
	return that.provider.Storage("primary")
}

var _, _ = any(&RueidisAutoConfig{}).(Stereotype)
