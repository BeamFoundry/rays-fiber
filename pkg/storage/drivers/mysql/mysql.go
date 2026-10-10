package mysql

import (
	_ "github.com/BeamFoundry/rays-fiber/pkg/storage/drivers/mysql/internal"

	"github.com/gofiber/storage/mysql/v2"
)

// MySQLConfigOverride is an interface for a Component must implements to adapt the [mysql.Config]
// that will be created.
//
// The instance name (like "primary") is provided with a pointer to a config object created.
type MySQLConfigOverride interface {
	MySQLConfigOverride(instance string, config *mysql.Config)
}
