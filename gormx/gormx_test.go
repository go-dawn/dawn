package gormx

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func Test_Gormx_Paginate(t *testing.T) {
	at := assert.New(t)

	gdb := MockGdb(t)

	args := fasthttp.AcquireArgs()

	type Data struct {
		ID        uint32    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		F         string    `json:"f"`
	}

	args.Set("page", "1")
	p, err := Paginate(gdb, NewIndexQuery(args), &[]Fake{}, &[]Data{})
	at.NotNil(err)

	at.Nil(gdb.AutoMigrate(&Fake{}))
	fakers := []Fake{{F: "f0"}, {F: "f1"}, {F: "f2"}}
	at.Nil(gdb.Create(&fakers).Error)

	args.Reset()
	args.Set("pageSize", "2")

	p, err = Paginate(gdb, NewIndexQuery(args), &[]Fake{}, &[]Data{})

	at.Nil(err)
	at.Equal(1, p.Page)
	at.Equal(2, p.PageSize)
	at.Equal(3, p.Total)

	list, ok := p.Data.([]Data)
	at.Len(list, 2)
	at.True(ok)
	at.Equal("f0", list[0].F)
	at.Equal("f1", list[1].F)
}
