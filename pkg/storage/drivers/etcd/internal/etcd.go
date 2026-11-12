package internal

import (
	"time"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	"github.com/BeamFoundry/rays-web-security/pkg/certs"

	"github.com/gofiber/storage/etcd/v2"

	"github.com/gofiber/fiber/v3"
)

type EtcdConfigOverride interface {
	EtcdConfigOverride(instance string, config *etcd.Config)
}

type EtcdClientTLSProvider interface {
	EtcdClientTLSProvider(instance string) certs.ClientTLSProvider
}

type etcdConfig struct {
	username    string
	password    string
	endpoints   []string
	dialTimeout time.Duration
}

func (this etcdConfig) Config() etcd.Config {
	config := etcd.ConfigDefault
	config.Endpoints = this.endpoints
	config.DialTimeout = this.dialTimeout
	config.Username = this.username
	config.Password = this.password
	return config
}

// ---

type EtcdConfigurationProperties struct {
	ConfigurationProperties `Name:"rays.fiber.storage"`
	//
	etcd map[string]etcdConfig
}

var _, _ = any(&EtcdConfigurationProperties{}).(Stereotype)

func (this *EtcdConfigurationProperties) Config(name string) etcd.Config {
	c, ok := this.etcd[name]
	if !ok {
		c = this.etcd["primary"]
	}
	return c.Config()
}

type etcdStorageProvider struct {
	cfg         *EtcdConfigurationProperties
	override    Optional[EtcdConfigOverride]
	tlsProvider Optional[EtcdClientTLSProvider]
}

func (this *etcdStorageProvider) Storage(name string) fiber.Storage {
	cfg := this.cfg.Config(name)
	if tlsProvider, ok := this.tlsProvider.Get(); ok {
		if clientTLS := tlsProvider.EtcdClientTLSProvider(name); clientTLS != nil {
			tlsConfig := RaiseVal(certs.ConfigureClientTLS(clientTLS))
			if tlsConfig != nil {
				cfg.TLS = tlsConfig
			}
		}
	}
	if config, ok := this.override.Get(); ok {
		config.EtcdConfigOverride(name, &cfg)
	}
	return etcd.New(cfg)
}

type EtcdAutoConfig struct {
	Configuration

	config      *EtcdConfigurationProperties    `@:"Inject"`
	override    Optional[EtcdConfigOverride]    `@:"Inject"`
	tlsProvider Optional[EtcdClientTLSProvider] `@:"Inject"`
}

func (this *EtcdAutoConfig) Etcd(_ struct {
	Qualifier `Name:"etcd"`
},
) StorageEngine {
	return &etcdStorageProvider{this.config, this.override, this.tlsProvider}
}

func (this *EtcdAutoConfig) EtcdDefault(
	that struct {
		Qualifier             `Name:"etcd"`
		ConditionalOnProperty `Name:"rays.fiber.storage.primary" HavingValue:"etcd"`
		provider              StorageEngine `@:"Inject,Qualifier{etcd}"`
	},
) fiber.Storage {
	return that.provider.Storage("primary")
}

var _, _ = any(&EtcdAutoConfig{}).(Stereotype)
