package usecase

import (
	"context"

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

type Variant interface {
	Find(ctx context.Context, filterParam abstraction.Filter) ([]dto.VariantResponse, abstraction.PaginationInfo, error)
	FindByID(ctx context.Context, payload dto.ByIDRequest) (dto.VariantResponse, error)
	Create(ctx context.Context, payload dto.CreateVariantRequest) (dto.VariantResponse, error)
	Update(ctx context.Context, payload dto.UpdateVariantRequest) (dto.VariantResponse, error)
	Delete(ctx context.Context, payload dto.ByIDRequest) (dto.VariantResponse, error)
}

type variant struct {
	Repo repository.Factory
	Cfg  *config.Configuration
}

func NewVariant(cfg *config.Configuration, f repository.Factory) Variant {
	return &variant{f, cfg}
}

func (u *variant) Find(ctx context.Context, filterParam abstraction.Filter) (result []dto.VariantResponse, pagination abstraction.PaginationInfo, err error) {
	variants, info, err := u.Repo.Variant.Find(ctx, filterParam, nil)
	if err != nil {
		return nil, pagination, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}
	pagination = *info

	for _, variant := range variants {
		result = append(result, dto.VariantResponse{
			VariantModel: variant,
		})
	}

	return result, pagination, nil
}

func (u *variant) FindByID(ctx context.Context, payload dto.ByIDRequest) (dto.VariantResponse, error) {
	var result dto.VariantResponse

	variant, err := u.Repo.Variant.FindByID(ctx, payload.ID)
	if err != nil {
		return result, err
	}

	result = dto.VariantResponse{
		VariantModel: *variant,
	}

	return result, nil
}

func (u *variant) Create(ctx context.Context, payload dto.CreateVariantRequest) (result dto.VariantResponse, err error) {
	var (
		variantID = uuid.New().String()
		variant   = model.VariantModel{
			Entity: model.Entity{
				ID: variantID,
			},
			VariantEntity: model.VariantEntity{
				Code:       str.GenCode(constant.CODE_VARIANT_PREFIX),
				Name:       payload.Name,
				ITL:        payload.ITL,
				UniqueCode: payload.UniqueCode,
				BrandID:    payload.BrandID,
			},
		}
	)

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		_, err = u.Repo.Variant.Create(ctx, variant)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.VariantResponse{
		VariantModel: variant,
	}

	return result, nil
}

func (u *variant) Update(ctx context.Context, payload dto.UpdateVariantRequest) (result dto.VariantResponse, err error) {
	var (
		variantData = &model.VariantModel{
			VariantEntity: model.VariantEntity{
				Name:       payload.Name,
				ITL:        payload.ITL,
				UniqueCode: payload.UniqueCode,
			},
		}
	)

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		_, err = u.Repo.Variant.UpdateByID(ctx, payload.ID, *variantData)
		if err != nil {
			return err
		}

		variantData, err = u.Repo.Variant.FindByID(ctx, payload.ID)
		if err != nil {
			return res.ErrorBuilder(res.Constant.Error.NotFound, err)
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.VariantResponse{
		VariantModel: *variantData,
	}

	return result, nil
}

func (u *variant) Delete(ctx context.Context, payload dto.ByIDRequest) (result dto.VariantResponse, err error) {
	var data *model.VariantModel

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		data, err = u.Repo.Variant.FindByID(ctx, payload.ID)
		if err != nil {
			return err
		}

		err = u.Repo.Variant.DeleteByID(ctx, payload.ID)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.VariantResponse{
		VariantModel: *data,
	}

	return result, nil
}
