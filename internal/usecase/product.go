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

type Product interface {
	Find(ctx context.Context, filterParam abstraction.Filter) ([]dto.ProductViewResponse, abstraction.PaginationInfo, error)
	FindByID(ctx context.Context, payload dto.ByIDRequest) (dto.ProductResponse, error)
	Create(ctx context.Context, payload dto.CreateProductRequest) (dto.ProductResponse, error)
	Update(ctx context.Context, payload dto.UpdateProductRequest) (dto.ProductResponse, error)
	Delete(ctx context.Context, payload dto.ByIDRequest) (dto.ProductResponse, error)
}

type product struct {
	Repo repository.Factory
	Cfg  *config.Configuration
}

func NewProduct(cfg *config.Configuration, f repository.Factory) Product {
	return &product{f, cfg}
}

func (u *product) Find(ctx context.Context, filterParam abstraction.Filter) (result []dto.ProductViewResponse, pagination abstraction.PaginationInfo, err error) {
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

	products, info, err := u.Repo.Product.Find(ctx, filterParam, search)
	if err != nil {
		return nil, pagination, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}
	pagination = *info

	for _, product := range products {
		result = append(result, dto.ProductViewResponse{
			ProductView: product,
		})
	}

	return result, pagination, nil
}

func (u *product) FindByID(ctx context.Context, payload dto.ByIDRequest) (dto.ProductResponse, error) {
	var result dto.ProductResponse

	product, err := u.Repo.Product.FindByID(ctx, payload.ID)
	if err != nil {
		return result, err
	}

	result = dto.ProductResponse{
		ProductModel: *product,
	}

	return result, nil
}

func (u *product) Create(ctx context.Context, payload dto.CreateProductRequest) (result dto.ProductResponse, err error) {
	var (
		productID = uuid.New().String()
		product   = model.ProductModel{
			Entity: model.Entity{
				ID: productID,
			},
			ProductEntity: model.ProductEntity{
				Code:       str.GenCode(constant.CODE_PRODUCT_PREFIX),
				Name:       payload.Name,
				PriceUSD:   payload.PriceUSD,
				PriceFinal: payload.PriceFinal,
				BrandID:    payload.BrandID,
				VariantID:  payload.VariantID,
				PacketID:   payload.PacketID,
			},
		}
	)

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		_, err = u.Repo.Product.Create(ctx, product)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.ProductResponse{
		ProductModel: product,
	}

	return result, nil
}

func (u *product) Update(ctx context.Context, payload dto.UpdateProductRequest) (result dto.ProductResponse, err error) {
	var (
		productData = &model.ProductModel{
			ProductEntity: model.ProductEntity{
				Name:       payload.Name,
				PriceUSD:   payload.PriceUSD,
				PriceFinal: payload.PriceFinal,
			},
		}
	)

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		_, err = u.Repo.Product.UpdateByID(ctx, payload.ID, *productData)
		if err != nil {
			return err
		}

		productData, err = u.Repo.Product.FindByID(ctx, payload.ID)
		if err != nil {
			return res.ErrorBuilder(res.Constant.Error.NotFound, err)
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.ProductResponse{
		ProductModel: *productData,
	}

	return result, nil
}

func (u *product) Delete(ctx context.Context, payload dto.ByIDRequest) (result dto.ProductResponse, err error) {
	var data *model.ProductModel

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		data, err = u.Repo.Product.FindByID(ctx, payload.ID)
		if err != nil {
			return err
		}

		err = u.Repo.Product.DeleteByID(ctx, payload.ID)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.ProductResponse{
		ProductModel: *data,
	}

	return result, nil
}
