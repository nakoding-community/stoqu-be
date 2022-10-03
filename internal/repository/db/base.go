package db

import (
	"context"
	"fmt"
	"math"
	"strings"

	"gitlab.com/stoqu/stoqu-be/pkg/util/ctxval"

	abstraction "gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"

	"github.com/jackc/pgconn"
	"gorm.io/gorm"
)

type (
	Base[T any] interface {
		getConn(ctx context.Context) *gorm.DB
		maskError(err error) error
		buildFilterSort(ctx context.Context, name string, query *gorm.DB, filterParam abstraction.Filter)
		buildPagination(ctx context.Context, query *gorm.DB, pagination abstraction.Pagination) *abstraction.PaginationInfo

		Find(ctx context.Context, filterParam abstraction.Filter) ([]T, *abstraction.PaginationInfo, error)
		FindByID(ctx context.Context, id string) (*T, error)
		FindByCode(ctx context.Context, code string) (*T, error)
		FindByName(ctx context.Context, name string) (*T, error)
		Create(ctx context.Context, data T) (T, error)
		Creates(ctx context.Context, data []T) ([]T, error)
		UpdateByID(ctx context.Context, id string, data T) (T, error)
		DeleteByID(ctx context.Context, id string) error
	}

	base[T any] struct {
		conn       *gorm.DB
		entity     T
		entityName string
	}
)

func NewBase[T any](conn *gorm.DB, entity T, entityName string) Base[T] {
	return &base[T]{
		conn,
		entity,
		entityName,
	}
}

func (m *base[T]) getConn(ctx context.Context) *gorm.DB {
	tx := ctxval.GetTrxValue(ctx)
	fmt.Println("DEBUG HERE - get conn", tx)
	if tx != nil {
		return tx
	}
	return m.conn
}

func (m *base[T]) buildFilterSort(ctx context.Context, name string, query *gorm.DB, filterParam abstraction.Filter) {
	for _, filter := range filterParam.Query {
		query.Where(filter.Field+" = ?", filter.Value)
	}

	for i := range filterParam.SortBy {
		sortBys := strings.Split(filterParam.SortBy[i], ",")
		for _, sortBy := range sortBys {
			prefix := sortBy[:1]
			sortType := "asc"
			if prefix == "-" {
				sortType = "desc"
				sortBy = sortBy[1:]
			}

			sortArg := fmt.Sprintf("%s.%s %s", name, sortBy, sortType)
			query = query.Order(sortArg)
		}
	}
}

func (m *base[T]) buildPagination(ctx context.Context, tx *gorm.DB, pagination abstraction.Pagination) *abstraction.PaginationInfo {
	info := &abstraction.PaginationInfo{}
	limit := 10
	if pagination.Limit != nil {
		limit = *pagination.Limit
	}

	page := 0
	if pagination.Page != nil && *pagination.Page >= 0 {
		page = *pagination.Page
	}

	tx.Count(&info.Count)
	offset := (page - 1) * limit
	tx.Limit(limit).Offset(offset)
	info.TotalPage = int64(math.Ceil(float64(info.Count) / float64(limit)))
	info.Pagination = pagination

	return info
}

func (m *base[T]) maskError(err error) error {
	if err != nil {
		// not found
		if err == gorm.ErrRecordNotFound {
			return res.ErrorBuilder(res.Constant.Error.NotFound, err)
		}

		if pqErr, ok := err.(*pgconn.PgError); ok {
			// duplicate data
			if pqErr.Code == "23505" {
				return res.ErrorBuilder(res.Constant.Error.Duplicate, err)
			}
		}

		return res.ErrorBuilder(res.Constant.Error.UnprocessableEntity, err)
	}

	return nil
}

func (m *base[T]) Find(ctx context.Context, filterParam abstraction.Filter) ([]T, *abstraction.PaginationInfo, error) {
	query := m.getConn(ctx).Model(m.entity)

	m.buildFilterSort(ctx, m.entityName, query, filterParam)
	info := m.buildPagination(ctx, query, filterParam.Pagination)

	result := []T{}
	err := query.WithContext(ctx).Find(&result).Error

	if err != nil {
		return nil, info, err
	}
	return result, info, nil
}

func (m *base[T]) FindByID(ctx context.Context, id string) (*T, error) {
	query := m.getConn(ctx).Model(m.entity)
	result := new(T)
	err := query.WithContext(ctx).Where("id", id).First(result).Error
	if err != nil {
		return nil, m.maskError(err)
	}
	return result, nil
}

func (m *base[T]) FindByCode(ctx context.Context, code string) (*T, error) {
	query := m.getConn(ctx).Model(m.entity)
	result := new(T)
	err := query.WithContext(ctx).Where("code", code).First(result).Error
	if err != nil {
		return nil, m.maskError(err)
	}
	return result, nil
}

func (m *base[T]) FindByName(ctx context.Context, name string) (*T, error) {
	query := m.getConn(ctx).Model(m.entity)
	result := new(T)
	err := query.WithContext(ctx).Where("name", name).First(result).Error
	if err != nil {
		return nil, m.maskError(err)
	}
	return result, nil
}

func (m *base[T]) Create(ctx context.Context, data T) (T, error) {
	query := m.getConn(ctx).Model(m.entity)
	err := query.WithContext(ctx).Create(&data).Error
	return data, m.maskError(err)
}

func (m *base[T]) Creates(ctx context.Context, data []T) ([]T, error) {
	err := m.getConn(ctx).Model(m.entity).WithContext(ctx).Create(&data).Error
	return data, m.maskError(err)
}

func (m *base[T]) UpdateByID(ctx context.Context, id string, data T) (T, error) {
	err := m.getConn(ctx).Model(&data).WithContext(ctx).Where("id = ?", id).Updates(&data).Error
	return data, m.maskError(err)
}

func (m *base[T]) DeleteByID(ctx context.Context, id string) error {
	err := m.getConn(ctx).WithContext(ctx).Where("id = ?", id).Delete(m.entity).Error
	return m.maskError(err)
}
