package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	configPath = "./"
	configName = "foo"
	key        = "foo"
	value      = "bar"

	nonExistKey = "non"

	mergeCfg = map[string]interface{}{
		"merge": "cfg",
	}

	fu  forUnmarshal
	sub SubConfig
)

type (
	forUnmarshal struct {
		S string
		SubConfig
	}
	SubConfig struct {
		B bool
	}
)

func Test_Load_Panic(t *testing.T) {
	t.Parallel()

	assert.Panics(t, func() {
		nonConfigName := "non config name"
		Load(configPath, nonConfigName)
	})
}

func Test_Load(t *testing.T) {
	reset()

	Load(configPath, configName)
	assert.Equal(t, value, GetString(key))

	defaultFileStorePath := "./data/dawn_store.db"
	assert.Equal(t, defaultFileStorePath, GetString("cache.file.path", defaultFileStorePath))
}

func Test_AllGetFunctions(t *testing.T) {
	Load(configPath, configName)

	assert.Equal(t, "iface", Get("iface"))
	assert.Equal(t, "di", Get(nonExistKey, "di"))

	assert.Equal(t, "iface", GetValue("iface"))
	assert.Equal(t, "di", GetValue(nonExistKey, "di"))

	assert.Equal(t, "s", GetString("string"))
	assert.Equal(t, "ds", GetString(nonExistKey, "ds"))

	assert.Equal(t, true, GetBool("Bool"))
	assert.Equal(t, true, GetBool(nonExistKey, true))

	assert.Equal(t, time.Second, GetDuration("Duration"))
	assert.Equal(t, time.Minute, GetDuration(nonExistKey, time.Minute))

	Time, _ := time.Parse("2006-01-02 15:04:05", "2020-03-07 12:31:19")
	assert.Equal(t, Time, GetTime("Time"))
	now := time.Now()
	assert.Equal(t, now, GetTime(nonExistKey, now))

	assert.Equal(t, 1, GetInt("Int"))
	assert.Equal(t, 2, GetInt(nonExistKey, 2))

	assert.Equal(t, int64(1), GetInt64("Int"))
	assert.Equal(t, int64(2), GetInt64(nonExistKey, 2))

	assert.Equal(t, 1.1, GetFloat64("Float64"))
	assert.Equal(t, 2.2, GetFloat64(nonExistKey, 2.2))

	assert.Equal(t, map[string]interface{}{"string": "Map"}, GetStringMap("StringMap"))
	assert.Equal(t, map[string]interface{}{"k1": "v1"},
		GetStringMap(nonExistKey, map[string]interface{}{"K1": "v1"}))

	assert.Equal(t, map[string]string{"string": "String"},
		GetStringMapString("StringMapString"))
	assert.Equal(t, map[string]string{"K1": "v1"},
		GetStringMapString(nonExistKey, map[string]string{"K1": "v1"}))

	assert.Equal(t, []string{"s1", "s2"}, GetStringSlice("StringSlice"))
	assert.Equal(t, []string{"s3", "s4"}, GetStringSlice(nonExistKey, []string{"s3", "s4"}))
}

func Test_AllSettings(t *testing.T) {
	reset()
	assert.Len(t, AllSettings(), 0)
}

func Test_Unmarshal(t *testing.T) {
	reset()

	Set("S", value)
	err := Unmarshal(&fu)
	assert.Nil(t, err)
	assert.Equal(t, value, fu.S)
}

func Test_UnmarshalKey(t *testing.T) {
	reset()

	Set("SubConfig.B", true)
	err := UnmarshalKey("SubConfig", &sub)
	assert.Nil(t, err)
	assert.True(t, sub.B)
}

func Test_MergeConfigMap(t *testing.T) {
	reset()

	MergeConfigMap(mergeCfg)
	assert.Equal(t, mergeCfg["merge"], GetString("merge"))
}

func Test_Sub(t *testing.T) {
	reset()

	Set("SubConfig.B", true)
	c := Sub("SubConfig")
	assert.True(t, c.GetBool("B"))
}

func Test_Has(t *testing.T) {
	reset()

	Set("SubConfig.B", true)

	assert.True(t, Has("SubConfig"))
	assert.True(t, Has("SubConfig.B"))
	assert.False(t, Has("B"))
}

func Test_LoadAll(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		assert.NotNil(t, LoadAll("./testdata/error"))
	})

	t.Run("success", func(t *testing.T) {
		assert.Nil(t, LoadAll("./testdata/all"))
		assert.True(t, global.Has("http"))
		assert.True(t, global.Has("others.1"))
	})

	t.Run("env", func(t *testing.T) {
		assert.Nil(t, LoadAll("./testdata/all"))
		assert.Equal(t, false, global.GetBool("app.debug"))

		LoadEnv("DAWN")

		require.NoError(t, os.Setenv("DAWN_APP_DEBUG", "true"))

		assert.Equal(t, true, global.GetBool("app.debug"))
	})
}

func reset() {
	global = New()
}
