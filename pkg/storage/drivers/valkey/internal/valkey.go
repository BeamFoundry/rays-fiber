package internal

import (
	"time"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	"github.com/BeamFoundry/rays-web-security/pkg/certs"

	"github.com/gofiber/storage/valkey"

	"github.com/gofiber/fiber/v3"
)

type ValkeyConfigOverride interface {
	ValkeyConfigOverride(instance string, config *valkey.Config)
}

type ValkeyClientTLSProvider interface {
	ValkeyClientTLSProvider(instance string) certs.ClientTLSProvider
}

type valkeyConfig struct {
	password            string
	clientName          string
	url                 string
	username            string
	initAddress         []string
	readBufferEachConn  int
	cacheSizeEachConn   int
	ringScaleEachConn   int
	selectDB            int
	writeBufferEachConn int
	blockingPoolSize    int
	pipelineMultiplex   int
	cacheTTL            time.Duration
	disableRetry        bool
	disableCache        bool
	alwaysPipelining    bool
	reset               bool
}

func (this *valkeyConfig) Config() valkey.Config {
	config := valkey.ConfigDefault
	config.Username = this.username
	config.Password = this.password
	config.ClientName = this.clientName
	config.URL = this.url
	config.SelectDB = this.selectDB
	config.InitAddress = this.initAddress
	config.CacheSizeEachConn = this.cacheSizeEachConn
	config.RingScaleEachConn = this.ringScaleEachConn
	config.ReadBufferEachConn = this.readBufferEachConn
	config.WriteBufferEachConn = this.writeBufferEachConn
	config.BlockingPoolSize = this.blockingPoolSize
	config.PipelineMultiplex = this.pipelineMultiplex
	config.DisableRetry = this.disableRetry
	config.DisableCache = this.disableCache
	config.AlwaysPipelining = this.alwaysPipelining
	config.Reset = this.reset
	config.CacheTTL = this.cacheTTL
	return config
}

// ---

type ValkeyConfigurationProperties struct {
	ConfigurationProperties `Name:"rays.fiber.storage"`
	//
	valkey map[string]valkeyConfig
}

func (this *ValkeyConfigurationProperties) Config(name string) valkey.Config {
	c, ok := this.valkey[name]
	if !ok {
		c = this.valkey["primary"]
	}
	return c.Config()
}

var _, _ = any(&ValkeyConfigurationProperties{}).(Stereotype)

type valkeyStorageProvider struct {
	cfg         *ValkeyConfigurationProperties
	override    Optional[ValkeyConfigOverride]
	tlsProvider Optional[ValkeyClientTLSProvider]
}

func (this *valkeyStorageProvider) Storage(name string) fiber.Storage {
	cfg := this.cfg.Config(name)
	if tlsProvider, ok := this.tlsProvider.Get(); ok {
		if clientTLS := tlsProvider.ValkeyClientTLSProvider(name); clientTLS != nil {
			tlsConfig := RaiseVal(certs.ConfigureClientTLS(clientTLS))
			if tlsConfig != nil {
				cfg.TLSConfig = tlsConfig
			}
		}
	}
	if config, ok := this.override.Get(); ok {
		config.ValkeyConfigOverride(name, &cfg)
	}
	return valkey.New(cfg)
}

type ValkeyAutoConfig struct {
	Configuration

	config      *ValkeyConfigurationProperties    `@:"Inject"`
	override    Optional[ValkeyConfigOverride]    `@:"Inject"`
	tlsProvider Optional[ValkeyClientTLSProvider] `@:"Inject"`
}

func (this *ValkeyAutoConfig) Valkey(_ struct {
	Qualifier `Name:"valkey"`
},
) StorageEngine {
	return &valkeyStorageProvider{this.config, this.override, this.tlsProvider}
}

func (this *ValkeyAutoConfig) ValkeyDefault(
	that struct {
		Qualifier             `Name:"valkey"`
		ConditionalOnProperty `Name:"rays.fiber.storage.primary" HavingValue:"valkey"`
		provider              StorageEngine `@:"Inject,Qualifier{valkey}"`
	},
) fiber.Storage {
	return that.provider.Storage("primary")
}

var _, _ = any(&ValkeyAutoConfig{}).(Stereotype)
