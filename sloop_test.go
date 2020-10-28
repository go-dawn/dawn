package dawn

import (
	"os"
	"testing"
	"time"

	"github.com/go-dawn/dawn/config"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	m = mockModule{}
)

func Test_Sloop_New(t *testing.T) {
	t.Parallel()

	app := fiber.New()

	s := New(Config{App: app}).AddModulers(m)

	assert.Equal(t, app, s.app)
	assert.Len(t, s.mods, 1)
	assert.Equal(t, "anonymous", s.mods[0].String())
}

func Test_Sloop_Default(t *testing.T) {
	t.Parallel()

	config.Set("debug", true)
	s := Default(fiber.Config{})

	require.NotNil(t, s.app)
	assert.Len(t, s.app.Stack()[0], 1)
}

func Test_Sloop_AddModulers(t *testing.T) {
	t.Parallel()

	s := New().AddModulers(m)

	assert.Len(t, s.mods, 1)
	assert.Equal(t, "anonymous", s.mods[0].String())
}

func Test_Sloop_Run(t *testing.T) {
	t.Parallel()

	assert.NotNil(t, New().Run(""))

	s := New(Config{App: fiber.New()}).AddModulers(m)

	go func() {
		time.Sleep(time.Millisecond * 100)
		assert.NoError(t, s.app.Shutdown())
	}()

	assert.NoError(t, s.Run(""))
}

func Test_Sloop_RunTls(t *testing.T) {
	assert.NotNil(t, New().RunTls("", "", ""))

	s := New(Config{App: fiber.New()})

	t.Run("invalid addr", func(t *testing.T) {
		t.Parallel()

		assert.NotNil(t, s.RunTls(":99999", "./.github/testdata/ssl.pem", "./.github/testdata/ssl.key"))
	})

	t.Run("invalid ssl info", func(t *testing.T) {
		t.Parallel()

		assert.NotNil(t, s.RunTls("", "./.github/README.md", "./.github/README.md"))
	})

	t.Run("with ssl", func(t *testing.T) {
		t.Parallel()

		go func() {
			time.Sleep(time.Millisecond * 100)
			assert.NoError(t, s.app.Shutdown())
		}()

		assert.NoError(t, s.RunTls("", "./.github/testdata/ssl.pem", "./.github/testdata/ssl.key"))
	})
}

func Test_Sloop_Shutdown(t *testing.T) {
	t.Parallel()

	require.NotNil(t, (&Sloop{}).Shutdown())
	require.Nil(t, New(Config{App: fiber.New()}).Shutdown())
}

func Test_Sloop_Router(t *testing.T) {
	t.Parallel()

	require.Nil(t, (&Sloop{}).Router())
	require.NotNil(t, New(Config{App: fiber.New()}).Router())
}

func Test_Sloop_Watch(t *testing.T) {
	t.Parallel()

	s := &Sloop{
		sigCh: make(chan os.Signal, 1),
	}

	go s.Watch()

	select {
	case s.sigCh <- os.Interrupt:
	case <-time.NewTimer(time.Second).C:
		assert.Fail(t, "should receive signal")
	}
}
