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

func TestNewCurrency(t *testing.T) {
	type args struct {
		cfg *config.Configuration
		f   repository.Factory
	}
	tests := []struct {
		name string
		args args
		want Currency
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCurrency(tt.args.cfg, tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCurrency() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_currency_Find(t *testing.T) {
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
		wantResult     []dto.CurrencyResponse
		wantPagination abstraction.PaginationInfo
		wantErr        bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &currency{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			gotResult, gotPagination, err := u.Find(tt.args.ctx, tt.args.filterParam)
			if (err != nil) != tt.wantErr {
				t.Errorf("currency.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("currency.Find() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
			if !reflect.DeepEqual(gotPagination, tt.wantPagination) {
				t.Errorf("currency.Find() gotPagination = %v, want %v", gotPagination, tt.wantPagination)
			}
		})
	}
}

func Test_currency_FindByID(t *testing.T) {
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
		want    dto.CurrencyResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &currency{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			got, err := u.FindByID(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("currency.FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("currency.FindByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_currency_Create(t *testing.T) {
	type fields struct {
		Cfg  *config.Configuration
		Repo repository.Factory
	}
	type args struct {
		ctx     context.Context
		payload dto.CreateCurrencyRequest
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult dto.CurrencyResponse
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &currency{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			gotResult, err := u.Create(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("currency.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("currency.Create() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_currency_Update(t *testing.T) {
	type fields struct {
		Cfg  *config.Configuration
		Repo repository.Factory
	}
	type args struct {
		ctx     context.Context
		payload dto.UpdateCurrencyRequest
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult dto.CurrencyResponse
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &currency{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			gotResult, err := u.Update(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("currency.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("currency.Update() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_currency_Delete(t *testing.T) {
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
		wantResult dto.CurrencyResponse
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &currency{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			gotResult, err := u.Delete(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("currency.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("currency.Delete() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_currency_Convert(t *testing.T) {
	type fields struct {
		Cfg  *config.Configuration
		Repo repository.Factory
	}
	type args struct {
		ctx     context.Context
		payload dto.ConvertCurrencyRequest
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult dto.ConvertCurrencyResponse
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &currency{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			gotResult, err := u.Convert(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("currency.Convert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("currency.Convert() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
