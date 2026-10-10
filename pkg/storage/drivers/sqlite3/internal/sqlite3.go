package internal3

import (
	"time"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"

	"github.com/gofiber/storage/sqlite3/v2"

	"github.com/gofiber/fiber/v3"
)

type SQLite3ConfigOverride interface {
	SQLite3ConfigOverride(instance string, config *sqlite3.Config)
}

type sqlite3Config struct {
	database        string
	table           string
	reset           bool
	gcInterval      time.Duration
	maxIdleConns    int
	maxOpenConns    int
	connMaxLifetime time.Duration
}

func (this sqlite3Config) Config() sqlite3.Config {
	config := sqlite3.ConfigDefault
	config.Database = this.database
	config.Table = this.table
	config.Reset = this.reset
	config.GCInterval = this.gcInterval
	config.MaxIdleConns = this.maxIdleConns
	config.MaxOpenConns = this.maxOpenConns
	config.ConnMaxLifetime = this.connMaxLifetime
	return config
}

// ---

type SQLite3ConfigurationProperties struct {
	ConfigurationProperties `Name:"rays.fiber.storage"`
	//
	sqlite3 map[string]sqlite3Config
}

func (this *SQLite3ConfigurationProperties) Config(name string) sqlite3.Config {
	if c, ok := this.sqlite3[name]; ok {
		return c.Config()
	} else {
		return this.sqlite3["primary"].Config()
	}
}

var _, _ = any(&SQLite3ConfigurationProperties{}).(Stereotype)

type sqlite3StorageProvider struct {
	cfg      *SQLite3ConfigurationProperties
	override Optional[SQLite3ConfigOverride]
}

func (this *sqlite3StorageProvider) Storage(name string) fiber.Storage {
	cfg := this.cfg.Config(name)
	if config, ok := this.override.Get(); ok {
		config.SQLite3ConfigOverride(name, &cfg)
	}
	return sqlite3.New(cfg)
}

type SQLite3AutoConfig struct {
	Configuration

	config   *SQLite3ConfigurationProperties `@:"Inject"`
	override Optional[SQLite3ConfigOverride] `@:"Inject"`
}

func (this *SQLite3AutoConfig) Sqlite3(_ struct {
	Qualifier `Name:"sqlite3"`
}) StorageEngine {
	return &sqlite3StorageProvider{this.config, this.override}
}

func (this *SQLite3AutoConfig) Sqlite3Default(
	that struct {
		Qualifier             `Name:"sqlite3"`
		ConditionalOnProperty `Name:"rays.fiber.storage.primary" HavingValue:"sqlite3"`
		provider              StorageEngine `@:"Inject,Qualifier{sqlite3}"`
	},
) fiber.Storage {
	return that.provider.Storage("primary")
}

var _, _ = any(&SQLite3AutoConfig{}).(Stereotype)
