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

type StockLookup interface {
	Find(ctx context.Context, filterParam abstraction.Filter) ([]dto.StockLookupResponse, abstraction.PaginationInfo, error)
	FindByID(ctx context.Context, payload dto.ByIDRequest) (dto.StockLookupResponse, error)
	Create(ctx context.Context, payload dto.CreateStockLookupRequest) (dto.StockLookupResponse, error)
	Update(ctx context.Context, payload dto.UpdateStockLookupRequest) (dto.StockLookupResponse, error)
	Delete(ctx context.Context, payload dto.ByIDRequest) (dto.StockLookupResponse, error)
}

type stockLookup struct {
	Repo repository.Factory
	Cfg  *config.Configuration
}

func NewStockLookup(cfg *config.Configuration, f repository.Factory) StockLookup {
	return &stockLookup{f, cfg}
}

func (u *stockLookup) Find(ctx context.Context, filterParam abstraction.Filter) (result []dto.StockLookupResponse, pagination abstraction.PaginationInfo, err error) {
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

	stockLookups, info, err := u.Repo.StockLookup.Find(ctx, filterParam, search)
	if err != nil {
		return nil, pagination, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}
	pagination = *info

	for _, stockLookup := range stockLookups {
		result = append(result, dto.StockLookupResponse{
			StockLookupModel: stockLookup,
		})
	}

	return result, pagination, nil
}

func (u *stockLookup) FindByID(ctx context.Context, payload dto.ByIDRequest) (dto.StockLookupResponse, error) {
	var result dto.StockLookupResponse

	stockLookup, err := u.Repo.StockLookup.FindByID(ctx, payload.ID)
	if err != nil {
		return result, err
	}

	result = dto.StockLookupResponse{
		StockLookupModel: *stockLookup,
	}

	return result, nil
}

func (u *stockLookup) Create(ctx context.Context, payload dto.CreateStockLookupRequest) (result dto.StockLookupResponse, err error) {
	var (
		stockLookupID = uuid.New().String()
		stockLookup   = model.StockLookupModel{
			Entity: model.Entity{
				ID: stockLookupID,
			},
			StockLookupEntity: model.StockLookupEntity{
				Code:                 str.GenCode(constant.CODE_STOCK_LOOKUP_PREFIX),
				IsSeal:               payload.IsSeal,
				Value:                payload.Value,
				RemainingValue:       payload.RemainingValue,
				RemainingValueBefore: payload.RemainingValueBefore,
				StockID:              payload.StockID,
			},
		}
	)

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		_, err = u.Repo.StockLookup.Create(ctx, stockLookup)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.StockLookupResponse{
		StockLookupModel: stockLookup,
	}

	return result, nil
}

func (u *stockLookup) Update(ctx context.Context, payload dto.UpdateStockLookupRequest) (result dto.StockLookupResponse, err error) {
	var (
		stockLookupData = &model.StockLookupModel{
			StockLookupEntity: model.StockLookupEntity{
				IsSeal:               payload.IsSeal,
				Value:                payload.Value,
				RemainingValue:       payload.RemainingValue,
				RemainingValueBefore: payload.RemainingValueBefore,
			},
		}
	)

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		_, err = u.Repo.StockLookup.UpdateByID(ctx, payload.ID, *stockLookupData)
		if err != nil {
			return err
		}

		stockLookupData, err = u.Repo.StockLookup.FindByID(ctx, payload.ID)
		if err != nil {
			return res.ErrorBuilder(res.Constant.Error.NotFound, err)
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.StockLookupResponse{
		StockLookupModel: *stockLookupData,
	}

	return result, nil
}

func (u *stockLookup) Delete(ctx context.Context, payload dto.ByIDRequest) (result dto.StockLookupResponse, err error) {
	var data *model.StockLookupModel

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		data, err = u.Repo.StockLookup.FindByID(ctx, payload.ID)
		if err != nil {
			return err
		}

		err = u.Repo.StockLookup.DeleteByID(ctx, payload.ID)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.StockLookupResponse{
		StockLookupModel: *data,
	}

	return result, nil
}
