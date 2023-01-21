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

func TestNewProduct(t *testing.T) {
	type args struct {
		cfg *config.Configuration
		f   repository.Factory
	}
	tests := []struct {
		name string
		args args
		want Product
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewProduct(tt.args.cfg, tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProduct() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_product_Find(t *testing.T) {
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
		wantResult     []dto.ProductViewResponse
		wantPagination abstraction.PaginationInfo
		wantErr        bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &product{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			gotResult, gotPagination, err := u.Find(tt.args.ctx, tt.args.filterParam)
			if (err != nil) != tt.wantErr {
				t.Errorf("product.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("product.Find() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
			if !reflect.DeepEqual(gotPagination, tt.wantPagination) {
				t.Errorf("product.Find() gotPagination = %v, want %v", gotPagination, tt.wantPagination)
			}
		})
	}
}

func Test_product_FindByID(t *testing.T) {
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
		want    dto.ProductResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &product{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			got, err := u.FindByID(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("product.FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("product.FindByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_product_Create(t *testing.T) {
	type fields struct {
		Cfg  *config.Configuration
		Repo repository.Factory
	}
	type args struct {
		ctx     context.Context
		payload dto.CreateProductRequest
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult dto.ProductResponse
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &product{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			gotResult, err := u.Create(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("product.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("product.Create() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_product_Update(t *testing.T) {
	type fields struct {
		Cfg  *config.Configuration
		Repo repository.Factory
	}
	type args struct {
		ctx     context.Context
		payload dto.UpdateProductRequest
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult dto.ProductResponse
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &product{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			gotResult, err := u.Update(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("product.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("product.Update() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_product_Delete(t *testing.T) {
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
		wantResult dto.ProductResponse
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &product{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			gotResult, err := u.Delete(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("product.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("product.Delete() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
