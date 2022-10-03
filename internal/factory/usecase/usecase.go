package usecase

import (
	"gitlab.com/stoqu/stoqu-be/internal/config"
	"gitlab.com/stoqu/stoqu-be/internal/factory/repository"
	"gitlab.com/stoqu/stoqu-be/internal/usecase"
)

type Factory struct {
	Auth usecase.Auth
	User usecase.User
}

func Init(cfg *config.Configuration, r repository.Factory) Factory {
	f := Factory{}
	f.Auth = usecase.NewAuth(cfg, r)
	f.User = usecase.NewUser(cfg, r)

	return f
}
