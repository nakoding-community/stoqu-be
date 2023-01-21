package usecase

import (
	"context"
	"reflect"
	"testing"

	"gitlab.com/stoqu/stoqu-be/internal/config"
	"gitlab.com/stoqu/stoqu-be/internal/factory/repository"
	"gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	"gitlab.com/stoqu/stoqu-be/internal/model/dto"
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"
)

func TestNewStock(t *testing.T) {
	type args struct {
		cfg *config.Configuration
		f   repository.Factory
	}
	tests := []struct {
		name string
		args args
		want Stock
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStock(tt.args.cfg, tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stock_Find(t *testing.T) {
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
		wantResult     []dto.StockViewResponse
		wantPagination abstraction.PaginationInfo
		wantErr        bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &stock{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			gotResult, gotPagination, err := u.Find(tt.args.ctx, tt.args.filterParam)
			if (err != nil) != tt.wantErr {
				t.Errorf("stock.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("stock.Find() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
			if !reflect.DeepEqual(gotPagination, tt.wantPagination) {
				t.Errorf("stock.Find() gotPagination = %v, want %v", gotPagination, tt.wantPagination)
			}
		})
	}
}

func Test_stock_FindByID(t *testing.T) {
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
		want    dto.StockResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &stock{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			got, err := u.FindByID(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("stock.FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("stock.FindByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_stock_Transaction(t *testing.T) {
	type fields struct {
		Cfg  *config.Configuration
		Repo repository.Factory
	}
	type args struct {
		ctx     context.Context
		payload dto.TransactionStockRequest
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult dto.StockTransactionResponse
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &stock{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			gotResult, err := u.Transaction(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("stock.Transaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("stock.Transaction() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_stock_Convertion(t *testing.T) {
	type fields struct {
		Cfg  *config.Configuration
		Repo repository.Factory
	}
	type args struct {
		ctx     context.Context
		payload dto.ConvertionStockRequest
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult dto.StockConvertionResponse
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &stock{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			gotResult, err := u.Convertion(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("stock.Convertion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("stock.Convertion() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_stock_Movement(t *testing.T) {
	type fields struct {
		Cfg  *config.Configuration
		Repo repository.Factory
	}
	type args struct {
		ctx     context.Context
		payload dto.MovementStockRequest
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult dto.StockMovementResponse
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &stock{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			gotResult, err := u.Movement(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("stock.Movement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("stock.Movement() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_stock_History(t *testing.T) {
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
		wantResult     []dto.StockHistoryResponse
		wantPagination abstraction.PaginationInfo
		wantErr        bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &stock{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			gotResult, gotPagination, err := u.History(tt.args.ctx, tt.args.filterParam)
			if (err != nil) != tt.wantErr {
				t.Errorf("stock.History() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("stock.History() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
			if !reflect.DeepEqual(gotPagination, tt.wantPagination) {
				t.Errorf("stock.History() gotPagination = %v, want %v", gotPagination, tt.wantPagination)
			}
		})
	}
}

func Test_stock_upsertStockRack(t *testing.T) {
	type fields struct {
		Cfg  *config.Configuration
		Repo repository.Factory
	}
	type args struct {
		ctx     context.Context
		stockID string
		rackID  string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult *model.StockRackModel
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &stock{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			gotResult, err := u.upsertStockRack(tt.args.ctx, tt.args.stockID, tt.args.rackID)
			if (err != nil) != tt.wantErr {
				t.Errorf("stock.upsertStockRack() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("stock.upsertStockRack() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_stock_TransactionProcess(t *testing.T) {
	type fields struct {
		Cfg  *config.Configuration
		Repo repository.Factory
	}
	type args struct {
		ctx     context.Context
		payload dto.TransactionStockRequest
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult dto.StockTransactionResponse
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &stock{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			gotResult, err := u.TransactionProcess(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("stock.TransactionProcess() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("stock.TransactionProcess() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_stock_transactionTypeIn(t *testing.T) {
	type fields struct {
		Cfg  *config.Configuration
		Repo repository.Factory
	}
	type args struct {
		ctx     context.Context
		wrapper *transactionDataWrapper
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
			u := &stock{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			if err := u.transactionTypeIn(tt.args.ctx, tt.args.wrapper); (err != nil) != tt.wantErr {
				t.Errorf("stock.transactionTypeIn() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_stock_transactionTypeOut(t *testing.T) {
	type fields struct {
		Cfg  *config.Configuration
		Repo repository.Factory
	}
	type args struct {
		ctx     context.Context
		wrapper *transactionDataWrapper
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
			u := &stock{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			if err := u.transactionTypeOut(tt.args.ctx, tt.args.wrapper); (err != nil) != tt.wantErr {
				t.Errorf("stock.transactionTypeOut() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_stock_convertionValidateData(t *testing.T) {
	type fields struct {
		Cfg  *config.Configuration
		Repo repository.Factory
	}
	type args struct {
		ctx     context.Context
		payload dto.ConvertionStockRequest
		wrapper *convertionDataWrapper
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
			u := &stock{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			if err := u.convertionValidateData(tt.args.ctx, tt.args.payload, tt.args.wrapper); (err != nil) != tt.wantErr {
				t.Errorf("stock.convertionValidateData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_stock_convertionValidateCalculation(t *testing.T) {
	type fields struct {
		Cfg  *config.Configuration
		Repo repository.Factory
	}
	type args struct {
		ctx     context.Context
		payload dto.ConvertionStockRequest
		wrapper *convertionDataWrapper
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
			u := &stock{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			if err := u.convertionValidateCalculation(tt.args.ctx, tt.args.payload, tt.args.wrapper); (err != nil) != tt.wantErr {
				t.Errorf("stock.convertionValidateCalculation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_stock_convertionProcess(t *testing.T) {
	type fields struct {
		Cfg  *config.Configuration
		Repo repository.Factory
	}
	type args struct {
		ctx     context.Context
		payload dto.ConvertionStockRequest
		wrapper *convertionDataWrapper
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
			u := &stock{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			if err := u.convertionProcess(tt.args.ctx, tt.args.payload, tt.args.wrapper); (err != nil) != tt.wantErr {
				t.Errorf("stock.convertionProcess() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_stock_convertionMutation(t *testing.T) {
	type fields struct {
		Cfg  *config.Configuration
		Repo repository.Factory
	}
	type args struct {
		ctx     context.Context
		wrapper *convertionDataWrapper
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
			u := &stock{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			if err := u.convertionMutation(tt.args.ctx, tt.args.wrapper); (err != nil) != tt.wantErr {
				t.Errorf("stock.convertionMutation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
