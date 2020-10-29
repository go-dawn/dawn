package gormx

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Fake struct{ gorm.Model }

func MockGdb(t assert.TestingT) *gorm.DB {
	db, err := gorm.Open(
		sqlite.Open("file:dawn?mode=memory&cache=shared&_fk=1"),
		&gorm.Config{DryRun: true})

	assert.Nil(t, err)
	return db
}

func Test_Scope_Paginate(t *testing.T) {
	gdb := MockGdb(t)

	t.Run("int", func(t *testing.T) {
		stat := gdb.Scopes(ScopePaginate(1, 10)).Find(&Fake{}).Statement

		assert.Equal(t, "SELECT * FROM `fakes` LIMIT 10", stat.SQL.String())
	})

	t.Run("offset", func(t *testing.T) {
		stat := gdb.Scopes(ScopePaginate(2, 10)).Find(&Fake{}).Statement

		assert.Equal(t, "SELECT * FROM `fakes` LIMIT 10 OFFSET 10", stat.SQL.String())
	})
}

func Test_Scope_Search(t *testing.T) {
	gdb := MockGdb(t)

	t.Run("empty object", func(t *testing.T) {
		stat := gdb.Scopes(ScopeSearch([]byte(`{}`))).Find(&Fake{}).Statement

		assert.Equal(t, "SELECT * FROM `fakes`", stat.SQL.String())
	})

	t.Run("number", func(t *testing.T) {
		stat := gdb.Scopes(ScopeSearch([]byte(`{"id":1}`))).Find(&Fake{}).Statement

		assert.Equal(t, "SELECT * FROM `fakes` WHERE id = ?", stat.SQL.String())
	})

	t.Run("like", func(t *testing.T) {
		stat := gdb.Scopes(ScopeSearch([]byte(`{"name":"k"}`))).Find(&Fake{}).Statement

		assert.Equal(t, "SELECT * FROM `fakes` WHERE name LIKE ?", stat.SQL.String())
		assert.Equal(t, []interface{}{"%k%"}, stat.Vars)
	})

	t.Run("in", func(t *testing.T) {
		t.Run("number", func(t *testing.T) {
			stat := gdb.Scopes(ScopeSearch([]byte(`{"name":[1.1,2.2,3.3]}`))).
				Find(&Fake{}).Statement

			assert.Equal(t, "SELECT * FROM `fakes` WHERE name IN (?,?,?)", stat.SQL.String())
			assert.Equal(t, []interface{}{1.1, 2.2, 3.3}, stat.Vars)
		})
		t.Run("string", func(t *testing.T) {
			stat := gdb.Scopes(ScopeSearch([]byte(`{"name":["1","2","3"]}`))).
				Find(&Fake{}).Statement
			assert.Equal(t, "SELECT * FROM `fakes` WHERE name IN (?,?,?)", stat.SQL.String())
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
					stat := gdb.Scopes(ScopeSearch([]byte(search))).Find(&Fake{}).Statement

					exp := fmt.Sprintf("SELECT * FROM `fakes` WHERE c %s ?", opt)
					assert.Equal(t, exp, stat.SQL.String())
					assert.Equal(t, []interface{}{val}, stat.Vars)
				})
			}
		}

		t.Run("><", func(t *testing.T) {
			stat := gdb.Scopes(ScopeSearch([]byte(`{"c$><":["2020-01-01", "2020-03-01"]}`))).Find(&Fake{}).Statement

			assert.Equal(t, "SELECT * FROM `fakes` WHERE c BETWEEN ? AND ?", stat.SQL.String())
			assert.Equal(t, []interface{}{"2020-01-01", "2020-03-01"}, stat.Vars)
		})

		t.Run("bool", func(t *testing.T) {
			stat := gdb.Scopes(ScopeSearch([]byte(`{"ok":true}`))).Find(&Fake{}).Statement

			assert.Equal(t, "SELECT * FROM `fakes` WHERE ok = ?", stat.SQL.String())
			assert.Equal(t, []interface{}{true}, stat.Vars)
		})
	})
}

func Benchmark_Scope_Search(b *testing.B) {
	gdb := MockGdb(b)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		gdb.Session(&gorm.Session{DryRun: false}).
			Scopes(ScopeSearch([]byte(`{"ID":1,"created_at$><":["2020-01-01", "2020-03-01"]}`))).
			Find(&Fake{})
	}
}

func Test_Scope_Sort(t *testing.T) {
	gdb := MockGdb(t)

	t.Run("no sort", func(t *testing.T) {
		stat := gdb.Scopes(ScopeSort([]byte(""))).Find(&Fake{}).Statement

		assert.Equal(t, "SELECT * FROM `fakes`", stat.SQL.String())
	})

	t.Run("asc sort", func(t *testing.T) {
		stat := gdb.Scopes(ScopeSort([]byte("name"))).Find(&Fake{}).Statement

		assert.Equal(t, "SELECT * FROM `fakes` ORDER BY name", stat.SQL.String())
	})

	t.Run("desc sort", func(t *testing.T) {
		stat := gdb.Scopes(ScopeSort([]byte("-name"))).Find(&Fake{}).Statement

		assert.Equal(t, "SELECT * FROM `fakes` ORDER BY name desc", stat.SQL.String())
	})

	t.Run("two column sort", func(t *testing.T) {
		stat := gdb.Scopes(ScopeSort([]byte("-name,key"))).Find(&Fake{}).Statement

		assert.Equal(t, "SELECT * FROM `fakes` ORDER BY name desc,key", stat.SQL.String())
	})
}
