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

func TestNewUser(t *testing.T) {
	type args struct {
		cfg *config.Configuration
		f   repository.Factory
	}
	tests := []struct {
		name string
		args args
		want User
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUser(tt.args.cfg, tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_user_Find(t *testing.T) {
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
		wantResult     []dto.UserResponse
		wantPagination abstraction.PaginationInfo
		wantErr        bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &user{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			gotResult, gotPagination, err := u.Find(tt.args.ctx, tt.args.filterParam)
			if (err != nil) != tt.wantErr {
				t.Errorf("user.Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("user.Find() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
			if !reflect.DeepEqual(gotPagination, tt.wantPagination) {
				t.Errorf("user.Find() gotPagination = %v, want %v", gotPagination, tt.wantPagination)
			}
		})
	}
}

func Test_user_FindByID(t *testing.T) {
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
		want    dto.UserResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &user{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			got, err := u.FindByID(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("user.FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("user.FindByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_user_Create(t *testing.T) {
	type fields struct {
		Cfg  *config.Configuration
		Repo repository.Factory
	}
	type args struct {
		ctx     context.Context
		payload dto.CreateUserRequest
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult dto.UserResponse
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &user{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			gotResult, err := u.Create(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("user.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("user.Create() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_user_Update(t *testing.T) {
	type fields struct {
		Cfg  *config.Configuration
		Repo repository.Factory
	}
	type args struct {
		ctx     context.Context
		payload dto.UpdateUserRequest
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult dto.UserResponse
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &user{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			gotResult, err := u.Update(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("user.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("user.Update() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func Test_user_Delete(t *testing.T) {
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
		wantResult dto.UserResponse
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &user{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			gotResult, err := u.Delete(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("user.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("user.Delete() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
