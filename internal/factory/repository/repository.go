package repository

import (
	"cloud.google.com/go/firestore"
	"firebase.google.com/go/messaging"
	el "github.com/olivere/elastic/v7"
	"gitlab.com/stoqu/stoqu-be/internal/config"
	dbRepository "gitlab.com/stoqu/stoqu-be/internal/repository/db"
	fbRepository "gitlab.com/stoqu/stoqu-be/internal/repository/firebase"
	"gorm.io/gorm"
)

type Factory struct {
	Db        *gorm.DB
	Es        *el.Client
	Fcm       *messaging.Client
	Firestore *firestore.Client

	Role        dbRepository.Role
	User        dbRepository.User
	UserProfile dbRepository.UserProfile

	Unit                 dbRepository.Unit
	Packet               dbRepository.Packet
	ConvertionUnit       dbRepository.ConvertionUnit
	Currency             dbRepository.Currency
	ReminderStock        dbRepository.ReminderStock
	ReminderStockHistory dbRepository.ReminderStockHistory
	Rack                 dbRepository.Rack

	Brand   dbRepository.Brand
	Variant dbRepository.Variant
	Product dbRepository.Product

	Stock              dbRepository.Stock
	StockRack          dbRepository.StockRack
	StockLookup        dbRepository.StockLookup
	StockTrx           dbRepository.StockTrx
	StockTrxItem       dbRepository.StockTrxItem
	StockTrxItemLookup dbRepository.StockTrxItemLookup

	OrderTrx           dbRepository.OrderTrx
	OrderTrxItem       dbRepository.OrderTrxItem
	OrderTrxItemLookup dbRepository.OrderTrxItemLookup
	OrderTrxStatus     dbRepository.OrderTrxStatus
	OrderTrxReceipt    dbRepository.OrderTrxReceipt

	OrderFs fbRepository.OrderFs
}

func Init(cfg *config.Configuration, db *gorm.DB, fsClient *firestore.Client) Factory {
	f := Factory{}

	// db
	f.Db = db
	f.Role = dbRepository.NewRole(f.Db)
	f.User = dbRepository.NewUser(f.Db)
	f.UserProfile = dbRepository.NewUserProfile(f.Db)

	f.Unit = dbRepository.NewUnit(f.Db)
	f.Packet = dbRepository.NewPacket(f.Db)
	f.ConvertionUnit = dbRepository.NewConvertionUnit(f.Db)
	f.Currency = dbRepository.NewCurrency(f.Db)
	f.ReminderStock = dbRepository.NewReminderStock(f.Db)
	f.ReminderStockHistory = dbRepository.NewReminderStockHistory(f.Db)
	f.Rack = dbRepository.NewRack(f.Db)

	f.Brand = dbRepository.NewBrand(f.Db)
	f.Variant = dbRepository.NewVariant(f.Db)
	f.Product = dbRepository.NewProduct(f.Db)

	f.Stock = dbRepository.NewStock(f.Db)
	f.StockRack = dbRepository.NewStockRack(f.Db)
	f.StockLookup = dbRepository.NewStockLookup(f.Db)
	f.StockTrx = dbRepository.NewStockTrx(f.Db)
	f.StockTrxItem = dbRepository.NewStockTrxItem(f.Db)
	f.StockTrxItemLookup = dbRepository.NewStockTrxItemLookup(f.Db)

	f.OrderTrx = dbRepository.NewOrderTrx(f.Db)
	f.OrderTrxItem = dbRepository.NewOrderTrxItem(f.Db)
	f.OrderTrxItemLookup = dbRepository.NewOrderTrxItemLookup(f.Db)
	f.OrderTrxStatus = dbRepository.NewOrderTrxStatus(f.Db)
	f.OrderTrxReceipt = dbRepository.NewOrderTrxReceipt(f.Db)

	// firestore
	f.Firestore = fsClient
	f.OrderFs = fbRepository.NewOrderFs(f.Firestore)

	return f
}
