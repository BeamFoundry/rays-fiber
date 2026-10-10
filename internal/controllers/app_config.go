package controllers

import (
	"time"

	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/core"

	"github.com/gofiber/fiber/v3"
)

//nolint:unused,govet //CopyPriveToPublic + Properties
type FiberAppConfiguration struct {
	// UnportedFields :
	// ColorScheme Colors
	// ErrorHandler fiber.ErrorHandler
	// JSONDecoder utils.JSONUnmarshal
	// JSONEncoder utils.JSONMarshal
	// StructValidator validator.StructValidator
	// TrustProxyConfig http.TrustProxyConfig
	// Views fiber.Views
	// XMLDecoder utils.XMLUnmarshal
	// XMLEncoder utils.XMLMarshal
	appName                      string
	bodyLimit                    int
	caseSensitive                bool
	compressedFileSuffixes       map[string]string
	concurrency                  int
	disableDefaultContentType    bool
	disableDefaultDate           bool
	disableHeaderNormalizing     bool
	disableKeepalive             bool
	disablePreParseMultipartForm bool
	enableIPValidation           bool
	enableSplittingOnParsers     bool
	getOnly                      bool
	idleTimeout                  time.Duration
	immutable                    bool
	passLocalsToViews            bool
	proxyHeader                  string
	readBufferSize               int
	readTimeout                  time.Duration
	reduceMemoryUsage            bool
	requestMethods               []string
	serverHeader                 string
	streamRequestBody            bool
	strictRouting                bool
	trustProxy                   bool
	unescapePath                 bool
	viewsLayout                  string
	writeBufferSize              int
	writeTimeout                 time.Duration
}

// TODO: simplify this
//
//nolint:gocritic // Constructor
func (this FiberAppConfiguration) Config() fiber.Config {
	config := fiber.Config{}
	CopyPrivateToPublic(&this, &config)
	config.GETOnly = this.getOnly
	return config
}

//nolint:unused //Properties
type FiberListenerConfiguration struct {
	/*
		certClientFile        string
		certFile              string
		certKeyFile           string
	*/
	listenerNetwork       string
	shutdownTimeout       time.Duration
	tlsMinVersion         uint16
	disableStartupMessage bool
	enablePrefork         bool
	enablePrintRoutes     bool
}

func (this FiberListenerConfiguration) Config() fiber.ListenConfig {
	config := fiber.ListenConfig{}
	CopyPrivateToPublic(&this, &config)
	config.TLSMinVersion = this.tlsMinVersion
	return config
}

type FiberServerConfiguration struct {
	address  string
	config   FiberAppConfiguration
	listener FiberListenerConfiguration
}

type FiberServersConfigurationProperties struct {
	ConfigurationProperties `Name:"rays.fiber"`
	servers                 map[string]FiberServerConfiguration
}

func (this *FiberServersConfigurationProperties) Exists(name string) bool {
	_, ok := this.servers[name]
	return ok
}

func (this *FiberServersConfigurationProperties) Address(name string) string {
	return this.servers[name].address
}

func (this *FiberServersConfigurationProperties) Config(name string) fiber.Config {
	return this.servers[name].config.Config()
}

func (this *FiberServersConfigurationProperties) Listener(name string) fiber.ListenConfig {
	return this.servers[name].listener.Config()
}

var _, _ = any(&FiberServersConfigurationProperties{}).(Stereotype)

type FiberAppServerConfig struct {
	Component
	config             *FiberServersConfigurationProperties `@:"Inject"`
	configCustomizer   map[string][]FiberConfig             `@:"Inject{GroupBy=Server,Default=primary}"`
	listenerCustomizer map[string][]FiberListenConfig       `@:"Inject{GroupBy=Server,Default=primary}"`
}

func (this *FiberAppServerConfig) Exists(name string) bool {
	return this.config.Exists(name)
}

func (this *FiberAppServerConfig) Address(name string) string {
	return this.config.Address(name)
}

func (this *FiberAppServerConfig) Config(name string) fiber.Config {
	config := this.config.Config(name)
	if apps, ok := this.configCustomizer[name]; ok {
		for _, app := range apps {
			app.FiberConfig(&config)
		}
	}
	return config
}

func (this *FiberAppServerConfig) Listener(name string) fiber.ListenConfig {
	listener := this.config.Listener(name)
	if lns, ok := this.listenerCustomizer[name]; ok {
		for _, ln := range lns {
			ln.FiberListenConfig(&listener)
		}
	}
	return listener
}
