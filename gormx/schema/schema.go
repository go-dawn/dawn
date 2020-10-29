package schema

import (
	"time"

	"github.com/go-dawn/dawn/config"

	"github.com/valyala/fasthttp"
)

type Dao struct {
	ID        uint32 `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

// Fields alias for map[string]interface
type Fields = map[string]interface{}

// IndexQuery wraps helper functions of query string
type IndexQuery struct {
	args *fasthttp.Args
}

// NewIndexQuery wraps *fasthttp.Args
func NewIndexQuery(v *fasthttp.Args) IndexQuery {
	return IndexQuery{v}
}

// Search gets search data in json format
func (q IndexQuery) Search() []byte {
	return q.args.Peek("search")
}

// Sort gets sort params
func (q IndexQuery) Sort() []byte {
	return q.args.Peek("sort")
}

// Page gets index page
func (q IndexQuery) Page() int {
	n := q.args.GetUintOrZero("page")
	if n <= 0 {
		return 1
	}
	return n
}

// PageSize gets page size of index
func (q IndexQuery) PageSize() int {
	n := q.args.GetUintOrZero("pageSize")
	if n <= 0 {
		return config.GetInt("http.pageSize", 15)
	}
	return n
}

// PageInfo gets page and page size
func (q IndexQuery) PageInfo() (int, int) {
	return q.Page(), q.PageSize()
}

// Pagination contains page, page size, total count
// and list data
type Pagination struct {
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
	Total    int         `json:"total"`
	Data     interface{} `json:"data"`
}
