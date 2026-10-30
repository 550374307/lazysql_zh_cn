package components

import (
	"fmt"
	"strings"

	"github.com/jorgerojas26/lazysql/app"
	"github.com/jorgerojas26/lazysql/drivers"
	"github.com/jorgerojas26/lazysql/helpers"
	"github.com/jorgerojas26/lazysql/internal/i18n"
	"github.com/jorgerojas26/lazysql/models"
)

func InitFromArg(connectionString string, readOnly bool) error {
	parsed, err := helpers.ParseConnectionString(connectionString)
	if err != nil {
		return fmt.Errorf("%s: %s", i18n.T("could not parse connection string"), err)
	}
	DBName := strings.Split(parsed.Normalize(",", "NULL", 0), ",")[3]

	if DBName == "NULL" {
		DBName = ""
	}

	connection := models.Connection{
		Name:     "",
		Provider: parsed.Driver,
		DBName:   DBName,
		URL:      connectionString,
		ReadOnly: readOnly,
	}

	poolConfig, err := app.App.Config().EffectiveConnectionPool(connection)
	if err != nil {
		return fmt.Errorf("%s: %v", i18n.T("invalid connection pool configuration"), err)
	}

	var newDBDriver drivers.Driver
	switch connection.Provider {
	case drivers.DriverMySQL:
		newDBDriver = &drivers.MySQL{PoolConfig: poolConfig}
	case drivers.DriverPostgres:
		newDBDriver = &drivers.Postgres{PoolConfig: poolConfig}
	case drivers.DriverSqlite:
		newDBDriver = &drivers.SQLite{}
	case drivers.DriverMSSQL:
		newDBDriver = &drivers.MSSQL{PoolConfig: poolConfig}
	case drivers.DriverClickHouse:
		newDBDriver = &drivers.ClickHouse{PoolConfig: poolConfig}
	default:
		return fmt.Errorf("%s: %s", i18n.T("could not handle database driver"), connection.Provider)
	}

	err = newDBDriver.Connect(app.App.Context(), connection.URL)
	if err != nil {
		return fmt.Errorf("%s: %s", i18n.T("could not connect to database"), err)
	}
	initialTelemetryEngine = newDBDriver.GetProvider()
	initialTelemetryReadOnly = connection.ReadOnly
	mainPages.AddAndSwitchToPage(connection.URL, NewHomePage(connection, newDBDriver).Flex, true)

	return nil
}
