package internal

import (
	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"

	"github.com/gofiber/storage/redis/v3"

	"github.com/gofiber/fiber/v3"
)

type RedisConfigOverride interface {
	RedisConfigOverride(instance string, config *redis.Config)
}

type redisConfig struct {
	host                string
	port                int
	username            string
	password            string
	database            int
	url                 string
	addrs               []string
	masterName          string
	clientName          string
	sentinelUsername    string
	sentinelPassword    string
	reset               bool
	poolSize            int
	isClusterMode       bool
	skipConnectionCheck bool
}

func (this redisConfig) Config() redis.Config {
	config := redis.ConfigDefault
	config.Host = this.host
	config.Port = this.port
	config.Username = this.username
	config.Password = this.password
	config.Database = this.database
	config.URL = this.url
	config.Addrs = this.addrs
	config.MasterName = this.masterName
	config.ClientName = this.clientName
	config.SentinelUsername = this.sentinelUsername
	config.SentinelPassword = this.sentinelPassword
	config.Reset = this.reset
	config.PoolSize = this.poolSize
	config.IsClusterMode = this.isClusterMode
	config.SkipConnectionCheck = this.skipConnectionCheck
	return config
}

type RedisConfigurationProperties struct {
	ConfigurationProperties `Name:"rays.fiber.storage"`
	//
	redis map[string]redisConfig
}

func (this *RedisConfigurationProperties) Config(name string) redis.Config {
	c, ok := this.redis[name]
	if !ok {
		c = this.redis["primary"]
	}
	return c.Config()
}

var _, _ = any(&RedisConfigurationProperties{}).(Stereotype)

type redisStorageProvider struct {
	cfg      *RedisConfigurationProperties
	override Optional[RedisConfigOverride]
}

func (this *redisStorageProvider) Storage(name string) fiber.Storage {
	cfg := this.cfg.Config(name)
	if config, ok := this.override.Get(); ok {
		config.RedisConfigOverride(name, &cfg)
	}
	return redis.New(cfg)
}

type RedisAutoConfig struct {
	Configuration

	config   *RedisConfigurationProperties `@:"Inject"`
	override Optional[RedisConfigOverride] `@:"Inject"`
}

func (this *RedisAutoConfig) Redis(_ struct {
	Qualifier `Name:"redis"`
}) StorageEngine {
	return &redisStorageProvider{this.config, this.override}
}

func (this *RedisAutoConfig) RedisDefault(
	that struct {
		Qualifier             `Name:"redis"`
		ConditionalOnProperty `Name:"rays.fiber.storage.primary" HavingValue:"redis"`
		provider              StorageEngine `@:"Inject,Qualifier{redis}"`
	},
) fiber.Storage {
	return that.provider.Storage("primary")
}

var _, _ = any(&RedisAutoConfig{}).(Stereotype)
