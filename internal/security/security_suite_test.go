package security_test

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"testing"

	_ "github.com/BeamFoundry/rays"
	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/i18n"
	. "github.com/BeamFoundry/rays/pkg/i18n/locales"
	. "github.com/BeamFoundry/rays/pkg/lang"
	. "github.com/BeamFoundry/rays/pkg/testing"
	. "github.com/BeamFoundry/rays/pkg/web"

	_ "github.com/BeamFoundry/rays-fiber"
	. "github.com/BeamFoundry/rays-fiber/internal/security"
	. "github.com/BeamFoundry/rays-fiber/pkg/core"
	. "github.com/BeamFoundry/rays-fiber/pkg/middleware"
	. "github.com/BeamFoundry/rays-fiber/pkg/testing"

	_ "github.com/BeamFoundry/rays-web-security"
	. "github.com/BeamFoundry/rays-web-security/pkg/config"
	. "github.com/BeamFoundry/rays-web-security/pkg/core"
	. "github.com/BeamFoundry/rays-web-security/pkg/rememberme"
	. "github.com/BeamFoundry/rays-web-security/pkg/testing"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/gofiber/fiber/v3/middleware/session"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSecurity(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rays-Fiber/Internal/Security Suite")
}

//---

type TestController struct {
	Controller
}

func (this *TestController) Home(
	_ struct {
		RequestMapping `Path:"/"`
	},
	c fiber.Ctx,
) error {
	return c.SendString("Home")
}

func (this *TestController) User(
	_ struct {
		RequestMapping `Path:"/user"`
	},
	c fiber.Ctx,
) error {
	return c.SendString("User")
}

func (this *TestController) Admin(
	_ struct {
		RequestMapping `Path:"/admin"`
	},
	c fiber.Ctx,
) error {
	return c.SendString("Admin")
}

func (this *TestController) APIGet(
	_ struct {
		RequestMapping `Path:"/api/get"`
	},
	c fiber.Ctx,
) error {
	return c.SendString("API GET")
}

func (this *TestController) APIPost(
	_ struct {
		RequestMapping `Path:"/api/post" Method:"POST"`
	},
	c fiber.Ctx,
) error {
	return c.SendString("API POST")
}

func (this *TestController) APIPostNoCsrf(
	_ struct {
		RequestMapping `Path:"/api/nocsrf" Method:"POST"`
	},
	c fiber.Ctx,
) error {
	return c.SendString("API POST")
}

// ---
type LocaleConfig struct {
	Configuration
	defaultLocale          Locale
	localeContextHolder    *ThreadLocaleContextHolder
	messageSource          *DefaultMessageSource
	messageSourceValidator *DefaultMessageSourceValidator
}

func (this LocaleConfig) New() *LocaleConfig {
	this.defaultLocale = Locales.EN

	this.localeContextHolder = ThreadLocaleContextHolder{}.New(this.defaultLocale)

	this.messageSource = DefaultMessageSource{}.New(this.defaultLocale)
	this.messageSourceValidator = DefaultMessageSourceValidator{}.New(this.localeContextHolder, this.messageSource)

	return &this
}

func (this *LocaleConfig) PostConstruct() Error {
	bundle := this.messageSource.DefaultMessageBundle()
	if err := bundle.LoadValidationTranslation(ValidationTranslations.EN); err != nil {
		return NewErrorFrom(err)
	}
	return nil
}

func (this *LocaleConfig) LocaleContextHolder() *ThreadLocaleContextHolder {
	return this.localeContextHolder
}

func (this *LocaleConfig) MessageSource() *DefaultMessageSource {
	return this.messageSource
}

func (this *LocaleConfig) MessageSourceValidator() *DefaultMessageSourceValidator {
	return this.messageSourceValidator
}

//--

type TestSecurityConfiguration struct {
	Configuration
	encoder             PasswordEncoder                   `@:"Inject"`
	securityMiddlewares Foundry[*FiberSecurityAppHandler] `@:"Inject"`
}

func (this *TestSecurityConfiguration) PostConstruct() Error {
	ctx := context.Background()
	ctx = context.WithValue(ctx, PrototypeContextProxyType, this)
	_, err := this.securityMiddlewares.Forge(ctx)
	Expect(err).To(HaveOccurred())
	return nil
}

func (this *TestSecurityConfiguration) UserDetailsService() UserDetailsService {
	userDetailsService := InMemoryUserDetailsManager{}.New(this.encoder)
	userDetailsService.NewUser().
		WithUsername("user").
		WithPassword("user").
		WithAuthorities([]GrantedAuthority{Authority("ROLE_USER")}).Build()
	userDetailsService.NewUser().
		WithUsername("admin").
		WithPassword("admin").
		WithAuthorities([]GrantedAuthority{Authority("ROLE_ADMIN")}).Build()
	return userDetailsService
}

func (this *TestSecurityConfiguration) BasicAuth() *HttpSecurity {
	return (&HttpSecurityBuilder{}).
		Match("/api/get", "/api/post").
		AnyRequest().HasAuthority("ROLE_ADMIN").
		Basic().
		Build()
}

func (this *TestSecurityConfiguration) BasicAuthNoCsrf() *HttpSecurity {
	return (&HttpSecurityBuilder{}).
		Match("/api/nocsrf").
		AnyRequest().HasAuthority("ROLE_ADMIN").
		Csrf().Disable().
		Basic().
		Build()
}

func (this *TestSecurityConfiguration) HttpSecurity(
	that struct {
		rememberme RememberMeService `@:"Inject"`
	},
) *HttpSecurity {
	return (&HttpSecurityBuilder{}).
		AuthorizedHTTPRequests().
		Paths("/", "/login", "/logout").PermitAll().
		Paths("/user").HasAuthority("ROLE_USER").
		Paths("/admin").HasAuthority("ROLE_ADMIN").
		And().AnyRequest().DenyAll().
		FormLogin().ViewName("testdata/views/login").
		And().RememberMe(that.rememberme).
		Build()
}

func (this *TestSecurityConfiguration) RememberMe(
	that struct {
		userDetailsService UserDetailsService `@:"Inject"`
	},
) RememberMeService {
	return TokenBasedRememberMeService{}.New("secret", that.userDetailsService)
}

type TestSecurityApphandlerMdw struct {
	Component
}

// SessionSecurityConfig implements [middleware.SessionSecurityConfig].
func (t *TestSecurityApphandlerMdw) SessionSecurityConfig(_ string, _ *session.Config) {
}

// CsrfSecurityConfig implements [middleware.CsrfSecurityConfig].
func (t *TestSecurityApphandlerMdw) CsrfSecurityConfig(_ string, _ *csrf.Config) {
}

var _ CsrfSecurityConfig = &TestSecurityApphandlerMdw{}

var _ SessionSecurityConfig = &TestSecurityApphandlerMdw{}

var _ = Describe("FiberSecurityAppHandler", func() {
	Context("Without Locale", func() {
		setAuth := func(req *http.Request) {
			req.SetBasicAuth("admin", "admin")
		}
		BeforeEach(func() {
			AddBeamFor(&TestSecurityAppConfigCustomizer{})
			AddBeamFor(&TestSecurityAppHandlers{})
			AddBeamFor(&TestSecurityPathAppHandlers{})
			AddBeamFor(&TestSecurityConfiguration{})
			AddBeamFor(&TestController{})
			AddBeamFor(&TestSecurityApphandlerMdw{})
			GetContainer().EnableScanning()
			err := GetContainer().Initialize()
			Expect(err).NotTo(HaveOccurred())
		})
		It("can get home without login", func() {
			appProvider := GetBeamScopedInstance[FiberAppProvider]()
			client := TestClient{}.New(appProvider.Default())
			resp := client.GET("/")
			Expect(resp.StatusCode).To(Equal(fiber.StatusOK))
			body, err := io.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(body)).To(Equal("Home"))
		})
		It("can get the login page", func() {
			appProvider := GetBeamScopedInstance[FiberAppProvider]()
			client := TestClient{}.New(appProvider.Default())
			resp := client.GET("/login")
			Expect(resp.StatusCode).To(Equal(fiber.StatusOK))
		})
		It("can get /api/get with basic login", func() {
			appProvider := GetBeamScopedInstance[FiberAppProvider]()
			client := TestClient{}.New(appProvider.Default())
			resp := client.GET("/api/get")
			Expect(resp.StatusCode).To(Equal(fiber.StatusUnauthorized))

			resp = client.GET("/api/get", setAuth)
			Expect(resp.StatusCode).To(Equal(fiber.StatusOK))
		})
		It("can post /api/post with basic login & csrf", func() {
			appProvider := GetBeamScopedInstance[FiberAppProvider]()
			client := TestClient{}.New(appProvider.Default())
			resp := client.GET("/api/get")
			Expect(resp.StatusCode).To(Equal(fiber.StatusUnauthorized))

			form := url.Values{}
			form.Set("data", "data")

			resp = client.POST("/api/post", form)
			Expect(resp.StatusCode).To(Equal(fiber.StatusForbidden))

			resp = client.POST("/api/post", form, setAuth)
			Expect(resp.StatusCode).To(Equal(fiber.StatusForbidden))

			setCsrf := func(req *http.Request) {
				req.Header.Set("X-Csrf-Token", client.GetCookieValue(CsrfCookieName))
			}
			resp = client.POST("/api/post", form, setAuth, setCsrf)
			Expect(resp.StatusCode).To(Equal(fiber.StatusOK)) // FIXME: ???
		})
		It("can post /api/nocsrf with basic login", func() {
			appProvider := GetBeamScopedInstance[FiberAppProvider]()
			client := TestClient{}.New(appProvider.Default())
			form := url.Values{}
			form.Set("data", "data")
			resp := client.POST("/api/nocsrf", form)
			Expect(resp.StatusCode).To(Equal(fiber.StatusUnauthorized))

			resp = client.POST("/api/nocsrf", form, setAuth)
			Expect(resp.StatusCode).To(Equal(fiber.StatusOK))
		})
		It("can not get /user without login", func() {
			appProvider := GetBeamScopedInstance[FiberAppProvider]()
			client := TestClient{}.New(appProvider.Default())

			resp := client.GET("/user")
			Expect(resp.StatusCode).To(Equal(fiber.StatusFound))
			Expect(resp.Header.Get("Location")).To(Equal("/login"))

			resp = client.GET("/login")
			Expect(resp.StatusCode).To(Equal(fiber.StatusOK))

			form := url.Values{}
			form.Set("username", "user")
			form.Set("_csrf", client.GetCookieValue(CsrfCookieName))
			form.Set("remember-me", "true")
			resp = client.POST("/login", form)
			Expect(resp.StatusCode).To(Equal(fiber.StatusOK)) // FIXME: ???

			form.Set("password", "bad")
			resp = client.POST("/login", form)
			Expect(resp.StatusCode).To(Equal(fiber.StatusOK)) // FIXME: ???

			form.Set("password", "user")
			resp = client.POST("/login", form)
			Expect(resp.StatusCode).To(Equal(fiber.StatusFound)) // FIXME: ???
			Expect(resp.Header.Get("Location")).To(Equal("/user"))

			resp = client.GET("/user")
			Expect(resp.StatusCode).To(Equal(fiber.StatusOK))
			body, err := io.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(body)).To(Equal("User"))

			// Logout
			resp = client.GET("/logout")
			Expect(resp.StatusCode).To(Equal(fiber.StatusFound))
			Expect(resp.Header.Get("Location")).To(Equal("/"))

			/*
				resp = client.GET("/user")
				Expect(resp.StatusCode).To(Equal(fiber.StatusOK))
				body, err = io.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(body)).To(Equal("User"))
			*/
		})
	})

	Context("With Locale", func() {
		BeforeEach(func() {
			AddBeamFor(&TestSecurityAppConfigCustomizer{})
			AddBeamFor(&TestSecurityAppHandlers{})
			AddBeamFor(&TestSecurityPathAppHandlers{})
			AddBeamFor(&TestSecurityConfiguration{})
			AddBeamFor(&TestController{})
			AddBeamFor(&LocaleConfig{})
			GetContainer().EnableScanning()
			err := GetContainer().Initialize()
			Expect(err).NotTo(HaveOccurred())
		})
		It("can not get /user without login", func() {
			appProvider := GetBeamScopedInstance[FiberAppProvider]()
			app := appProvider.Default()

			client := TestClient{}.New(app)

			resp := client.GET("/user")
			Expect(resp.StatusCode).To(Equal(fiber.StatusFound))
			Expect(resp.Header.Get("Location")).To(Equal("/login"))

			resp = client.GET("/login")
			Expect(resp.StatusCode).To(Equal(fiber.StatusOK))

			form := url.Values{}
			form.Set("username", "user")
			form.Set("remember-me", "true")
			resp = client.POST("/login", form)
			Expect(resp.StatusCode).To(Equal(fiber.StatusForbidden))

			resp = client.GET("/user")
			Expect(resp.StatusCode).To(Equal(fiber.StatusFound))
			Expect(resp.Header.Get("Location")).To(Equal("/login"))

			form.Set("_csrf", client.GetCookieValue(CsrfCookieName))
			resp = client.POST("/login", form)
			Expect(resp.StatusCode).To(Equal(fiber.StatusOK)) // FIXME: ???

			form.Set("password", "user")
			resp = client.POST("/login", form)
			Expect(resp.StatusCode).To(Equal(fiber.StatusFound)) // FIXME: ???
			Expect(resp.Header.Get("Location")).To(Equal("/user"))

			resp = client.GET("/user")
			Expect(resp.StatusCode).To(Equal(fiber.StatusOK))
			body, err := io.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(string(body)).To(Equal("User"))
		})
	})
})
