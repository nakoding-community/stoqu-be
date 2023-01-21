package usecase

import (
	"context"
	"reflect"
	"testing"

	"gitlab.com/stoqu/stoqu-be/internal/config"
	"gitlab.com/stoqu/stoqu-be/internal/factory/repository"
	"gitlab.com/stoqu/stoqu-be/internal/model/dto"
)

func TestNewDashboard(t *testing.T) {
	type args struct {
		cfg *config.Configuration
		f   repository.Factory
	}
	tests := []struct {
		name string
		args args
		want Dashboard
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDashboard(tt.args.cfg, tt.args.f); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDashboard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dashboard_Count(t *testing.T) {
	type fields struct {
		Cfg  *config.Configuration
		Repo repository.Factory
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantResult dto.DashboardResponse
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &dashboard{
				Cfg:  tt.fields.Cfg,
				Repo: tt.fields.Repo,
			}
			gotResult, err := u.Count(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("dashboard.Count() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("dashboard.Count() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
