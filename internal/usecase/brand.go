package usecase

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"gitlab.com/stoqu/stoqu-be/internal/config"

	"gitlab.com/stoqu/stoqu-be/internal/factory/repository"
	"gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	"gitlab.com/stoqu/stoqu-be/internal/model/dto"
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"
	"gitlab.com/stoqu/stoqu-be/pkg/constant"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"
	"gitlab.com/stoqu/stoqu-be/pkg/util/str"
	"gitlab.com/stoqu/stoqu-be/pkg/util/trxmanager"
)

type Brand interface {
	Find(ctx context.Context, filterParam abstraction.Filter) ([]dto.BrandResponse, abstraction.PaginationInfo, error)
	FindByID(ctx context.Context, payload dto.ByIDRequest) (dto.BrandResponse, error)
	Create(ctx context.Context, payload dto.CreateBrandRequest) (dto.BrandResponse, error)
	Update(ctx context.Context, payload dto.UpdateBrandRequest) (dto.BrandResponse, error)
	Delete(ctx context.Context, payload dto.ByIDRequest) (dto.BrandResponse, error)
}

type brand struct {
	Repo repository.Factory
	Cfg  *config.Configuration
}

func NewBrand(cfg *config.Configuration, f repository.Factory) Brand {
	return &brand{f, cfg}
}

func (u *brand) Find(ctx context.Context, filterParam abstraction.Filter) (result []dto.BrandResponse, pagination abstraction.PaginationInfo, err error) {
	var search *abstraction.Search
	if filterParam.Search != "" {
		searchQuery := "lower(name) LIKE ?"
		searchValStr := "%" + strings.ToLower(filterParam.Search) + "%"
		search = &abstraction.Search{
			Query: searchQuery,
			Args:  []interface{}{searchValStr},
		}
	}
	brands, info, err := u.Repo.Brand.Find(ctx, filterParam, search)
	if err != nil {
		return nil, pagination, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}
	pagination = *info

	for _, brand := range brands {
		variants, _, err := u.Repo.Variant.Find(ctx, abstraction.Filter{
			Query: []abstraction.FilterQuery{
				{
					Field: "brand_id",
					Value: brand.ID,
				},
			},
		}, nil)
		if err == nil {
			brand.Variants = append(brand.Variants, variants...)
		}

		result = append(result, dto.BrandResponse{
			BrandModel: brand,
		})
	}

	return result, pagination, nil
}

func (u *brand) FindByID(ctx context.Context, payload dto.ByIDRequest) (dto.BrandResponse, error) {
	var result dto.BrandResponse

	brand, err := u.Repo.Brand.FindByID(ctx, payload.ID)
	if err != nil {
		return result, err
	}

	result = dto.BrandResponse{
		BrandModel: *brand,
	}

	return result, nil
}

func (u *brand) Create(ctx context.Context, payload dto.CreateBrandRequest) (result dto.BrandResponse, err error) {
	var (
		brandID = uuid.New().String()
		brand   = model.BrandModel{
			Entity: model.Entity{
				ID: brandID,
			},
			BrandEntity: model.BrandEntity{
				Code:       str.GenCode(constant.CODE_BRAND_PREFIX),
				Name:       payload.Name,
				SupplierID: payload.SupplierID,
			},
		}
	)

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		_, err = u.Repo.Brand.Create(ctx, brand)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.BrandResponse{
		BrandModel: brand,
	}

	return result, nil
}

func (u *brand) Update(ctx context.Context, payload dto.UpdateBrandRequest) (result dto.BrandResponse, err error) {
	var (
		brandData = &model.BrandModel{
			BrandEntity: model.BrandEntity{
				Name:       payload.Name,
				SupplierID: payload.SupplierID,
			},
		}
	)

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		_, err = u.Repo.Brand.UpdateByID(ctx, payload.ID, *brandData)
		if err != nil {
			return err
		}

		brandData, err = u.Repo.Brand.FindByID(ctx, payload.ID)
		if err != nil {
			return res.ErrorBuilder(res.Constant.Error.NotFound, err)
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.BrandResponse{
		BrandModel: *brandData,
	}

	return result, nil
}

func (u *brand) Delete(ctx context.Context, payload dto.ByIDRequest) (result dto.BrandResponse, err error) {
	var data *model.BrandModel

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		data, err = u.Repo.Brand.FindByID(ctx, payload.ID)
		if err != nil {
			return err
		}

		err = u.Repo.Brand.DeleteByID(ctx, payload.ID)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.BrandResponse{
		BrandModel: *data,
	}

	return result, nil
}
