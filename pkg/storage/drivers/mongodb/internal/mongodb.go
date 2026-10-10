package internal

import (
	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"

	"github.com/gofiber/storage/mongodb/v2"

	"github.com/gofiber/fiber/v3"
)

type MongoDBConfigOverride interface {
	MongoDBConfigOverride(instance string, config *mongodb.Config)
}

type MongoDBConfig struct {
	connectionURI string
	host          string
	port          int
	username      string
	password      string
	database      string
	collection    string
	reset         bool
}

func (this MongoDBConfig) Config() mongodb.Config {
	return mongodb.Config{
		ConnectionURI: this.connectionURI,
		Host:          this.host,
		Port:          this.port,
		Username:      this.username,
		Password:      this.password,
		Database:      this.database,
		Collection:    this.collection,
		Reset:         this.reset,
	}
}

// ---

type MongoDBConfigurationProperties struct {
	ConfigurationProperties `Name:"rays.fiber.storage"`
	//
	mongodb map[string]MongoDBConfig
}

func (this *MongoDBConfigurationProperties) Config(name string) mongodb.Config {
	if c, ok := this.mongodb[name]; ok {
		return c.Config()
	} else {
		return this.mongodb["primary"].Config()
	}
}

// ---

type mongodbStorageProvider struct {
	cfg      *MongoDBConfigurationProperties
	override Optional[MongoDBConfigOverride]
}

func (this *mongodbStorageProvider) Storage(name string) fiber.Storage {
	cfg := this.cfg.Config(name)
	if config, ok := this.override.Get(); ok {
		config.MongoDBConfigOverride(name, &cfg)
	}
	return mongodb.New(cfg)
}

type MongoDBAutoConfig struct {
	Configuration

	config   *MongoDBConfigurationProperties `@:"Inject"`
	override Optional[MongoDBConfigOverride] `@:"Inject"`
}

func (this *MongoDBAutoConfig) Mongodb(_ struct {
	Qualifier `Name:"mongodb"`
}) StorageEngine {
	return &mongodbStorageProvider{this.config, this.override}
}

func (this *MongoDBAutoConfig) MongodbDefault(
	that struct {
		Qualifier             `Name:"mongodb"`
		ConditionalOnProperty `Name:"rays.fiber.storage.primary" HavingValue:"mongodb"`
		provider              StorageEngine `@:"Inject,Qualifier{mongodb}"`
	},
) fiber.Storage {
	return that.provider.Storage("primary")
}

var _, _ = any(&MongoDBAutoConfig{}).(Stereotype)
