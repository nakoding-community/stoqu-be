package db

import (
	"context"

	abstraction "gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"

	"gorm.io/gorm"
)

type (
	ConvertionUnit interface {
		// !TODO mockgen doesn't support embedded interface yet
		// !TODO but already discussed in this thread https://github.com/golang/mock/issues/621, lets wait for the release
		// Base[model.ConvertionUnitModel]

		// Base
		Find(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]model.ConvertionUnitModel, *abstraction.PaginationInfo, error)
		FindByID(ctx context.Context, id string) (*model.ConvertionUnitModel, error)
		FindByCode(ctx context.Context, code string) (*model.ConvertionUnitModel, error)
		FindByName(ctx context.Context, name string) (*model.ConvertionUnitModel, error)
		Create(ctx context.Context, data model.ConvertionUnitModel) (model.ConvertionUnitModel, error)
		Creates(ctx context.Context, data []model.ConvertionUnitModel) ([]model.ConvertionUnitModel, error)
		UpdateByID(ctx context.Context, id string, data model.ConvertionUnitModel) (model.ConvertionUnitModel, error)
		DeleteByID(ctx context.Context, id string) error
		Count(ctx context.Context) (int64, error)
	}

	convertionUnit struct {
		Base[model.ConvertionUnitModel]
	}
)

func NewConvertionUnit(conn *gorm.DB) ConvertionUnit {
	model := model.ConvertionUnitModel{}
	base := NewBase(conn, model, model.TableName())
	return &convertionUnit{
		base,
	}
}
