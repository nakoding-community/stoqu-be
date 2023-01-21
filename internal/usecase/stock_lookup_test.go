package usecase

import (
	"context"
	"reflect"
	"testing"

	"gitlab.com/stoqu/stoqu-be/internal/config"
	"gitlab.com/stoqu/stoqu-be/internal/factory/repository"
	"gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	"gitlab.com/stoqu/stoqu-be/internal/model/dto"
)

func TestNewStockLookup(t *testing.T) {
	type args struct {
		cfg *config.Configuration
		f   repository.Factory
	}
	tests := []struct {
		name string
		args args
		want StockLookup
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStockLookup(tt.args.cfg, tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStockLookup() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stockLookup_Find(t *testing.T) {
	type fields struct {
		Cfg  *config.Configuration
		Repo repository.Factory
	}
	type args struct {
		ctx         context.Context
		filterParam abstraction.Filter
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantResult     []dto.StockLookupViewResponse
		wantPagination abstraction.PaginationInfo
		wantErr        bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &stockLookup{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			gotResult, gotPagination, err := u.Find(tt.args.ctx, tt.args.filterParam)
			if (err != nil) != tt.wantErr {
				t.Errorf("stockLookup.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("stockLookup.Find() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
			if !reflect.DeepEqual(gotPagination, tt.wantPagination) {
				t.Errorf("stockLookup.Find() gotPagination = %v, want %v", gotPagination, tt.wantPagination)
			}
		})
	}
}

func Test_stockLookup_FindByID(t *testing.T) {
	type fields struct {
		Cfg  *config.Configuration
		Repo repository.Factory
	}
	type args struct {
		ctx     context.Context
		payload dto.ByIDRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    dto.StockLookupResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &stockLookup{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			got, err := u.FindByID(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("stockLookup.FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("stockLookup.FindByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stockLookup_Create(t *testing.T) {
	type fields struct {
		Cfg  *config.Configuration
		Repo repository.Factory
	}
	type args struct {
		ctx     context.Context
		payload dto.CreateStockLookupRequest
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult dto.StockLookupResponse
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &stockLookup{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			gotResult, err := u.Create(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("stockLookup.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("stockLookup.Create() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_stockLookup_Update(t *testing.T) {
	type fields struct {
		Cfg  *config.Configuration
		Repo repository.Factory
	}
	type args struct {
		ctx     context.Context
		payload dto.UpdateStockLookupRequest
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult dto.StockLookupResponse
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &stockLookup{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			gotResult, err := u.Update(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("stockLookup.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("stockLookup.Update() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_stockLookup_Delete(t *testing.T) {
	type fields struct {
		Cfg  *config.Configuration
		Repo repository.Factory
	}
	type args struct {
		ctx     context.Context
		payload dto.ByIDRequest
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult dto.StockLookupResponse
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &stockLookup{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			gotResult, err := u.Delete(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("stockLookup.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("stockLookup.Delete() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
