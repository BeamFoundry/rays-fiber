package internal

import (
	"strconv"
	"strings"
	"time"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"

	as "github.com/aerospike/aerospike-client-go/v8"
	"github.com/gofiber/storage/aerospike"

	"github.com/gofiber/fiber/v3"
)

type AerospikeConfigOverride interface {
	AerospikeConfigOverride(instance string, config *aerospike.Config)
}

// ---

type aerospikeConfig struct {
	hosts                    []string
	namespace                string
	setName                  string
	reset                    bool
	expiration               time.Duration
	schemaVersion            int
	schemaDescription        string
	forceSchemaUpdate        bool
	initialConnectionTimeout time.Duration
}

func (this aerospikeConfig) Config() aerospike.Config {
	hosts := []*as.Host{}
	for _, hostport := range this.hosts {
		if host, p, ok := strings.Cut(hostport, ":"); ok {
			if port, err := strconv.Atoi(p); err == nil {
				hosts = append(hosts, as.NewHost(host, port))
			}
		}
	}
	config := aerospike.ConfigDefault
	config.Hosts = hosts
	config.Namespace = this.namespace
	config.SetName = this.setName
	config.Reset = this.reset
	config.Expiration = this.expiration
	config.SchemaVersion = this.schemaVersion
	config.SchemaDescription = this.schemaDescription
	config.ForceSchemaUpdate = this.forceSchemaUpdate
	config.InitialConnectionTimeout = this.initialConnectionTimeout
	return config
}

// ---

type AerospikeConfigurationProperties struct {
	ConfigurationProperties `Name:"rays.fiber.storage"`
	//
	aerospike map[string]aerospikeConfig
}

var _, _ = any(&AerospikeConfigurationProperties{}).(Stereotype)

func (this *AerospikeConfigurationProperties) Config(name string) aerospike.Config {
	if c, ok := this.aerospike[name]; ok {
		return c.Config()
	} else {
		return this.aerospike["primary"].Config()
	}
}

type aerospikeStorageProvider struct {
	cfg      *AerospikeConfigurationProperties
	override Optional[AerospikeConfigOverride]
}

func (this *aerospikeStorageProvider) Storage(name string) fiber.Storage {
	cfg := this.cfg.Config(name)
	if config, ok := this.override.Get(); ok {
		config.AerospikeConfigOverride(name, &cfg)
	}
	return aerospike.New(cfg)
}

// ---

type AerospikeAutoConfig struct {
	Configuration

	config   *AerospikeConfigurationProperties `@:"Inject"`
	override Optional[AerospikeConfigOverride] `@:"Inject"`
}

func (this *AerospikeAutoConfig) Aerospike(_ struct {
	Qualifier `Name:"aerospike"`
}) StorageEngine {
	return &aerospikeStorageProvider{this.config, this.override}
}

func (this *AerospikeAutoConfig) AerospikeDefault(
	that struct {
		Qualifier             `Name:"aerospike"`
		ConditionalOnProperty `Name:"rays.fiber.storage.primary" HavingValue:"aerospike"`
		provider              StorageEngine `@:"Inject,Qualifier{aerospike}"`
	},
) fiber.Storage {
	return that.provider.Storage("primary")
}

var _, _ = any(&AerospikeAutoConfig{}).(Stereotype)
