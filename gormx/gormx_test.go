package gormx

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Test_Gormx_Paginate(t *testing.T) {
	at := assert.New(t)

	gdb := mockGdb(t)

	args := fasthttp.AcquireArgs()

	type Data struct {
		ID        uint32    `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		F         string    `json:"f"`
	}

	args.Set("page", "1")
	p, err := Paginate(gdb, args, &[]Fake{}, &[]Data{})
	at.NotNil(err)

	at.Nil(gdb.AutoMigrate(&Fake{}))
	fakers := []Fake{{F: "f0"}, {F: "f1"}, {F: "f2"}}
	at.Nil(gdb.Create(&fakers).Error)

	args.Reset()
	args.Set("pageSize", "2")
	p, err = Paginate(gdb, args, &[]Fake{}, &[]Data{})

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

func Test_Scope_Paginate(t *testing.T) {
	gdb := dryRunSession(t)

	t.Run("int", func(t *testing.T) {
		stat := gdb.Scopes(scopePaginate(1, 10)).Find(&Fake{}).Statement

		assert.Equal(t, "SELECT * FROM `fakes` WHERE `fakes`.`deleted_at` IS NULL LIMIT 10", stat.SQL.String())
	})

	t.Run("offset", func(t *testing.T) {
		stat := gdb.Scopes(scopePaginate(2, 10)).Find(&Fake{}).Statement

		assert.Equal(t, "SELECT * FROM `fakes` WHERE `fakes`.`deleted_at` IS NULL LIMIT 10 OFFSET 10", stat.SQL.String())
	})
}

func Test_Scope_Search(t *testing.T) {
	gdb := dryRunSession(t)

	t.Run("empty object", func(t *testing.T) {
		stat := gdb.Scopes(scopeSearch([]byte(`{}`))).Find(&Fake{}).Statement

		assert.Equal(t, "SELECT * FROM `fakes` WHERE `fakes`.`deleted_at` IS NULL", stat.SQL.String())
	})

	t.Run("number", func(t *testing.T) {
		stat := gdb.Scopes(scopeSearch([]byte(`{"id":1}`))).Find(&Fake{}).Statement

		assert.Equal(t, "SELECT * FROM `fakes` WHERE id = ? AND `fakes`.`deleted_at` IS NULL", stat.SQL.String())
	})

	t.Run("like", func(t *testing.T) {
		stat := gdb.Scopes(scopeSearch([]byte(`{"name":"k"}`))).Find(&Fake{}).Statement

		assert.Equal(t, "SELECT * FROM `fakes` WHERE name LIKE ? AND `fakes`.`deleted_at` IS NULL", stat.SQL.String())
		assert.Equal(t, []interface{}{"%k%"}, stat.Vars)
	})

	t.Run("in", func(t *testing.T) {
		t.Run("number", func(t *testing.T) {
			stat := gdb.Scopes(scopeSearch([]byte(`{"name":[1.1,2.2,3.3]}`))).
				Find(&Fake{}).Statement

			assert.Equal(t, "SELECT * FROM `fakes` WHERE name IN (?,?,?) AND `fakes`.`deleted_at` IS NULL", stat.SQL.String())
			assert.Equal(t, []interface{}{1.1, 2.2, 3.3}, stat.Vars)
		})
		t.Run("string", func(t *testing.T) {
			stat := gdb.Scopes(scopeSearch([]byte(`{"name":["1","2","3"]}`))).
				Find(&Fake{}).Statement
			assert.Equal(t, "SELECT * FROM `fakes` WHERE name IN (?,?,?) AND `fakes`.`deleted_at` IS NULL", stat.SQL.String())
			assert.Equal(t, []interface{}{"1", "2", "3"}, stat.Vars)
		})
	})

	t.Run("operator", func(t *testing.T) {
		for _, opt := range []string{"<", "<=", ">", ">="} {
			for _, val := range []interface{}{1.1, "3"} {
				name := fmt.Sprintf("%s %v", opt, val)
				t.Run(name, func(t *testing.T) {
					var search string
					search = fmt.Sprintf(`{"c$%s":%v}`, opt, val)
					if val == "3" {
						search = fmt.Sprintf(`{"c$%s":"%s"}`, opt, val)
					}
					stat := gdb.Scopes(scopeSearch([]byte(search))).Find(&Fake{}).Statement

					exp := fmt.Sprintf("SELECT * FROM `fakes` WHERE c %s ? AND `fakes`.`deleted_at` IS NULL", opt)
					assert.Equal(t, exp, stat.SQL.String())
					assert.Equal(t, []interface{}{val}, stat.Vars)
				})
			}
		}

		t.Run("><", func(t *testing.T) {
			stat := gdb.Scopes(scopeSearch([]byte(`{"c$><":["2020-01-01", "2020-03-01"]}`))).Find(&Fake{}).Statement

			assert.Equal(t, "SELECT * FROM `fakes` WHERE (c BETWEEN ? AND ?) AND `fakes`.`deleted_at` IS NULL", stat.SQL.String())
			assert.Equal(t, []interface{}{"2020-01-01", "2020-03-01"}, stat.Vars)
		})

		t.Run("bool", func(t *testing.T) {
			stat := gdb.Scopes(scopeSearch([]byte(`{"ok":true}`))).Find(&Fake{}).Statement

			assert.Equal(t, "SELECT * FROM `fakes` WHERE ok = ? AND `fakes`.`deleted_at` IS NULL", stat.SQL.String())
			assert.Equal(t, []interface{}{true}, stat.Vars)
		})
	})
}

func Benchmark_Scope_Search(b *testing.B) {
	gdb := dryRunSession(b)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		gdb.Session(&gorm.Session{DryRun: false}).
			Scopes(scopeSearch([]byte(`{"ID":1,"created_at$><":["2020-01-01", "2020-03-01"]}`))).
			Find(&Fake{})
	}
}

func Test_Scope_Sort(t *testing.T) {
	gdb := dryRunSession(t)

	t.Run("no sort", func(t *testing.T) {
		stat := gdb.Scopes(scopeSort([]byte(""))).Find(&Fake{}).Statement

		assert.Equal(t, "SELECT * FROM `fakes` WHERE `fakes`.`deleted_at` IS NULL", stat.SQL.String())
	})

	t.Run("asc sort", func(t *testing.T) {
		stat := gdb.Scopes(scopeSort([]byte("name"))).Find(&Fake{}).Statement

		assert.Equal(t, "SELECT * FROM `fakes` WHERE `fakes`.`deleted_at` IS NULL ORDER BY name", stat.SQL.String())
	})

	t.Run("desc sort", func(t *testing.T) {
		stat := gdb.Scopes(scopeSort([]byte("-name"))).Find(&Fake{}).Statement

		assert.Equal(t, "SELECT * FROM `fakes` WHERE `fakes`.`deleted_at` IS NULL ORDER BY name desc", stat.SQL.String())
	})

	t.Run("two column sort", func(t *testing.T) {
		stat := gdb.Scopes(scopeSort([]byte("-name,key"))).Find(&Fake{}).Statement

		assert.Equal(t, "SELECT * FROM `fakes` WHERE `fakes`.`deleted_at` IS NULL ORDER BY name desc,key", stat.SQL.String())
	})
}

type Fake struct {
	Dao
	F string
}

func mockGdb(t assert.TestingT, dst ...interface{}) *gorm.DB {
	db, err := gorm.Open(
		sqlite.Open(fmt.Sprintf("file:%d?mode=memory&cache=shared&_fk=1", time.Now().UnixNano())),
		&gorm.Config{Logger: disabledLogger{}})

	assert.Nil(t, err)

	if len(dst) > 0 {
		assert.Nil(t, db.AutoMigrate(dst...))
	}

	return db
}

func dryRunSession(t assert.TestingT) *gorm.DB {
	return mockGdb(t).Session(&gorm.Session{DryRun: true, Logger: disabledLogger{}})
}

type disabledLogger struct{}

func (disabledLogger) LogMode(logger.LogLevel) logger.Interface {
	return disabledLogger{}
}
func (disabledLogger) Info(context.Context, string, ...interface{})                    {}
func (disabledLogger) Warn(context.Context, string, ...interface{})                    {}
func (disabledLogger) Error(context.Context, string, ...interface{})                   {}
func (disabledLogger) Trace(context.Context, time.Time, func() (string, int64), error) {}
