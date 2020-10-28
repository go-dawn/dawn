package redis

import (
	"testing"

	"github.com/go-dawn/dawn/config"
	"github.com/stretchr/testify/assert"
)

func Test_Redis_New(t *testing.T) {
	t.Parallel()

	moduler := New()
	_, ok := moduler.(*redisModule)
	assert.True(t, ok)
}

func Test_Redis_Module_Name(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "dawn:redis", m.String())
}

func Test_Redis_Init(t *testing.T) {
	t.Run("empty config", func(t *testing.T) {
		m.Init()()
	})

	t.Run("error", func(t *testing.T) {
		defer func() {
			assert.Contains(t, recover(), "dawn:redis failed to ping")
		}()
		config.Load("./", "redis")
		config.Set("Redis.Connections.Default.Addr", "127.0.0.1:99999")

		m.Init()()
	})
}

func Test_Redis_Conn(t *testing.T) {
	assert.Nil(t, Conn("non"))
}
