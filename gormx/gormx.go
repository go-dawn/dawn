package gormx

import (
	"reflect"
	"time"

	"github.com/go-dawn/dawn/config"
	"github.com/jinzhu/copier"
	"github.com/valyala/fasthttp"
	"gorm.io/gorm"
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

// Paginate gets pagination based on index query and transfer
// dao to domain object
func Paginate(scope *gorm.DB, q IndexQuery, daos, data interface{}) (*Pagination, error) {
	page, pageSize := q.PageInfo()

	scope = scope.Scopes(
		ScopeSearch(q.Search()),
		ScopeSort(q.Sort()),
		ScopePaginate(page, pageSize),
	)

	var (
		total   int64
		errChan = make(chan error)
	)

	go func() {
		errChan <- scope.Find(daos).Error
	}()

	// get total count
	err1, err2 := <-errChan, scope.Model(daos).Count(&total).Error

	if err1 != nil && err1 != gorm.ErrRecordNotFound {
		return nil, err1
	}

	if err2 != nil {
		return nil, err2
	}

	if err := copier.Copy(data, daos); err != nil {
		return nil, err
	}

	return &Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    int(total),
		// transfer pointer to slice
		Data: reflect.ValueOf(data).Elem().Interface(),
	}, nil
}
