package internal

import (
	"time"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"

	"github.com/gofiber/storage/couchbase/v2"

	"github.com/gofiber/fiber/v3"
)

type CouchbaseConfigOverride interface {
	CouchbaseConfigOverride(instance string, config *couchbase.Config)
}

type CouchbaseConfig struct {
	username          string
	password          string
	host              string
	bucket            string
	connectionTimeout time.Duration
	kvTimeout         time.Duration
}

func (this CouchbaseConfig) Config() couchbase.Config {
	config := couchbase.ConfigDefault
	config.Username = this.username
	config.Password = this.password
	config.Host = this.host
	config.Bucket = this.bucket
	config.ConnectionTimeout = this.connectionTimeout
	config.KVTimeout = this.kvTimeout
	return config
}

// ---

type CouchbaseConfigurationProperties struct {
	ConfigurationProperties `Name:"rays.fiber.storage"`
	//
	couchbase map[string]CouchbaseConfig
}

func (this *CouchbaseConfigurationProperties) Config(name string) couchbase.Config {
	if c, ok := this.couchbase[name]; ok {
		return c.Config()
	} else {
		return this.couchbase["primary"].Config()
	}
}

var _, _ = (any)(&CouchbaseConfigurationProperties{}).(Stereotype)

type couchbaseStorageProvider struct {
	cfg      *CouchbaseConfigurationProperties
	override Optional[CouchbaseConfigOverride]
}

func (this *couchbaseStorageProvider) Storage(name string) fiber.Storage {
	cfg := this.cfg.Config(name)
	if config, ok := this.override.Get(); ok {
		config.CouchbaseConfigOverride(name, &cfg)
	}
	return couchbase.New(cfg)
}

type CouchbaseAutoConfig struct {
	Configuration

	config   *CouchbaseConfigurationProperties `@:"Inject"`
	override Optional[CouchbaseConfigOverride] `@:"Inject"`
}

func (this *CouchbaseAutoConfig) Couchbase(_ struct {
	Qualifier `Name:"couchbase"`
}) StorageEngine {
	return &couchbaseStorageProvider{this.config, this.override}
}

func (this *CouchbaseAutoConfig) CouchbaseDefault(
	that struct {
		Qualifier             `Name:"couchbase"`
		ConditionalOnProperty `Name:"rays.fiber.storage.primary" HavingValue:"couchbase"`
		provider              StorageEngine `@:"Inject,Qualifier{couchbase}"`
	},
) fiber.Storage {
	return that.provider.Storage("primary")
}

var _, _ = any(&CouchbaseAutoConfig{}).(Stereotype)
