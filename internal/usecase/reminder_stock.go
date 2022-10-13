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

type ReminderStock interface {
	Find(ctx context.Context, filterParam abstraction.Filter) ([]dto.ReminderStockResponse, abstraction.PaginationInfo, error)
	FindByID(ctx context.Context, payload dto.ByIDRequest) (dto.ReminderStockResponse, error)
	Create(ctx context.Context, payload dto.CreateReminderStockRequest) (dto.ReminderStockResponse, error)
	Update(ctx context.Context, payload dto.UpdateReminderStockRequest) (dto.ReminderStockResponse, error)
	Delete(ctx context.Context, payload dto.ByIDRequest) (dto.ReminderStockResponse, error)
}

type reminderStock struct {
	Repo repository.Factory
	Cfg  *config.Configuration
}

func NewReminderStock(cfg *config.Configuration, f repository.Factory) ReminderStock {
	return &reminderStock{f, cfg}
}

func (u *reminderStock) Find(ctx context.Context, filterParam abstraction.Filter) (result []dto.ReminderStockResponse, pagination abstraction.PaginationInfo, err error) {
	var search *abstraction.Search
	if filterParam.Search != "" {
		searchQuery := "lower(code) LIKE ? OR lower(name) LIKE ?"
		searchVal := "%" + strings.ToLower(filterParam.Search) + "%"
		search = &abstraction.Search{
			Query: searchQuery,
			Args:  []interface{}{searchVal, searchVal},
		}
	}

	reminderStocks, info, err := u.Repo.ReminderStock.Find(ctx, filterParam, search)
	if err != nil {
		return nil, pagination, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}
	pagination = *info

	for _, reminderStock := range reminderStocks {
		result = append(result, dto.ReminderStockResponse{
			ReminderStockModel: reminderStock,
		})
	}

	return result, pagination, nil
}

func (u *reminderStock) FindByID(ctx context.Context, payload dto.ByIDRequest) (dto.ReminderStockResponse, error) {
	var result dto.ReminderStockResponse

	reminderStock, err := u.Repo.ReminderStock.FindByID(ctx, payload.ID)
	if err != nil {
		return result, err
	}

	result = dto.ReminderStockResponse{
		ReminderStockModel: *reminderStock,
	}

	return result, nil
}

func (u *reminderStock) Create(ctx context.Context, payload dto.CreateReminderStockRequest) (result dto.ReminderStockResponse, err error) {
	var (
		reminderStockID = uuid.New().String()
		reminderStock   = model.ReminderStockModel{
			Entity: model.Entity{
				ID: reminderStockID,
			},
			ReminderStockEntity: model.ReminderStockEntity{
				Code:     str.GenCode(constant.CODE_REMINDER_STOCK_PREFIX),
				Name:     payload.Name,
				MinStock: payload.MinStock,
			},
		}
	)

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		_, err = u.Repo.ReminderStock.Create(ctx, reminderStock)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.ReminderStockResponse{
		ReminderStockModel: reminderStock,
	}

	return result, nil
}

func (u *reminderStock) Update(ctx context.Context, payload dto.UpdateReminderStockRequest) (result dto.ReminderStockResponse, err error) {
	var (
		reminderStockData = &model.ReminderStockModel{
			ReminderStockEntity: model.ReminderStockEntity{
				Name:     payload.Name,
				MinStock: payload.MinStock,
			},
		}
	)

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		_, err = u.Repo.ReminderStock.UpdateByID(ctx, payload.ID, *reminderStockData)
		if err != nil {
			return err
		}

		reminderStockData, err = u.Repo.ReminderStock.FindByID(ctx, payload.ID)
		if err != nil {
			return res.ErrorBuilder(res.Constant.Error.NotFound, err)
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.ReminderStockResponse{
		ReminderStockModel: *reminderStockData,
	}

	return result, nil
}

func (u *reminderStock) Delete(ctx context.Context, payload dto.ByIDRequest) (result dto.ReminderStockResponse, err error) {
	var data *model.ReminderStockModel

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		data, err = u.Repo.ReminderStock.FindByID(ctx, payload.ID)
		if err != nil {
			return err
		}

		err = u.Repo.ReminderStock.DeleteByID(ctx, payload.ID)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.ReminderStockResponse{
		ReminderStockModel: *data,
	}

	return result, nil
}
