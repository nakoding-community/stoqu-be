package usecase

import (
	"context"
	"encoding/json"
	"math"
	"net/http"
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

type Currency interface {
	Find(ctx context.Context, filterParam abstraction.Filter) ([]dto.CurrencyResponse, abstraction.PaginationInfo, error)
	FindByID(ctx context.Context, payload dto.ByIDRequest) (dto.CurrencyResponse, error)
	Create(ctx context.Context, payload dto.CreateCurrencyRequest) (dto.CurrencyResponse, error)
	Update(ctx context.Context, payload dto.UpdateCurrencyRequest) (dto.CurrencyResponse, error)
	Delete(ctx context.Context, payload dto.ByIDRequest) (dto.CurrencyResponse, error)
	Convert(ctx context.Context, payload dto.ConvertCurrencyRequest) (result dto.ConvertCurrencyResponse, err error)
}

type currency struct {
	Repo repository.Factory
	Cfg  *config.Configuration
}

func NewCurrency(cfg *config.Configuration, f repository.Factory) Currency {
	return &currency{f, cfg}
}

func (u *currency) Find(ctx context.Context, filterParam abstraction.Filter) (result []dto.CurrencyResponse, pagination abstraction.PaginationInfo, err error) {
	var search *abstraction.Search
	if filterParam.Search != "" {
		searchQuery := "lower(code) LIKE ? OR lower(name) LIKE ?"
		searchVal := "%" + strings.ToLower(filterParam.Search) + "%"
		search = &abstraction.Search{
			Query: searchQuery,
			Args:  []interface{}{searchVal, searchVal},
		}
	}

	currencies, info, err := u.Repo.Currency.Find(ctx, filterParam, search)
	if err != nil {
		return nil, pagination, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}
	pagination = *info

	for _, currency := range currencies {
		result = append(result, dto.CurrencyResponse{
			CurrencyModel: currency,
		})
	}

	return result, pagination, nil
}

func (u *currency) FindByID(ctx context.Context, payload dto.ByIDRequest) (dto.CurrencyResponse, error) {
	var result dto.CurrencyResponse

	currency, err := u.Repo.Currency.FindByID(ctx, payload.ID)
	if err != nil {
		return result, err
	}

	result = dto.CurrencyResponse{
		CurrencyModel: *currency,
	}

	return result, nil
}

func (u *currency) Create(ctx context.Context, payload dto.CreateCurrencyRequest) (result dto.CurrencyResponse, err error) {
	var (
		currencyID = uuid.New().String()
		currency   = model.CurrencyModel{
			Entity: model.Entity{
				ID: currencyID,
			},
			CurrencyEntity: model.CurrencyEntity{
				Code:   str.GenCode(constant.CODE_CURRENCY_PREFIX),
				Name:   payload.Name,
				IsAuto: payload.IsAuto,
				Value:  payload.Value,
			},
		}
	)

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		_, err = u.Repo.Currency.Create(ctx, currency)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.CurrencyResponse{
		CurrencyModel: currency,
	}

	return result, nil
}

func (u *currency) Update(ctx context.Context, payload dto.UpdateCurrencyRequest) (result dto.CurrencyResponse, err error) {
	var (
		currencyData = &model.CurrencyModel{
			CurrencyEntity: model.CurrencyEntity{
				Name:   payload.Name,
				IsAuto: payload.IsAuto,
				Value:  payload.Value,
			},
		}
	)

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		_, err = u.Repo.Currency.UpdateByID(ctx, payload.ID, *currencyData)
		if err != nil {
			return err
		}

		currencyData, err = u.Repo.Currency.FindByID(ctx, payload.ID)
		if err != nil {
			return res.ErrorBuilder(res.Constant.Error.NotFound, err)
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.CurrencyResponse{
		CurrencyModel: *currencyData,
	}

	return result, nil
}

func (u *currency) Delete(ctx context.Context, payload dto.ByIDRequest) (result dto.CurrencyResponse, err error) {
	var data *model.CurrencyModel

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		data, err = u.Repo.Currency.FindByID(ctx, payload.ID)
		if err != nil {
			return err
		}

		err = u.Repo.Currency.DeleteByID(ctx, payload.ID)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.CurrencyResponse{
		CurrencyModel: *data,
	}

	return result, nil
}

func (u *currency) Convert(ctx context.Context, payload dto.ConvertCurrencyRequest) (result dto.ConvertCurrencyResponse, err error) {
	currency, err := u.Repo.Currency.FindByName(ctx, "IDR")
	if err != nil {
		return result, err
	}

	if currency.IsAuto {
		url := "https://cdn.jsdelivr.net/gh/fawazahmed0/currency-api@1/latest/currencies/usd/idr.json"
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req = req.WithContext(ctx)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return result, err
		}
		defer resp.Body.Close()

		var currResp dto.CurrencyNowAPI
		if err := json.NewDecoder(resp.Body).Decode(&currResp); err != nil {
			return result, err
		}

		result.IDR = int(math.Ceil(float64(payload.USD) * float64(currResp.Idr)))
	} else {
		result.IDR = int(math.Ceil(float64(payload.USD) * float64(currency.Value)))
	}

	return result, nil
}
