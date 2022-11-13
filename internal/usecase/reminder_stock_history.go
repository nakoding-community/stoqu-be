package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gitlab.com/stoqu/stoqu-be/internal/config"

	"gitlab.com/stoqu/stoqu-be/internal/factory/repository"
	"gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	"gitlab.com/stoqu/stoqu-be/internal/model/dto"
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"
	"gitlab.com/stoqu/stoqu-be/pkg/constant"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"
	"gitlab.com/stoqu/stoqu-be/pkg/util/trxmanager"
)

type ReminderStockHistory interface {
	Find(ctx context.Context, filterParam abstraction.Filter) ([]dto.ReminderStockHistoryResponse, abstraction.PaginationInfo, error)
	FindByID(ctx context.Context, payload dto.ByIDRequest) (dto.ReminderStockHistoryResponse, error)
	Update(ctx context.Context, payload dto.UpdateReminderStockHistoryRequest) (dto.ReminderStockHistoryResponse, error)
	Delete(ctx context.Context, payload dto.ByIDRequest) (dto.ReminderStockHistoryResponse, error)

	// Custom
	UpdateBulkRead(ctx context.Context, payload dto.UpdateReminderStockHistoryBulkReadRequest) (result dto.ReminderStockHistoryBulkReadResponse, err error)
	CountUnread(ctx context.Context) (result dto.ReminderStockHistoryCountUnreadResponse, err error)
	GenerateRecurring(ctx context.Context, reminderType string) error
}

type reminderStockHistory struct {
	Cfg  *config.Configuration
	Repo repository.Factory
}

func NewReminderStockHistory(cfg *config.Configuration, f repository.Factory) ReminderStockHistory {
	return &reminderStockHistory{cfg, f}
}

func (u *reminderStockHistory) Find(ctx context.Context, filterParam abstraction.Filter) (result []dto.ReminderStockHistoryResponse, pagination abstraction.PaginationInfo, err error) {
	var search *abstraction.Search
	if filterParam.Search != "" {
		searchQuery := "lower(title) LIKE ? OR lower(body) LIKE ?"
		searchVal := "%" + strings.ToLower(filterParam.Search) + "%"
		search = &abstraction.Search{
			Query: searchQuery,
			Args:  []interface{}{searchVal, searchVal},
		}
	}

	reminderStockHistorys, info, err := u.Repo.ReminderStockHistory.Find(ctx, filterParam, search)
	if err != nil {
		return nil, pagination, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}
	pagination = *info

	for _, reminderStockHistory := range reminderStockHistorys {
		result = append(result, dto.ReminderStockHistoryResponse{
			ReminderStockHistoryModel: reminderStockHistory,
		})
	}

	return result, pagination, nil
}

func (u *reminderStockHistory) FindByID(ctx context.Context, payload dto.ByIDRequest) (dto.ReminderStockHistoryResponse, error) {
	var result dto.ReminderStockHistoryResponse

	reminderStockHistory, err := u.Repo.ReminderStockHistory.FindByID(ctx, payload.ID)
	if err != nil {
		return result, err
	}

	result = dto.ReminderStockHistoryResponse{
		ReminderStockHistoryModel: *reminderStockHistory,
	}

	return result, nil
}

func (u *reminderStockHistory) Update(ctx context.Context, payload dto.UpdateReminderStockHistoryRequest) (result dto.ReminderStockHistoryResponse, err error) {
	var (
		reminderStockHistoryData = &model.ReminderStockHistoryModel{
			ReminderStockHistoryEntity: model.ReminderStockHistoryEntity{
				IsRead: payload.IsRead,
			},
		}
	)

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		_, err = u.Repo.ReminderStockHistory.UpdateByID(ctx, payload.ID, *reminderStockHistoryData)
		if err != nil {
			return err
		}

		reminderStockHistoryData, err = u.Repo.ReminderStockHistory.FindByID(ctx, payload.ID)
		if err != nil {
			return res.ErrorBuilder(res.Constant.Error.NotFound, err)
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.ReminderStockHistoryResponse{
		ReminderStockHistoryModel: *reminderStockHistoryData,
	}

	return result, nil
}

func (u *reminderStockHistory) Delete(ctx context.Context, payload dto.ByIDRequest) (result dto.ReminderStockHistoryResponse, err error) {
	var data *model.ReminderStockHistoryModel

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		data, err = u.Repo.ReminderStockHistory.FindByID(ctx, payload.ID)
		if err != nil {
			return err
		}

		err = u.Repo.ReminderStockHistory.DeleteByID(ctx, payload.ID)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.ReminderStockHistoryResponse{
		ReminderStockHistoryModel: *data,
	}

	return result, nil
}

func (u *reminderStockHistory) UpdateBulkRead(ctx context.Context, payload dto.UpdateReminderStockHistoryBulkReadRequest) (result dto.ReminderStockHistoryBulkReadResponse, err error) {
	err = u.Repo.ReminderStockHistory.UpdateBulkRead(ctx, payload.IDs)
	if err != nil {
		return result, err
	}

	result.Status = constant.STATUS_SUCCESS
	return result, nil
}

func (u *reminderStockHistory) CountUnread(ctx context.Context) (result dto.ReminderStockHistoryCountUnreadResponse, err error) {
	count, err := u.Repo.ReminderStockHistory.CountUnread(ctx)
	if err != nil {
		return result, err
	}
	result.Count = count

	return result, nil
}

func (u *reminderStockHistory) GenerateRecurring(ctx context.Context, reminderType string) error {
	reminderStocks, _, err := u.Repo.ReminderStock.Find(ctx, abstraction.Filter{}, &abstraction.Search{})
	if err != nil {
		return err
	}
	if len(reminderStocks) == 0 {
		return res.ErrorBuilder(res.Constant.Error.NotFound, errors.New("find reminder stock not found"))
	}

	reminderStock := reminderStocks[0]
	if reminderType != reminderStock.ReminderType {
		return nil
	}

	stocks, err := u.Repo.Stock.FindByTotalLessThan(ctx, reminderStock.MinStock)
	if err != nil {
		return err
	}

	histories := []model.ReminderStockHistoryModel{}
	for _, stock := range stocks {
		histories = append(histories, model.ReminderStockHistoryModel{
			ReminderStockHistoryEntity: model.ReminderStockHistoryEntity{
				Title: "Reminder Stock",
				Body:  fmt.Sprintf("Stok produk %s brand %s, tersisa %d", stock.ProductName, stock.BrandName, stock.Total),
			},
		})
	}
	_, err = u.Repo.ReminderStockHistory.Creates(ctx, histories)
	if err != nil {
		return err
	}

	return nil
}
