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
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"
	"gitlab.com/stoqu/stoqu-be/pkg/util/trxmanager"
)

type Stock interface {
	Find(ctx context.Context, filterParam abstraction.Filter) ([]dto.StockViewResponse, abstraction.PaginationInfo, error)
	FindByID(ctx context.Context, payload dto.ByIDRequest) (dto.StockResponse, error)
	Create(ctx context.Context, payload dto.CreateStockRequest) (dto.StockResponse, error)
}

type stock struct {
	Repo repository.Factory
	Cfg  *config.Configuration
}

func NewStock(cfg *config.Configuration, f repository.Factory) Stock {
	return &stock{f, cfg}
}

func (u *stock) Find(ctx context.Context, filterParam abstraction.Filter) (result []dto.StockViewResponse, pagination abstraction.PaginationInfo, err error) {
	var search *abstraction.Search
	if filterParam.Search != "" {
		searchQuery := `
			lower(products.code) LIKE ? OR 
			lower(products.name) LIKE ? OR 
			products.price_usd = ? OR 
			products.price_final = ? OR 
			lower(brands.name) LIKE ? OR 
			lower(users.name) LIKE ? OR 
			lower(variants.name) LIKE ? OR 
			lower(units.name) LIKE ? OR
			packets.value = ?
		`
		searchVal := "%" + strings.ToLower(filterParam.Search) + "%"
		searchValFloat, _ := strconv.ParseFloat(filterParam.Search, 64)
		search = &abstraction.Search{
			Query: searchQuery,
			Args: []interface{}{
				searchVal,
				searchVal,
				searchValFloat,
				searchValFloat,
				searchVal,
				searchVal,
				searchVal,
				searchVal,
				searchValFloat,
			},
		}
	}

	products, info, err := u.Repo.Stock.Find(ctx, filterParam, search)
	if err != nil {
		return nil, pagination, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}
	pagination = *info

	for _, product := range products {
		result = append(result, dto.StockViewResponse{
			StockView: product,
		})
	}

	return result, pagination, nil
}

func (u *stock) FindByID(ctx context.Context, payload dto.ByIDRequest) (dto.StockResponse, error) {
	var result dto.StockResponse

	product, err := u.Repo.Stock.FindByID(ctx, payload.ID)
	if err != nil {
		return result, err
	}

	result = dto.StockResponse{
		StockModel: *product,
	}

	return result, nil
}

func (u *stock) Create(ctx context.Context, payload dto.CreateStockRequest) (result dto.StockResponse, err error) {
	var (
		productID = uuid.New().String()
		product   = model.StockModel{
			Entity: model.Entity{
				ID: productID,
			},
			StockEntity: model.StockEntity{
				ProductID: payload.ProductID,
			},
		}
	)

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		_, err = u.Repo.Stock.Create(ctx, product)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.StockResponse{
		StockModel: product,
	}

	return result, nil
}
