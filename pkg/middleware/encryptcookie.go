package middleware

import (
	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/lang"

	. "github.com/BeamFoundry/rays-fiber/pkg/middleware/internal"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/encryptcookie"
)

// EncryptCookieNext is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [EncryptCookie], to indicates when to skip it this middleware.
//
// Please see the [encryptcookie] middleware configuration documentation.
//
// [encryptcookie]: https://docs.gofiber.io/middleware/encryptcookie/#config
type EncryptCookieNext interface {
	EncryptCookieNext() func(fiber.Ctx) bool
}

// EncryptCookieNext is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [EncryptCookie], to indicates a custom function to encrypt cookies.
//
// Please see the [encryptcookie] middleware configuration documentation.
//
// [encryptcookie]: https://docs.gofiber.io/middleware/encryptcookie/#config
type EncryptCookieEncryptor interface {
	EncryptCookieEncryptor() func(name, decryptedString, key string) (string, error)
}

// EncryptCookieNext is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [EncryptCookie], to indicates a custom function to decrypt cookies.
//
// Please see the [encryptcookie] middleware configuration documentation.
//
// [encryptcookie]: https://docs.gofiber.io/middleware/encryptcookie/#config
type EncryptCookieDecryptor interface {
	EncryptCookieDecryptor() func(name, encryptedString, key string) (string, error)
}

// EncryptCookieConfig is an interface a [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController]
// or a [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements,
// if it uses [EncryptCookie], to review the [encryptcookie.Config] that will be configured by [EncryptCookie]
//
// Please see the [encryptcookie] middleware configuration documentation.
//
// [encryptcookie]: https://docs.gofiber.io/middleware/encryptcookie/#config
type EncryptCookieConfig interface {
	EncryptCookieConfig(config *encryptcookie.Config)
}

type encryptCookieConfigurationProperties struct {
	key    *string
	except []string
}

type encryptCookieConfig struct {
	encryptcookie encryptCookieConfigurationProperties
}

var _, _ = any(&FiberMwConfig[encryptCookieConfig]{}).(Stereotype)

func (this *encryptCookieConfig) Config() encryptcookie.Config {
	config := encryptcookie.ConfigDefault

	if this.encryptcookie.key != nil {
		config.Key = *this.encryptcookie.key
	} else {
		config.Key = encryptcookie.GenerateKey(32)
	}

	config.Except = this.encryptcookie.except

	return config
}

// [EncryptCookie] is a middleware to be defined directly as an anonymous member inside a
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter].
//
// # Example
//
//	type ApiController struct {
//	  FiberController
//	  RequestMapping `Path:"/api"`
//	  EncryptCookie
//	}
//
//	func (this *ApiController) EncryptCookieConfig(config *encryptcookie.Config) {
//	  // changing configuration for this controller.
//	}
//
// # Configuration
//
// Its configuration is done through Rays configuration properties and through interfaces the parent
// [github.com/BeamFoundry/rays-fiber/pkg/core.FiberController] or
// [github.com/BeamFoundry/rays-fiber/pkg/core.FilterPathFilter] may implements.
// Those interfaces are [EncryptCookieNext], [EncryptCookieEncryptor] and [EncryptCookieDecryptor].
//
// Also, the parent container may implements [EncryptCookieConfig] to receive a pointer to the [encryptcookie.Config]
// that will be used.
//
// Please see those interface definitions and the [encryptcookie] middleware documentation.
//
// # Properties mapping :
//
//	rays:
//		fiber:
//			servers:
//				primary:
//					encryptcookie:
//	 					key: null
//	 					except: []
//
// [encryptcookie]: https://docs.gofiber.io/middleware/encryptcookie/#config
type EncryptCookie struct {
	config *FiberMwConfig[encryptCookieConfig] `@:"Inject"`
	cfg    *encryptcookie.Config
}

// PostConstruct() implements [github.com/BeamFoundry/rays/pkg/core.InnerBeamPostConstruct]
func (this *EncryptCookie) PostConstruct(parent any) Error {
	this.cfg = new(this.config.GetConfig(parent).Config())
	if next, ok := (parent).(EncryptCookieNext); ok {
		this.cfg.Next = next.EncryptCookieNext()
	}
	if encryptor, ok := (parent).(EncryptCookieEncryptor); ok {
		this.cfg.Encryptor = encryptor.EncryptCookieEncryptor()
	}
	if decryptor, ok := (parent).(EncryptCookieDecryptor); ok {
		this.cfg.Decryptor = decryptor.EncryptCookieDecryptor()
	}
	if config, ok := (parent).(EncryptCookieConfig); ok {
		config.EncryptCookieConfig(this.cfg)
	}
	return nil
}

// Handler() implements [github.com/BeamFoundry/rays-fiber/pkg/core.FiberHandler]
func (this *EncryptCookie) Handler() fiber.Handler {
	return encryptcookie.New(*this.cfg)
}
