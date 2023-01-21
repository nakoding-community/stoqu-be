package usecase

import (
	"context"
	"reflect"
	"testing"

	"gitlab.com/stoqu/stoqu-be/internal/config"
	"gitlab.com/stoqu/stoqu-be/internal/factory/repository"
	"gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	"gitlab.com/stoqu/stoqu-be/internal/model/dto"
	"gitlab.com/stoqu/stoqu-be/internal/model/entity"
)

func TestNewOrder(t *testing.T) {
	type args struct {
		cfg     *config.Configuration
		f       repository.Factory
		stockUC Stock
	}
	tests := []struct {
		name string
		args args
		want Order
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOrder(tt.args.cfg, tt.args.f, tt.args.stockUC); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_order_Find(t *testing.T) {
	type fields struct {
		Cfg     *config.Configuration
		Repo    repository.Factory
		stockUC Stock
	}
	type args struct {
		ctx         context.Context
		filterParam abstraction.Filter
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantResult     []dto.OrderViewResponse
		wantPagination abstraction.PaginationInfo
		wantErr        bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &order{
				Cfg:     tt.fields.Cfg,
				Repo:    tt.fields.Repo,
				stockUC: tt.fields.stockUC,
			}
			gotResult, gotPagination, err := u.Find(tt.args.ctx, tt.args.filterParam)
			if (err != nil) != tt.wantErr {
				t.Errorf("order.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("order.Find() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
			if !reflect.DeepEqual(gotPagination, tt.wantPagination) {
				t.Errorf("order.Find() gotPagination = %v, want %v", gotPagination, tt.wantPagination)
			}
		})
	}
}

func Test_order_FindByID(t *testing.T) {
	type fields struct {
		Cfg     *config.Configuration
		Repo    repository.Factory
		stockUC Stock
	}
	type args struct {
		ctx     context.Context
		payload dto.ByIDRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    dto.OrderViewResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &order{
				Cfg:     tt.fields.Cfg,
				Repo:    tt.fields.Repo,
				stockUC: tt.fields.stockUC,
			}
			got, err := u.FindByID(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("order.FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("order.FindByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_order_FindDetailByID(t *testing.T) {
	type fields struct {
		Cfg     *config.Configuration
		Repo    repository.Factory
		stockUC Stock
	}
	type args struct {
		ctx     context.Context
		payload dto.ByIDRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    dto.OrderViewDetailResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &order{
				Cfg:     tt.fields.Cfg,
				Repo:    tt.fields.Repo,
				stockUC: tt.fields.stockUC,
			}
			got, err := u.FindDetailByID(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("order.FindDetailByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("order.FindDetailByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_order_Upsert(t *testing.T) {
	type fields struct {
		Cfg     *config.Configuration
		Repo    repository.Factory
		stockUC Stock
	}
	type args struct {
		ctx     context.Context
		payload dto.UpsertOrderRequest
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult dto.OrderUpsertResponse
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &order{
				Cfg:     tt.fields.Cfg,
				Repo:    tt.fields.Repo,
				stockUC: tt.fields.stockUC,
			}
			gotResult, err := u.Upsert(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("order.Upsert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("order.Upsert() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_order_buildMapStockLookups(t *testing.T) {
	type fields struct {
		Cfg     *config.Configuration
		Repo    repository.Factory
		stockUC Stock
	}
	type args struct {
		ctx     context.Context
		payload dto.UpsertOrderRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]entity.StockLookupModel
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &order{
				Cfg:     tt.fields.Cfg,
				Repo:    tt.fields.Repo,
				stockUC: tt.fields.stockUC,
			}
			got, err := u.buildMapStockLookups(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("order.buildMapStockLookups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("order.buildMapStockLookups() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_order_create(t *testing.T) {
	type fields struct {
		Cfg     *config.Configuration
		Repo    repository.Factory
		stockUC Stock
	}
	type args struct {
		ctx      context.Context
		payload  dto.UpsertOrderRequest
		orderTrx entity.OrderTrxModel
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &order{
				Cfg:     tt.fields.Cfg,
				Repo:    tt.fields.Repo,
				stockUC: tt.fields.stockUC,
			}
			if err := u.create(tt.args.ctx, tt.args.payload, tt.args.orderTrx); (err != nil) != tt.wantErr {
				t.Errorf("order.create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_order_update(t *testing.T) {
	type fields struct {
		Cfg     *config.Configuration
		Repo    repository.Factory
		stockUC Stock
	}
	type args struct {
		ctx      context.Context
		payload  dto.UpsertOrderRequest
		orderTrx entity.OrderTrxModel
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &order{
				Cfg:     tt.fields.Cfg,
				Repo:    tt.fields.Repo,
				stockUC: tt.fields.stockUC,
			}
			if err := u.update(tt.args.ctx, tt.args.payload, tt.args.orderTrx); (err != nil) != tt.wantErr {
				t.Errorf("order.update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
