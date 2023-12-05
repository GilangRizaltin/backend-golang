package models

import (
	"time"
)

type OrderModel struct {
	Id                int                `db:"No" valid:"-"`
	User              string             `db:"User" form:"User" json:"User" valid:"-"`
	Subtotal          *int               `db:"Subtotal" form:"Subtotal" json:"Subtotal" valid:"numeric"`
	Promo             string             `db:"Promo" form:"Promo" json:"Promo" valid:"-"`
	Percent_discount  float64            `db:"Discount_Percentage" form:"Percent_discount" json:"Percent_discount" valid:"-"`
	Flat_discount     int                `db:"Discount_Flat" form:"Flat_discount" json:"Flat_discount" valid:"-"`
	Serve             string             `db:"Serve" form:"Serve" json:"Serve" valid:"-"`
	Fee               int                `db:"Serving_Fee" form:"Fee" json:"Fee" valid:"-"`
	Tax               float64            `db:"Tax" form:"Tax" json:"Tax" valid:"-"`
	Total_transaction *int               `db:"Total_transaction" form:"Total_transaction" json:"Total_transaction" valid:"-"`
	Payment_type      string             `db:"Payment_type" form:"Payment_type" json:"Payment_type" valid:"-"`
	Status            string             `db:"Status" form:"Status" json:"Status" valid:"-"`
	Created_at        *time.Time         `db:"Date" valid:"-"`
	Product           []OrderDetailModel `form:"products" json:"products" valid:"-"`
}

type OrderDetailModel struct {
	Order_id         int    `db:"No Order" form:"order_id" json:"order_id" valid:"-"`
	Product_name     string `db:"Product_name" form:"Product_name" json:"Product_name" valid:"-"`
	Size             string `db:"Size" form:"Size" json:"Size" valid:"in(small|medium|large|Short|Regular|Grande|Venti)"`
	Hot_or_not       bool   `db:"Hot_or_not" form:"Hot_or_not" json:"Hot_or_not" valid:"-"`
	Price            int    `db:"Price" form:"Price" json:"Price" valid:"-"`
	Quantity         int    `db:"Quantity" form:"Quantity" json:"Quantity" valid:"-"`
	Subtotal_product int    `db:"Subtotal_product" form:"Subtotal_product" json:"Subtotal_product" valid:"-"`
}

type QueryParamsOrder struct {
	Status   string `form:"status" json:"status" valid:"in(On progress|Pending|Done|Cancelled),optional"`
	Sort     string `form:"sort" json:"sort" valid:"in(Newest|Oldest),optional"`
	Order_id int    `form:"order_id" json:"order_id" valid:"numeric,optional"`
	Page     int    `form:"page" json:"page" valid:"numeric,optional"`
}

type StatisticOrder struct {
	OrderDate     *time.Time  `db:"OrderDate" json:"OrderDate"`
	TotalQuantity interface{} `db:"TotalQuantity" json:"TotalQuantity"`
}

type OrderDataStatus struct {
	Status string `db:"Status" form:"Status" json:"Status" valid:"-"`
	Total  int    `db:"Total" form:"Total" json:"Total" valid:"-"`
}
