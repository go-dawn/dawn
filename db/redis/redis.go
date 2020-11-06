package redis

import (
	"context"
	"fmt"

	"github.com/go-dawn/dawn"
	"github.com/go-dawn/dawn/config"
	"github.com/go-redis/redis/v8"
)

var (
	m        = &Module{}
	fallback = "default"
)

type Module struct {
	dawn.Module
	conns    map[string]*redis.Client
	fallback string
}

// New gets the moduler
func New() *Module {
	return m
}

// String is module name
func (*Module) String() string {
	return "dawn:redis"
}

// Init does connection work to each database by config:
//  [Redis]
//  Default = "default"
//  [Redis.Connections]
//  [Redis.Connections.default]
//  Network = "tcp"
//  Addr = "127.0.0.1:6379"
//  Username = "username"
//  Password = "password"
//  DB = 0
//  MaxRetries = 5
//  DialTimeout = "5s"
//  ReadTimeout = "5s"
//  WriteTimeout = "5s"
//  PoolSize = 1024
//  MinIdleConns = 10
//  MaxConnAge = "1m"
//  PoolTimeout = "1m"
//  IdleTimeout = "1m"
//  IdleCheckFrequency = "1m"
func (m *Module) Init() dawn.Cleanup {
	m.conns = make(map[string]*redis.Client)

	// extract redis config
	c := config.Sub("redis")

	m.fallback = c.GetString("default", fallback)

	connsConfig := c.GetStringMap("connections")

	// connect each db in config
	for name := range connsConfig {
		cfg := c.Sub("connections." + name)
		m.conns[name] = connect(name, cfg)
	}

	return m.cleanup
}

func connect(name string, c *config.Config) (client *redis.Client) {
	addr := c.GetString("Addr", "127.0.0.1:6379")
	client = redis.NewClient(&redis.Options{
		Network:            c.GetString("Network"),
		Addr:               addr,
		Username:           c.GetString("Username"),
		Password:           c.GetString("Password"),
		DB:                 c.GetInt("DB"),
		MaxRetries:         c.GetInt("MaxRetries"),
		DialTimeout:        c.GetDuration("DialTimeout"),
		ReadTimeout:        c.GetDuration("ReadTimeout"),
		WriteTimeout:       c.GetDuration("WriteTimeout"),
		PoolSize:           c.GetInt("PoolSize"),
		MinIdleConns:       c.GetInt("MinIdleConns"),
		MaxConnAge:         c.GetDuration("MaxConnAge"),
		PoolTimeout:        c.GetDuration("PoolTimeout"),
		IdleTimeout:        c.GetDuration("IdleTimeout"),
		IdleCheckFrequency: c.GetDuration("IdleCheckFrequency"),
	})

	return
}

func (m *Module) Boot() {
	for name, client := range m.conns {
		if _, err := client.Ping(context.Background()).Result(); err != nil {
			panic(fmt.Sprintf("dawn:redis failed to ping %s(%s): %v", name, client.Options().Addr, err))
		}
	}
}

func (m *Module) cleanup() {
	// close every connections
	for _, client := range m.conns {
		_ = client.Close()
	}
}

// Conn gets redis connection by specific name or fallback
func Conn(name ...string) redis.Cmdable {
	n := m.fallback

	if len(name) > 0 && name[0] != "" {
		n = name[0]
	}

	return m.conns[n]
}
