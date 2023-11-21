package models

import "time"

type OrderModel struct {
	Id                int        `db:"No"`
	User              string     `db:"User" form:"User" json:"User"`
	Subtotal          int        `db:"Subtotal" form:"Subtotal" json:"Subtotal"`
	Promo             string     `db:"Promo" form:"Promo" json:"Promo"`
	Percent_discount  float64    `db:"Percent_discount" form:"Percent_discount" json:"Percent_discount"`
	Flat_discount     int        `db:"Flat_discount" form:"Flat_discount" json:"Flat_discount"`
	Serve             string     `db:"Serve" form:"Serve" json:"Serve"`
	Fee               int        `db:"Fee" form:"Fee" json:"Fee"`
	Tax               float64    `db:"Tax" form:"Tax" json:"Tax"`
	Total_transaction int        `db:"Total_transaction" form:"Total_transaction" json:"Total_transaction"`
	Payment_type      string     `db:"Payment_type" form:"Payment_type" json:"Payment_type"`
	Status            string     `db:"Status" form:"Status" json:"Status"`
	Created_at        *time.Time `db:"Created_at"`
	// Product           []OrderDetailModel
}

type OrderDetailModel struct {
	Order_id     int    `db:"No Order"`
	Product_name string `db:"Product_name" form:"Product_name" json:"Product_name"`
	Size         string `db:"Size" form:"Size" json:"Size"`
	Hot_or_not   bool   `db:"Hot_or_not" form:"Hot_or_not" json:"Hot_or_not"`
	Price        int    `db:"Price" form:"Price" json:"Price"`
	Quantity     int    `db:"Quantity" form:"Quantity" json:"Quantity"`
	// Subtotal_product int    `db:"Subtotal_product" form:"Subtotal_product" json:"Subtotal_product"`
}

// func (m *OrderProduct) SetProduct() {

// }
