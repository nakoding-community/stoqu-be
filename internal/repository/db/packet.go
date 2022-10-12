package db

import (
	"context"

	abstraction "gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"

	"gorm.io/gorm"
)

type (
	Packet interface {
		// !TODO mockgen doesn't support embedded interface yet
		// !TODO but already discussed in this thread https://github.com/golang/mock/issues/621, lets wait for the release
		// Base[model.PacketModel]

		// Base
		Find(ctx context.Context, filterParam abstraction.Filter, search *abstraction.Search) ([]model.PacketModel, *abstraction.PaginationInfo, error)
		FindByID(ctx context.Context, id string) (*model.PacketModel, error)
		FindByCode(ctx context.Context, code string) (*model.PacketModel, error)
		FindByName(ctx context.Context, name string) (*model.PacketModel, error)
		Create(ctx context.Context, data model.PacketModel) (model.PacketModel, error)
		Creates(ctx context.Context, data []model.PacketModel) ([]model.PacketModel, error)
		UpdateByID(ctx context.Context, id string, data model.PacketModel) (model.PacketModel, error)
		DeleteByID(ctx context.Context, id string) error
		Count(ctx context.Context) (int64, error)
	}

	packet struct {
		Base[model.PacketModel]
	}
)

func NewPacket(conn *gorm.DB) Packet {
	model := model.PacketModel{}
	base := NewBase(conn, model, model.TableName())
	return &packet{
		base,
	}
}
