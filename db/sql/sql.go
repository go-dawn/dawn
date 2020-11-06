package sql

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/go-dawn/dawn"
	"github.com/go-dawn/dawn/config"
	"github.com/go-dawn/pkg/deck"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	m        = &Module{}
	fallback = "testing"
)

type Module struct {
	dawn.Module
	conns    map[string]*gorm.DB
	fallback string
}

// New gets the moduler
func New() *Module {
	return m
}

// String is module name
func (*Module) String() string {
	return "dawn:sql"
}

// Init does connection work to each database by config:
//  [Sql]
//  Default = "testing"
//  [Sql.Connections]
//  [Sql.Connections.testing]
//  Driver = "sqlite"
//  [Sql.Connections.mysql]
//  Driver = "mysql"
func (m *Module) Init() dawn.Cleanup {
	m.conns = make(map[string]*gorm.DB)

	// extract sql config
	c := config.Sub("sql")

	m.fallback = c.GetString("default", fallback)

	connsConfig := c.GetStringMap("connections")

	if len(connsConfig) == 0 {
		m.conns[m.fallback] = connect(m.fallback, config.New())
		return m.cleanup
	}

	// connect each db in config
	for name := range connsConfig {
		cfg := c.Sub("connections." + name)
		m.conns[name] = connect(name, cfg)
	}

	return m.cleanup
}

// cleanup 	close every connections
func (m *Module) cleanup() {
	for _, gdb := range m.conns {
		if db, err := gdb.DB(); err == nil {
			_ = db.Close()
		}
	}
}

func connect(name string, c *config.Config) (db *gorm.DB) {
	driver := c.GetString("driver", "sqlite")

	var err error
	switch strings.ToLower(driver) {
	case "sqlite":
		db, err = resolveSqlite(c)
	case "mysql":
		db, err = resolveMysql(c)
	case "postgres":
		db, err = resolvePostgres(c)
	default:
		panic(fmt.Sprintf("dawn:sql unknown driver %s of %s", driver, name))
	}

	if err != nil || db == nil {
		panic(fmt.Sprintf("dawn:sql failed to connect %s(%s): %v", name, driver, err))
	}

	return
}

// Conn gets sql connection by specific name or fallback
func Conn(name ...string) *gorm.DB {
	n := m.fallback

	if len(name) > 0 && name[0] != "" {
		n = name[0]
	}

	return m.conns[n]
}

var l = deck.DisabledGormLogger{}

// resolveSqlite resolves sqlite connection with config:
// Driver = "sqlite"
// Database = "file:dawn?mode=memory&cache=shared&_fk=1"
// Prefix = "dawn_"
// Log = false
func resolveSqlite(c *config.Config) (*gorm.DB, error) {
	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   c.GetString("Prefix"),
			SingularTable: false,
		},
	}

	// disable logger
	if !c.GetBool("Log") {
		gormConfig.Logger = l
	}

	dbname := c.GetString("Database", "file:dawn?mode=memory&cache=shared&_fk=1")

	return gorm.Open(sqlite.Open(dbname), gormConfig)
}

// resolveMysql resolves mysql connection with config:
//Driver = "mysql"
//Username = "username"
//Password = "password"
//Host = "127.0.0.1"
//Port = "3306"
//Database = "database"
//Location = "Asia/Shanghai"
//Charset = "utf8mb4"
//ParseTime = true
//Prefix = "dawn_"
//Log = false
//MaxIdleConns = 10
//MaxOpenConns = 100
//ConnMaxLifetime = "5m"
func resolveMysql(c *config.Config) (*gorm.DB, error) {
	parseTime := "True"
	if !c.GetBool("ParseTime", true) {
		parseTime = "False"
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		c.GetString("Username"),
		c.GetString("Password"),
		c.GetString("Host"),
		c.GetString("Port"),
		c.GetString("Database"),
		c.GetString("Charset", "utf8mb4"),
		parseTime,
		url.QueryEscape(c.GetString("Location", "Asia/Shanghai")),
	)

	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   c.GetString("Prefix"),
			SingularTable: false,
		},
	}

	// disable logger
	if !c.GetBool("Log") {
		gormConfig.Logger = l
	}

	gdb, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err == nil || c.GetBool("Testing") {
		db, err := gdb.DB()
		if err != nil {
			return gdb, err
		}
		db.SetMaxIdleConns(c.GetInt("MaxIdleConns"))
		db.SetMaxOpenConns(c.GetInt("MaxOpenConns"))
		db.SetConnMaxLifetime(c.GetDuration("ConnMaxLifetime"))
	}

	return gdb, err
}

// resolvePostgres resolves postgres connection with config:
//Driver = "postgres"
//Username = "username"
//Password = "password"
//Host = "127.0.0.1"
//Port = "5432"
//Database = "database"
//Sslmode = "disable"
//TimeZone = "Asia/Shanghai"
//Prefix = "dawn_"
//Log = false
//MaxIdleConns = 10
//MaxOpenConns = 100
//ConnMaxLifetime = "5m"
func resolvePostgres(c *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s TimeZone=%s",
		c.GetString("Username"),
		c.GetString("Password"),
		c.GetString("Host"),
		c.GetString("Port"),
		c.GetString("Database"),
		c.GetString("Sslmode", "disable"),
		url.QueryEscape(c.GetString("TimeZone", "Asia/Shanghai")),
	)

	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   c.GetString("Prefix"),
			SingularTable: false,
		},
	}

	// disable logger
	if !c.GetBool("Log") {
		gormConfig.Logger = l
	}

	gdb, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err == nil || c.GetBool("Testing") {
		db, err := gdb.DB()
		if err != nil {
			return gdb, err
		}
		db.SetMaxIdleConns(c.GetInt("MaxIdleConns"))
		db.SetMaxOpenConns(c.GetInt("MaxOpenConns"))
		db.SetConnMaxLifetime(c.GetDuration("ConnMaxLifetime"))
	}

	return gdb, err
}
