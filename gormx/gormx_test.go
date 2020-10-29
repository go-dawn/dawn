package gormx

import (
	"testing"

	"github.com/go-dawn/dawn/gormx/schema"
	"github.com/stretchr/testify/assert"

	"github.com/valyala/fasthttp"
)

func Test_Gormx_Paginate(t *testing.T) {
	at := assert.New(t)

	gdb := MockGdb(t)

	type data struct{}

	p, err := Paginate(gdb, schema.NewIndexQuery(fasthttp.AcquireArgs()), &Fake{}, &data{})

	at.Nil(err)
	at.Equal(1, p.Page)
	at.Equal(15, p.PageSize)
	at.Equal(0, p.Total)
	// TODO: result
	//at.Len(p.Data, 0)
}
