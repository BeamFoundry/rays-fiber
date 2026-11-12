package internal

import (
	"time"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	"github.com/BeamFoundry/rays-web-security/pkg/certs"

	"github.com/gofiber/storage/rueidis"

	"github.com/gofiber/fiber/v3"
)

type RueidisConfigOverride interface {
	RueidisConfigOverride(instance string, config *rueidis.Config)
}

type RueidisClientTLSProvider interface {
	RueidisClientTLSProvider(instance string) certs.ClientTLSProvider
}

type rueidisConfig struct {
	password                string
	clientName              string
	url                     string
	username                string
	initAddress             []string
	readBufferEachConn      int
	cacheSizeEachConn       int
	ringScaleEachConn       int
	selectDB                int
	writeBufferEachConn     int
	blockingPoolSize        int
	pipelineMultiplex       int
	cacheTTL                time.Duration
	disableRetry            bool
	disableCache            bool
	alwaysPipelining        bool
	disableAlwaysPipelining bool
	reset                   bool
}

func (this *rueidisConfig) Config() rueidis.Config {
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
	c, ok := this.rueidis[name]
	if !ok {
		c = this.rueidis["primary"]
	}
	return c.Config()
}

var _, _ = any(&RueidisConfigurationProperties{}).(Stereotype)

type rueidisStorageProvider struct {
	cfg         *RueidisConfigurationProperties
	override    Optional[RueidisConfigOverride]
	tlsProvider Optional[RueidisClientTLSProvider]
}

func (this *rueidisStorageProvider) Storage(name string) fiber.Storage {
	cfg := this.cfg.Config(name)
	if tlsProvider, ok := this.tlsProvider.Get(); ok {
		if clientTLS := tlsProvider.RueidisClientTLSProvider(name); clientTLS != nil {
			tlsConfig := RaiseVal(certs.ConfigureClientTLS(clientTLS))
			if tlsConfig != nil {
				cfg.TLSConfig = tlsConfig
			}
		}
	}
	if config, ok := this.override.Get(); ok {
		config.RueidisConfigOverride(name, &cfg)
	}
	return rueidis.New(cfg)
}

type RueidisAutoConfig struct {
	Configuration

	config      *RueidisConfigurationProperties    `@:"Inject"`
	override    Optional[RueidisConfigOverride]    `@:"Inject"`
	tlsProvider Optional[RueidisClientTLSProvider] `@:"Inject"`
}

func (this *RueidisAutoConfig) Rueidis(_ struct {
	Qualifier `Name:"rueidis"`
},
) StorageEngine {
	return &rueidisStorageProvider{this.config, this.override, this.tlsProvider}
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
