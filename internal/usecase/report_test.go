package usecase

import (
	"context"
	"os"
	"reflect"
	"testing"

	"gitlab.com/stoqu/stoqu-be/internal/config"
	"gitlab.com/stoqu/stoqu-be/internal/factory/repository"
	"gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	"gitlab.com/stoqu/stoqu-be/internal/model/dto"
	"gitlab.com/stoqu/stoqu-be/internal/model/entity"
)

func TestNewReport(t *testing.T) {
	type args struct {
		cfg *config.Configuration
		f   repository.Factory
	}
	tests := []struct {
		name string
		args args
		want Report
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewReport(tt.args.cfg, tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewReport() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_report_FindOrder(t *testing.T) {
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
		wantResult     dto.OrderReportResponse
		wantPagination abstraction.PaginationInfo
		wantErr        bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &report{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			gotResult, gotPagination, err := u.FindOrder(tt.args.ctx, tt.args.filterParam)
			if (err != nil) != tt.wantErr {
				t.Errorf("report.FindOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("report.FindOrder() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
			if !reflect.DeepEqual(gotPagination, tt.wantPagination) {
				t.Errorf("report.FindOrder() gotPagination = %v, want %v", gotPagination, tt.wantPagination)
			}
		})
	}
}

func Test_report_orderExcelData(t *testing.T) {
	type fields struct {
		Cfg  *config.Configuration
		Repo repository.Factory
	}
	type args struct {
		ctx         context.Context
		filterParam abstraction.Filter
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []entity.OrderView
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &report{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			got, err := u.orderExcelData(tt.args.ctx, tt.args.filterParam)
			if (err != nil) != tt.wantErr {
				t.Errorf("report.orderExcelData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("report.orderExcelData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_report_FindOrderExcel(t *testing.T) {
	type fields struct {
		Cfg  *config.Configuration
		Repo repository.Factory
	}
	type args struct {
		ctx         context.Context
		filterParam abstraction.Filter
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantF   *os.File
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &report{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			gotF, err := u.FindOrderExcel(tt.args.ctx, tt.args.filterParam)
			if (err != nil) != tt.wantErr {
				t.Errorf("report.FindOrderExcel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotF, tt.wantF) {
				t.Errorf("report.FindOrderExcel() = %v, want %v", gotF, tt.wantF)
			}
		})
	}
}

func Test_report_FindOrderProduct(t *testing.T) {
	type fields struct {
		Cfg  *config.Configuration
		Repo repository.Factory
	}
	type args struct {
		ctx         context.Context
		filterParam abstraction.Filter
		query       dto.ProductReportQuery
	}
	tests := []struct {
		name           string
		fields         fields
		args           args
		wantResult     dto.OrderProductReportResponse
		wantPagination abstraction.PaginationInfo
		wantErr        bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &report{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			gotResult, gotPagination, err := u.FindOrderProduct(tt.args.ctx, tt.args.filterParam, tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("report.FindOrderProduct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("report.FindOrderProduct() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
			if !reflect.DeepEqual(gotPagination, tt.wantPagination) {
				t.Errorf("report.FindOrderProduct() gotPagination = %v, want %v", gotPagination, tt.wantPagination)
			}
		})
	}
}

func Test_report_orderProductExcelData(t *testing.T) {
	type fields struct {
		Cfg  *config.Configuration
		Repo repository.Factory
	}
	type args struct {
		ctx         context.Context
		filterParam abstraction.Filter
		query       dto.ProductReportQuery
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantOrders []entity.OrderViewProduct
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &report{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			gotOrders, err := u.orderProductExcelData(tt.args.ctx, tt.args.filterParam, tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("report.orderProductExcelData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOrders, tt.wantOrders) {
				t.Errorf("report.orderProductExcelData() = %v, want %v", gotOrders, tt.wantOrders)
			}
		})
	}
}

func Test_report_FindOrderProductExcel(t *testing.T) {
	type fields struct {
		Cfg  *config.Configuration
		Repo repository.Factory
	}
	type args struct {
		ctx         context.Context
		filterParam abstraction.Filter
		query       dto.ProductReportQuery
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantF   *os.File
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &report{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			gotF, err := u.FindOrderProductExcel(tt.args.ctx, tt.args.filterParam, tt.args.query)
			if (err != nil) != tt.wantErr {
				t.Errorf("report.FindOrderProductExcel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotF, tt.wantF) {
				t.Errorf("report.FindOrderProductExcel() = %v, want %v", gotF, tt.wantF)
			}
		})
	}
}

func Test_report_GenerateExcelReport(t *testing.T) {
	type fields struct {
		Cfg  *config.Configuration
		Repo repository.Factory
	}
	type args struct {
		ctx   context.Context
		excel dto.GenerateExcelReportInput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantF   *os.File
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &report{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			gotF, err := u.GenerateExcelReport(tt.args.ctx, tt.args.excel)
			if (err != nil) != tt.wantErr {
				t.Errorf("report.GenerateExcelReport() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotF, tt.wantF) {
				t.Errorf("report.GenerateExcelReport() = %v, want %v", gotF, tt.wantF)
			}
		})
	}
}
