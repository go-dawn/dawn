package redis

import (
	"testing"

	"github.com/go-dawn/dawn/config"

	"github.com/go-redis/redis/v8"

	"github.com/stretchr/testify/assert"
)

func Test_Redis_Module_Name(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "dawn:redis", New().String())
}

func Test_Redis_Module_Init(t *testing.T) {
	t.Parallel()

	config.Load("./", "redis")
	config.Set("Redis.Connections.Default.Addr", "127.0.0.1:99999")
	m := &Module{}

	m.Init()()
}

func Test_Redis_Module_Boot(t *testing.T) {
	m := &Module{
		conns: map[string]*redis.Client{
			fallback: redis.NewClient(&redis.Options{
				Addr: "127.0.0.1:99999",
			}),
		},
	}

	assert.Panics(t, m.Boot)
}

func Test_Redis_Cleanup(t *testing.T) {
	m := &Module{
		conns: map[string]*redis.Client{
			fallback: redis.NewClient(&redis.Options{}),
		},
	}

	m.cleanup()
}

func Test_Redis_Conn(t *testing.T) {
	assert.Nil(t, Conn("non"))
}
