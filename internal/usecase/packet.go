package usecase

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"gitlab.com/stoqu/stoqu-be/internal/config"

	"gitlab.com/stoqu/stoqu-be/internal/factory/repository"
	"gitlab.com/stoqu/stoqu-be/internal/model/abstraction"
	"gitlab.com/stoqu/stoqu-be/internal/model/dto"
	model "gitlab.com/stoqu/stoqu-be/internal/model/entity"
	"gitlab.com/stoqu/stoqu-be/pkg/constant"
	res "gitlab.com/stoqu/stoqu-be/pkg/util/response"
	"gitlab.com/stoqu/stoqu-be/pkg/util/str"
	"gitlab.com/stoqu/stoqu-be/pkg/util/trxmanager"
)

type Packet interface {
	Find(ctx context.Context, filterParam abstraction.Filter) ([]dto.PacketResponse, abstraction.PaginationInfo, error)
	FindByID(ctx context.Context, payload dto.ByIDRequest) (dto.PacketResponse, error)
	Create(ctx context.Context, payload dto.CreatePacketRequest) (dto.PacketResponse, error)
	Update(ctx context.Context, payload dto.UpdatePacketRequest) (dto.PacketResponse, error)
	Delete(ctx context.Context, payload dto.ByIDRequest) (dto.PacketResponse, error)
}

type packet struct {
	Repo repository.Factory
	Cfg  *config.Configuration
}

func NewPacket(cfg *config.Configuration, f repository.Factory) Packet {
	return &packet{f, cfg}
}

func (u *packet) Find(ctx context.Context, filterParam abstraction.Filter) (result []dto.PacketResponse, pagination abstraction.PaginationInfo, err error) {
	var search *abstraction.Search
	if filterParam.Search != "" {
		searchQuery := "lower(code) LIKE ? OR lower(name) LIKE ?"
		searchVal := "%" + strings.ToLower(filterParam.Search) + "%"
		search = &abstraction.Search{
			Query: searchQuery,
			Args:  []interface{}{searchVal, searchVal},
		}
	}

	packets, info, err := u.Repo.Packet.Find(ctx, filterParam, search)
	if err != nil {
		return nil, pagination, res.ErrorBuilder(res.Constant.Error.InternalServerError, err)
	}
	pagination = *info

	for _, packet := range packets {
		result = append(result, dto.PacketResponse{
			PacketModel: packet,
		})
	}

	return result, pagination, nil
}

func (u *packet) FindByID(ctx context.Context, payload dto.ByIDRequest) (dto.PacketResponse, error) {
	var result dto.PacketResponse

	packet, err := u.Repo.Packet.FindByID(ctx, payload.ID)
	if err != nil {
		return result, err
	}

	result = dto.PacketResponse{
		PacketModel: *packet,
	}

	return result, nil
}

func (u *packet) Create(ctx context.Context, payload dto.CreatePacketRequest) (result dto.PacketResponse, err error) {
	var (
		packetID = uuid.New().String()
		packet   = model.PacketModel{
			Entity: model.Entity{
				ID: packetID,
			},
			PacketEntity: model.PacketEntity{
				Code:   str.GenCode(constant.CODE_PACKET_PREFIX),
				Value:  payload.Value,
				UnitID: payload.UnitID,
			},
		}
	)

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		_, err = u.Repo.Packet.Create(ctx, packet)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.PacketResponse{
		PacketModel: packet,
	}

	return result, nil
}

func (u *packet) Update(ctx context.Context, payload dto.UpdatePacketRequest) (result dto.PacketResponse, err error) {
	var (
		packetData = &model.PacketModel{
			PacketEntity: model.PacketEntity{
				UnitID: payload.UnitID,
				Value:  payload.Value,
			},
		}
	)

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		_, err = u.Repo.Packet.UpdateByID(ctx, payload.ID, *packetData)
		if err != nil {
			return err
		}

		packetData, err = u.Repo.Packet.FindByID(ctx, payload.ID)
		if err != nil {
			return res.ErrorBuilder(res.Constant.Error.NotFound, err)
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.PacketResponse{
		PacketModel: *packetData,
	}

	return result, nil
}

func (u *packet) Delete(ctx context.Context, payload dto.ByIDRequest) (result dto.PacketResponse, err error) {
	var data *model.PacketModel

	if err = trxmanager.New(u.Repo.Db).WithTrx(ctx, func(ctx context.Context) error {
		data, err = u.Repo.Packet.FindByID(ctx, payload.ID)
		if err != nil {
			return err
		}

		err = u.Repo.Packet.DeleteByID(ctx, payload.ID)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return result, err
	}

	result = dto.PacketResponse{
		PacketModel: *data,
	}

	return result, nil
}
