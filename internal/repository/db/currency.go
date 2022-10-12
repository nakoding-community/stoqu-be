package db

import (
	"context"

	abstraction "gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"

	"gorm.io/gorm"
)

type (
	Currency interface {
		// !TODO mockgen doesn't support embedded interface yet
		// !TODO but already discussed in this thread https://github.com/golang/mock/issues/621, lets wait for the release
		// Base[model.CurrencyModel]

		// Base
		Find(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]model.CurrencyModel, *abstraction.PaginationInfo, error)
		FindByID(ctx context.Context, id string) (*model.CurrencyModel, error)
		FindByCode(ctx context.Context, code string) (*model.CurrencyModel, error)
		FindByName(ctx context.Context, name string) (*model.CurrencyModel, error)
		Create(ctx context.Context, data model.CurrencyModel) (model.CurrencyModel, error)
		Creates(ctx context.Context, data []model.CurrencyModel) ([]model.CurrencyModel, error)
		UpdateByID(ctx context.Context, id string, data model.CurrencyModel) (model.CurrencyModel, error)
		DeleteByID(ctx context.Context, id string) error
		Count(ctx context.Context) (int64, error)
	}

	currency struct {
		Base[model.CurrencyModel]
	}
)

func NewCurrency(conn *gorm.DB) Currency {
	model := model.CurrencyModel{}
	base := NewBase(conn, model, model.TableName())
	return &currency{
		base,
	}
}
