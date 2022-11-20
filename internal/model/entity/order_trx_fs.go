package entity

import "time"

type OrderTrxFs struct {
	ID           string    `json:"id"`
	TrxType      string    `json:"trx_type"`
	CustomerName string    `json:"customer_name"`
	SupplierName string    `json:"supplier_name"`
	PicName      string    `json:"pic_name"`
	Price        float64   `json:"price"`
	FinalPrice   float64   `json:"final_price"`
	Status       string    `json:"status"`
	StockStatus  string    `json:"stock_status"`
	CreatedAt    string    `json:"created_at"`
	PhoneNumber  string    `json:"phone_number"`
	IsRead       bool      `json:"is_read"`
	Code         string    `json:"code"`
	Created      time.Time `firestore:"created,serverTimestamp"`
	Updated      time.Time `firestore:"updated,omitempty"`
}

type OrderTrxTotalFs struct {
	ID         string `json:"id"`
	TotalOrder int64  `json:"total_order"`
}
