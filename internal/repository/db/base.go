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
		GetConn(ctx context.Context) *gorm.DB
		MaskError(err error) error
		BuildFilterSort(ctx context.Context, query *gorm.DB, filterParam abstraction.Filter)
		BuildPagination(ctx context.Context, query *gorm.DB, pagination abstraction.Pagination) *abstraction.PaginationInfo

		Find(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]T, *abstraction.PaginationInfo, error)
		FindByID(ctx context.Context, id string) (*T, error)
		FindByIDs(ctx context.Context, ids []string, sortBy string) ([]T, error)
		FindByCode(ctx context.Context, code string) (*T, error)
		FindByName(ctx context.Context, name string) (*T, error)
		Create(ctx context.Context, data T) (T, error)
		Creates(ctx context.Context, data []T) ([]T, error)
		UpdateByID(ctx context.Context, id string, data T) (T, error)
		DeleteByID(ctx context.Context, id string) error
		DeleteByIDs(ctx context.Context, ids []string) error
		Count(ctx context.Context) (int64, error)
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

func (m *base[T]) GetConn(ctx context.Context) *gorm.DB {
	tx := ctxval.GetTrxValue(ctx)
	if tx != nil {
		return tx
	}
	return m.conn
}

func (m *base[T]) BuildFilterSort(ctx context.Context, query *gorm.DB, filterParam abstraction.Filter) {
	for _, filter := range filterParam.Query {
		// handle default operator, if pass direct from usecase not handler
		if filter.Operator == "" {
			filter.Operator = "="
		}
		query.Where(filter.Field+" "+filter.Operator+" ?", filter.Value)
		// !TODO: why we dont add a full query instead ? so we able to customize the operator
		// ex: filter.Query = "order_trxs = ?"
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

			sortArg := fmt.Sprintf("%s %s", sortBy, sortType)
			query = query.Order(sortArg)
		}
	}
}

func (m *base[T]) BuildPagination(ctx context.Context, tx *gorm.DB, pagination abstraction.Pagination) *abstraction.PaginationInfo {
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

func (m *base[T]) MaskError(err error) error {
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

func (m *base[T]) Find(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]T, *abstraction.PaginationInfo, error) {
	query := m.GetConn(ctx).Model(m.entity)
	if search != nil {
		query = query.Where(search.Query, search.Args...)
	}

	m.BuildFilterSort(ctx, query, filterParam)
	info := m.BuildPagination(ctx, query, filterParam.Pagination)

	result := []T{}
	err := query.WithContext(ctx).Find(&result).Error

	if err != nil {
		return nil, info, m.MaskError(err)
	}
	return result, info, nil
}

func (m *base[T]) FindByID(ctx context.Context, id string) (*T, error) {
	query := m.GetConn(ctx).Model(m.entity).Where("id", id)
	result := new(T)
	err := query.WithContext(ctx).First(result).Error
	if err != nil {
		return nil, m.MaskError(err)
	}
	return result, nil
}

func (m *base[T]) FindByIDs(ctx context.Context, ids []string, sortBy string) ([]T, error) {
	query := m.GetConn(ctx).Model(m.entity).Where("id IN ?", ids)
	if sortBy != "" {
		query.Order(sortBy)
	}

	result := []T{}
	err := query.WithContext(ctx).Find(&result).Error
	if err != nil {
		return nil, m.MaskError(err)
	}
	return result, nil
}

func (m *base[T]) FindByCode(ctx context.Context, code string) (*T, error) {
	query := m.GetConn(ctx).Model(m.entity).Where("code", code)
	result := new(T)
	err := query.WithContext(ctx).First(result).Error
	if err != nil {
		return nil, m.MaskError(err)
	}
	return result, nil
}

func (m *base[T]) FindByName(ctx context.Context, name string) (*T, error) {
	query := m.GetConn(ctx).Model(m.entity).Where("name", name)
	result := new(T)
	err := query.WithContext(ctx).First(result).Error
	if err != nil {
		return nil, m.MaskError(err)
	}
	return result, nil
}

func (m *base[T]) Create(ctx context.Context, data T) (T, error) {
	query := m.GetConn(ctx).Model(m.entity)
	err := query.WithContext(ctx).Create(&data).Error
	return data, m.MaskError(err)
}

func (m *base[T]) Creates(ctx context.Context, data []T) ([]T, error) {
	err := m.GetConn(ctx).Model(m.entity).WithContext(ctx).Create(&data).Error
	return data, m.MaskError(err)
}

func (m *base[T]) UpdateByID(ctx context.Context, id string, data T) (T, error) {
	err := m.GetConn(ctx).Model(&data).WithContext(ctx).Where("id = ?", id).Updates(&data).Error
	return data, m.MaskError(err)
}

func (m *base[T]) DeleteByID(ctx context.Context, id string) error {
	err := m.GetConn(ctx).WithContext(ctx).Where("id = ?", id).Delete(m.entity).Error
	return m.MaskError(err)
}

func (m *base[T]) DeleteByIDs(ctx context.Context, ids []string) error {
	err := m.GetConn(ctx).WithContext(ctx).Where("id IN ?", ids).Delete(m.entity).Error
	return m.MaskError(err)
}

func (m *base[T]) Count(ctx context.Context) (int64, error) {
	var count int64
	err := m.GetConn(ctx).Model(m.entity).WithContext(ctx).Count(&count).Error
	return count, err
}
