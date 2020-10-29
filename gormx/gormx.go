package gormx

import (
	"reflect"

	"github.com/go-dawn/dawn/gormx/schema"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// Paginate gets pagination based on index query and transfer
// dao to domain object
func Paginate(scope *gorm.DB, q schema.IndexQuery, daos, data interface{}) (*schema.Pagination, error) {
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

	return &schema.Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    int(total),
		// transfer pointer to slice
		Data: reflect.ValueOf(data).Elem().Interface(),
	}, nil
}
