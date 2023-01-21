package usecase

import (
	"context"
	"reflect"
	"testing"

	"gitlab.com/stoqu/stoqu-be/internal/config"
	"gitlab.com/stoqu/stoqu-be/internal/factory/repository"
	"gitlab.com/stoqu/stoqu-be/internal/model/dto"
)

func TestNewAuth(t *testing.T) {
	type args struct {
		cfg *config.Configuration
		f   repository.Factory
	}
	tests := []struct {
		name string
		args args
		want Auth
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAuth(tt.args.cfg, tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_auth_Login(t *testing.T) {
	type fields struct {
		Cfg  *config.Configuration
		Repo repository.Factory
	}
	type args struct {
		ctx     context.Context
		payload dto.LoginAuthRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    dto.AuthLoginResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &auth{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			got, err := u.Login(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("auth.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("auth.Login() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_auth_Register(t *testing.T) {
	type fields struct {
		Cfg  *config.Configuration
		Repo repository.Factory
	}
	type args struct {
		ctx     context.Context
		payload dto.RegisterAuthRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    dto.AuthRegisterResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &auth{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			got, err := u.Register(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("auth.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("auth.Register() = %v, want %v", got, tt.want)
			}
		})
	}
}
