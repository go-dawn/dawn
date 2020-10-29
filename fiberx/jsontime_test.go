package fiberx

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJsonTime_MarshalJSON(t *testing.T) {
	t.Run("out of range", func(t *testing.T) {
		target := time.Unix(253402300800, 0)
		jt := JsonTime(target)

		_, err := jt.MarshalJSON()
		assert.NotNil(t, err)
	})

	t.Run("success", func(t *testing.T) {
		v := "2020-01-02 03:04:05"

		target, err := time.Parse(JsonTimeFormat, v)
		assert.Nil(t, err)

		jt := JsonTime(target)

		b, err := jt.MarshalJSON()
		assert.Nil(t, err)

		expect := `"` + v + `"`
		assert.Equal(t, expect, string(b))
	})
}

func TestJsonTime_UnmarshalJSON(t *testing.T) {
	t.Run("null", func(t *testing.T) {
		var (
			v  = "null"
			jt JsonTime
		)

		err := jt.UnmarshalJSON([]byte(v))

		assert.Nil(t, err)
		assert.Equal(t, JsonTime(time.Time{}), jt)
	})

	t.Run("illegal format", func(t *testing.T) {
		var (
			v  string
			jt JsonTime
		)

		err := jt.UnmarshalJSON([]byte(v))

		assert.NotNil(t, err)
	})

	t.Run("success", func(t *testing.T) {
		var (
			v  = `"2020-01-02 03:04:05"`
			jt JsonTime
		)

		err := jt.UnmarshalJSON([]byte(v))
		assert.Nil(t, err)

		assert.Contains(t, v, time.Time(jt).Format(JsonTimeFormat))
	})
}

func TestJsonTime_String(t *testing.T) {
	v := "2020-01-02 03:04:05"

	target, err := time.Parse(JsonTimeFormat, v)
	assert.Nil(t, err)

	jt := JsonTime(target)

	assert.Equal(t, v, jt.String())
}
