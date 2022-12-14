package usecase

import (
	"gitlab.com/stoqu/stoqu-be/internal/config"
	"gitlab.com/stoqu/stoqu-be/internal/factory/repository"
	"gitlab.com/stoqu/stoqu-be/internal/usecase"
)

type Factory struct {
	Role usecase.Role
	User usecase.User
	Auth usecase.Auth

	Unit                 usecase.Unit
	Packet               usecase.Packet
	ConvertionUnit       usecase.ConvertionUnit
	Currency             usecase.Currency
	ReminderStock        usecase.ReminderStock
	ReminderStockHistory usecase.ReminderStockHistory
	Rack                 usecase.Rack

	Brand   usecase.Brand
	Variant usecase.Variant
	Product usecase.Product

	Dashboard usecase.Dashboard

	Stock       usecase.Stock
	StockLookup usecase.StockLookup

	Report usecase.Report

	Order usecase.Order
}

func Init(cfg *config.Configuration, r repository.Factory) Factory {
	f := Factory{}

	f.Role = usecase.NewRole(cfg, r)
	f.User = usecase.NewUser(cfg, r)
	f.Auth = usecase.NewAuth(cfg, r)

	f.Unit = usecase.NewUnit(cfg, r)
	f.Packet = usecase.NewPacket(cfg, r)
	f.ConvertionUnit = usecase.NewConvertionUnit(cfg, r)
	f.Currency = usecase.NewCurrency(cfg, r)
	f.ReminderStock = usecase.NewReminderStock(cfg, r)
	f.ReminderStockHistory = usecase.NewReminderStockHistory(cfg, r)
	f.Rack = usecase.NewRack(cfg, r)

	f.Brand = usecase.NewBrand(cfg, r)
	f.Variant = usecase.NewVariant(cfg, r)
	f.Product = usecase.NewProduct(cfg, r)

	f.Dashboard = usecase.NewDashboard(cfg, r)

	f.Stock = usecase.NewStock(cfg, r)
	f.StockLookup = usecase.NewStockLookup(cfg, r)

	f.Report = usecase.NewReport(cfg, r)

	f.Order = usecase.NewOrder(cfg, r, f.Stock)

	return f
}
