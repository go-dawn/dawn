package schema

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/jinzhu/copier"
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

func TestDaoToDto(t *testing.T) {
	type Dto struct {
		CreatedAt JsonTime
	}

	type Dao struct {
		CreatedAt time.Time
	}

	value := "2020-01-02 03:04:05"
	target, err := time.Parse(JsonTimeFormat, value)
	assert.Nil(t, err)

	a := Dao{target}
	var b Dto

	assert.Nil(t, copier.Copy(&b, &a))

	ret, err := json.Marshal(&b)
	assert.Nil(t, err)
	assert.Equal(t, `{"CreatedAt":"2020-01-02 03:04:05"}`, string(ret))
}
