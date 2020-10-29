package gormx

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"gorm.io/gorm"
)

func ScopePaginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(scope *gorm.DB) *gorm.DB {
		offset := (page - 1) * pageSize

		return scope.Offset(offset).Limit(pageSize)
	}
}

func ScopeSearch(search []byte) func(db *gorm.DB) *gorm.DB {
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

func ScopeSort(sort []byte) func(db *gorm.DB) *gorm.DB {
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
