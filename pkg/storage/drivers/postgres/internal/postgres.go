package internal

import (
	"time"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"

	"github.com/gofiber/storage/postgres/v3"

	"github.com/gofiber/fiber/v3"
)

type PostgresConfigOverride interface {
	PostgresConfigOverride(instance string, config *postgres.Config)
}

type postgresConfig struct {
	connectionURI string
	host          string
	port          int
	username      string
	password      string
	database      string
	table         string
	sslMode       string
	reset         bool
	gcInterval    time.Duration
}

func (this postgresConfig) Config() postgres.Config {
	config := postgres.ConfigDefault
	config.ConnectionURI = this.connectionURI
	config.Host = this.host
	config.Port = this.port
	config.Username = this.username
	config.Password = this.password
	config.Database = this.database
	config.Table = this.table
	config.SSLMode = this.sslMode
	config.Reset = this.reset
	config.GCInterval = this.gcInterval
	return config
}

// ---

type PostgresConfigurationProperties struct {
	ConfigurationProperties `Name:"rays.fiber.storage"`
	//
	postgres map[string]postgresConfig
}

func (this *PostgresConfigurationProperties) Config(name string) postgres.Config {
	if c, ok := this.postgres[name]; ok {
		return c.Config()
	} else {
		return this.postgres["primary"].Config()
	}
}

var _, _ = any(&PostgresConfigurationProperties{}).(Stereotype)

type postgresStorageProvider struct {
	cfg      *PostgresConfigurationProperties
	override Optional[PostgresConfigOverride]
}

func (this *postgresStorageProvider) Storage(name string) fiber.Storage {
	cfg := this.cfg.Config(name)
	if config, ok := this.override.Get(); ok {
		config.PostgresConfigOverride(name, &cfg)
	}
	return postgres.New(cfg)
}

type PostgresAutoConfig struct {
	Configuration

	config   *PostgresConfigurationProperties `@:"Inject"`
	override Optional[PostgresConfigOverride] `@:"Inject"`
}

func (this *PostgresAutoConfig) Postgres(_ struct {
	Qualifier `Name:"postgres"`
}) StorageEngine {
	return &postgresStorageProvider{this.config, this.override}
}

func (this *PostgresAutoConfig) PostgresDefault(
	that struct {
		Qualifier             `Name:"postgres"`
		ConditionalOnProperty `Name:"rays.fiber.storage.primary" HavingValue:"postgres"`
		provider              StorageEngine `@:"Inject,Qualifier{postgres}"`
	},
) fiber.Storage {
	return that.provider.Storage("primary")
}

var _, _ = any(&PostgresAutoConfig{}).(Stereotype)
