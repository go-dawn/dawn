package gormx

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"

	"github.com/gofiber/fiber/v2"

	"github.com/go-dawn/dawn/config"
	"github.com/jinzhu/copier"
	"github.com/valyala/fasthttp"
	"gorm.io/gorm"
)

type Dao struct {
	ID        uint32 `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *gorm.DeletedAt `gorm:"index"`
}

// Pagination contains page, page size, total count
// and list data
type Pagination struct {
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
	Total    int         `json:"total"`
	Data     interface{} `json:"data"`
}

// Paginate gets pagination based on index query and transfer
// dao to domain object
func Paginate(scope *gorm.DB, c *fiber.Ctx, daos, data interface{}) (*Pagination, error) {
	q := indexQuery{c.Context().QueryArgs()}

	page, pageSize := q.pageInfo()

	scope = scope.Scopes(
		scopeSearch(q.search()),
		scopeSort(q.sort()),
		scopePaginate(page, pageSize),
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

// indexQuery wraps helper functions of query string
type indexQuery struct {
	*fasthttp.Args
}

// search gets search data in json format
func (q indexQuery) search() []byte {
	return q.Peek("search")
}

// sort gets sort params
func (q indexQuery) sort() []byte {
	return q.Peek("sort")
}

// page gets index page
func (q indexQuery) page() int {
	n := q.GetUintOrZero("page")
	if n <= 0 {
		return 1
	}
	return n
}

// pageSize gets page size of index
func (q indexQuery) pageSize() int {
	n := q.GetUintOrZero("pageSize")
	if n <= 0 {
		return config.GetInt("http.pageSize", 15)
	}
	return n
}

// pageInfo gets page and page size
func (q indexQuery) pageInfo() (int, int) {
	return q.page(), q.pageSize()
}

func scopePaginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(scope *gorm.DB) *gorm.DB {
		offset := (page - 1) * pageSize

		return scope.Offset(offset).Limit(pageSize)
	}
}

func scopeSearch(search []byte) func(db *gorm.DB) *gorm.DB {
	return func(scope *gorm.DB) *gorm.DB {
		params := map[string]interface{}{}
		if err := jsoniter.Unmarshal(search, &params); err == nil {
			//"name":"jone" => name like ?, %jone%
			//"name":["jone","kj"] => name in (?), ["jone","kj"]
			//"free": true | false => free = ?, true
			//"name$<>": "jone" name <> ? , "jone"
			//"date$<>":["2020-12-12",""] date not in (?), "", ""
			//"date$><":["2020-12-12",""] date between ? and ?, "", ""
			//"date$<":"" => date < ?
			//"date$<=":"" => date <= ?
			//"date$>=":"" => date >= ?
			//"date$>":"" => date > ?
			for key, val := range params {
				// 使用$分割列名和操作符
				strs := strings.Split(key, "$")
				// 不带操作符
				if len(strs) == 1 {
					column := strs[0]
					// 根据值的类型附加过滤条件
					switch reflect.TypeOf(val).Kind() {
					case reflect.String:
						// 字符串全部用like过滤
						scope = scope.Where(column+" LIKE ?", fmt.Sprintf("%%%v%%", val))
					case reflect.Bool, reflect.Float64:
						scope = scope.Where(column+" = ?", val)
					case reflect.Slice, reflect.Array:
						// 数组使用 IN
						if arr, ok := val.([]interface{}); ok {
							scope = scope.Where(column+" IN (?)", arr)
						}
					}
				}
				// 带操作符
				if len(strs) == 2 {
					column, opt := strs[0], strs[1]
					switch opt {
					// 比较操作符
					case ">", ">=", "<", "<=":
						switch val.(type) {
						// 只支持字符串和数字类型
						case string, float64:
							scope = scope.Where(fmt.Sprintf("%s %s ?", column, opt), val)
						}
					case "><": // between 操作符
						switch reflect.TypeOf(val).Kind() {
						case reflect.Slice, reflect.Array:
							// 只支持长度为2的数组
							if arr, ok := val.([]interface{}); ok && len(arr) == 2 {
								scope = scope.Where(strs[0]+" BETWEEN ? AND ?", arr[0], arr[1])
							}
						}
					}
				}

				// 忽略其他情况
			}
		}

		return scope
	}
}

var sortSep = []byte(",")

func scopeSort(sort []byte) func(db *gorm.DB) *gorm.DB {
	return func(scope *gorm.DB) *gorm.DB {
		if len(sort) == 0 {
			return scope
		}
		buf := new(bytes.Buffer)
		for _, key := range bytes.Split(sort, sortSep) {
			if key[0] == '-' {
				buf.Write(key[1:])
				buf.WriteString(" desc,")
			} else {
				buf.Write(key)
				buf.WriteByte(',')
			}
		}
		if buf.Len() != 0 {
			scope = scope.Order(buf.String()[:buf.Len()-1])
		}
		return scope
	}
}
