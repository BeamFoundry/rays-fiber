package internal

import (
	"time"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"

	"github.com/gofiber/storage/mysql/v2"

	"github.com/gofiber/fiber/v3"
)

type MySQLConfigOverride interface {
	MySQLConfigOverride(instance string, config *mysql.Config)
}

type mySQLConfig struct {
	connectionURI string
	host          string
	port          int
	username      string
	password      string
	database      string
	table         string
	reset         bool
	gcInterval    time.Duration
}

func (this mySQLConfig) Config() mysql.Config {
	config := mysql.ConfigDefault
	config.ConnectionURI = this.connectionURI
	config.Host = this.host
	config.Port = this.port
	config.Username = this.username
	config.Password = this.password
	config.Database = this.database
	config.Table = this.table
	config.Reset = this.reset
	config.GCInterval = this.gcInterval
	return config
}

// ---

type MySQLConfigurationProperties struct {
	ConfigurationProperties `Name:"rays.fiber.storage"`
	//
	mysql map[string]mySQLConfig
}

func (this *MySQLConfigurationProperties) Config(name string) mysql.Config {
	if c, ok := this.mysql[name]; ok {
		return c.Config()
	} else {
		return this.mysql["primary"].Config()
	}
}

var _, _ = any(&MySQLConfigurationProperties{}).(Stereotype)

type mysqlStorageProvider struct {
	cfg      *MySQLConfigurationProperties
	override Optional[MySQLConfigOverride]
}

func (this *mysqlStorageProvider) Storage(name string) fiber.Storage {
	cfg := this.cfg.Config(name)
	if config, ok := this.override.Get(); ok {
		config.MySQLConfigOverride(name, &cfg)
	}
	return mysql.New(cfg)
}

type MySQLAutoConfig struct {
	Configuration

	config   *MySQLConfigurationProperties `@:"Inject"`
	override Optional[MySQLConfigOverride] `@:"Inject"`
}

func (this *MySQLAutoConfig) MySQL(_ struct {
	Qualifier `Name:"mysql"`
}) StorageEngine {
	return &mysqlStorageProvider{this.config, this.override}
}

func (this *MySQLAutoConfig) MySQLDefault(
	that struct {
		Qualifier             `Name:"mysql"`
		ConditionalOnProperty `Name:"rays.fiber.storage.primary" HavingValue:"mysql"`
		provider              StorageEngine `@:"Inject,Qualifier{mysql}"`
	},
) fiber.Storage {
	return that.provider.Storage("primary")
}

var _, _ = any(&MySQLAutoConfig{}).(Stereotype)
