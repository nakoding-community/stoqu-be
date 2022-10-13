package usecase

import (
	"context"
	"strconv"
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

type ProductLookup interface {
	Find(ctx context.Context, filterParam abstraction.Filter) ([]dto.ProductLookupResponse, abstraction.PaginationInfo, error)
	FindByID(ctx context.Context, payload dto.ByIDRequest) (dto.ProductLookupResponse, error)
	Create(ctx context.Context, payload dto.CreateProductLookupRequest) (dto.ProductLookupResponse, error)
	Update(ctx context.Context, payload dto.UpdateProductLookupRequest) (dto.ProductLookupResponse, error)
	Delete(ctx context.Context, payload dto.ByIDRequest) (dto.ProductLookupResponse, error)
}

type productLookup struct {
	Repo repository.Factory
	Cfg  *config.Configuration
}

func NewProductLookup(cfg *config.Configuration, f repository.Factory) ProductLookup {
	return &productLookup{f, cfg}
}

func (u *productLookup) Find(ctx context.Context, filterParam abstraction.Filter) (result []dto.ProductLookupResponse, pagination abstraction.PaginationInfo, err error) {
	var search *abstraction.Search
	if filterParam.Search != "" {
		searchQuery := "lower(code) LIKE ? OR type_value = ? OR remaining_type_value = ? OR remaining_type_value_before = ?"
		searchVal := "%" + strings.ToLower(filterParam.Search) + "%"
		searchValFloat, _ := strconv.ParseFloat(filterParam.Search, 64)
		search = &abstraction.Search{
			Query: searchQuery,
			Args:  []interface{}{searchVal, searchValFloat, searchValFloat, searchValFloat},
		}
	}

	productLookups, info, err := u.Repo.ProductLookup.Find(ctx, filterParam, search)
	if err != nil {
		return nil, pagination, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}
	pagination = *info

	for _, productLookup := range productLookups {
		result = append(result, dto.ProductLookupResponse{
			ProductLookupModel: productLookup,
		})
	}

	return result, pagination, nil
}

func (u *productLookup) FindByID(ctx context.Context, payload dto.ByIDRequest) (dto.ProductLookupResponse, error) {
	var result dto.ProductLookupResponse

	productLookup, err := u.Repo.ProductLookup.FindByID(ctx, payload.ID)
	if err != nil {
		return result, err
	}

	result = dto.ProductLookupResponse{
		ProductLookupModel: *productLookup,
	}

	return result, nil
}

func (u *productLookup) Create(ctx context.Context, payload dto.CreateProductLookupRequest) (result dto.ProductLookupResponse, err error) {
	var (
		productLookupID = uuid.New().String()
		productLookup   = model.ProductLookupModel{
			Entity: model.Entity{
				ID: productLookupID,
			},
			ProductLookupEntity: model.ProductLookupEntity{
				Code:                     str.GenCode(constant.CODE_PRODUCT_LOOKUP_PREFIX),
				IsSeal:                   payload.IsSeal,
				TypeValue:                payload.TypeValue,
				RemainingTypeValue:       payload.RemainingTypeValue,
				RemainingTypeValueBefore: payload.RemainingTypeValueBefore,
				ProductID:                payload.ProductID,
			},
		}
	)

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		_, err = u.Repo.ProductLookup.Create(ctx, productLookup)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.ProductLookupResponse{
		ProductLookupModel: productLookup,
	}

	return result, nil
}

func (u *productLookup) Update(ctx context.Context, payload dto.UpdateProductLookupRequest) (result dto.ProductLookupResponse, err error) {
	var (
		productLookupData = &model.ProductLookupModel{
			ProductLookupEntity: model.ProductLookupEntity{
				IsSeal:                   payload.IsSeal,
				TypeValue:                payload.TypeValue,
				RemainingTypeValue:       payload.RemainingTypeValue,
				RemainingTypeValueBefore: payload.RemainingTypeValueBefore,
			},
		}
	)

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		_, err = u.Repo.ProductLookup.UpdateByID(ctx, payload.ID, *productLookupData)
		if err != nil {
			return err
		}

		productLookupData, err = u.Repo.ProductLookup.FindByID(ctx, payload.ID)
		if err != nil {
			return res.ErrorBuilder(res.Constant.Error.NotFound, err)
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.ProductLookupResponse{
		ProductLookupModel: *productLookupData,
	}

	return result, nil
}

func (u *productLookup) Delete(ctx context.Context, payload dto.ByIDRequest) (result dto.ProductLookupResponse, err error) {
	var data *model.ProductLookupModel

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		data, err = u.Repo.ProductLookup.FindByID(ctx, payload.ID)
		if err != nil {
			return err
		}

		err = u.Repo.ProductLookup.DeleteByID(ctx, payload.ID)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.ProductLookupResponse{
		ProductLookupModel: *data,
	}

	return result, nil
}
