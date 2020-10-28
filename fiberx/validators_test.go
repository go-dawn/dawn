package fiberx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMobile(t *testing.T) {
	t.Parallel()

	at := assert.New(t)

	at.Nil(V.Var("13888888888", "mobile"))
	at.NotNil(V.Var("23888888888", "mobile"))
	at.NotNil(V.Var("1388888888", "mobile"))
	at.NotNil(V.Var(13888888888, "mobile"))
}

func BenchmarkName(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = V.Var("13888888888", "mobile")
		}
	})
}
