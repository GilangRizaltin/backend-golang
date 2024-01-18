package models

import "time"

type ProductModel struct {
	Id              int         `db:"No" valid:"-"`
	Product_photo_1 interface{} `db:"Product_photo_1" form:"Product_photo_1" json:"Product_photo_1" valid:"optional"`
	Product_photo_2 interface{} `db:"Product_photo_2" form:"Product_photo_2" json:"Product_photo_2" valid:"optional"`
	Product_photo_3 interface{} `db:"Product_photo_3" form:"Product_photo_3" json:"Product_photo_3" valid:"optional"`
	Product_photo_4 interface{} `db:"Product_photo_4" form:"Product_photo_4" json:"Product_photo_4" valid:"optional"`
	Product_name    string      `db:"Product" form:"Product" json:"Product" valid:"matches(^[a-zA-Z ]+$)"`
	Category        string      `db:"Categories" form:"Categories" json:"Categories" valid:"in(Coffee|Non - Coffee|Food)"`
	Price_default   int         `db:"Price" form:"Price" json:"Price" valid:"numeric, optional"`
	Description     string      `db:"Description" form:"Description" json:"Description" valid:"-"`
	Created_at      *time.Time  `db:"Created_at" valid:"-"`
	// Updated_at      *time.Time  `db:"Updated_at"`
}

type QueryParamsProduct struct {
	ProductId       int    `form:"id" json:"id" valid:"numeric,optional"`
	ProductName     string `form:"search" json:"search" valid:"matches(^[a-zA-Z ]+$), optional"`
	MaximumPrice    int    `form:"maxprice" json:"maxprice" valid:"numeric,optional"`
	MinimumPrice    int    `form:"minprice" json:"minprice" valid:"numeric,optional"`
	ProductCategory string `form:"category" json:"category" valid:"in(Coffee|Non - Coffee|Food), optional"`
	Sort            string `form:"sort" json:"sort" valid:"in(Cheapest|Most Expensive|New Product|Oldest), optional"`
	Page            int    `form:"page" json:"page" valid:"numeric, optional"`
}

type UpdateProduct struct {
	Photo_index   []int  `form:"Photo_index" json:"Photo_index,omitempty" valid:"optional"`
	Product_name  string `form:"Product" json:"Product,omitempty" valid:"matches(^[a-zA-Z ]+$), optional"`
	Category      string `form:"Categories" json:"Categories,omitempty" valid:"in(Coffee|Non - Coffee|Food), optional"`
	Price_default int    `form:"Price" json:"Price,omitempty" valid:"numeric, optional"`
	Description   string `form:"Description" json:"Description,omitempty" valid:"optional"`
}

type PopularProduct struct {
	Product_Id    int         `db:"Id" form:"Id" json:"Id" valid:"optional"`
	Product_name  string      `db:"Product" form:"Product" json:"Product" valid:"matches(^[a-zA-Z ]+$),optional"`
	TotalQuantity interface{} `db:"Total_Quantity" json:"Total_Quantity"`
	Total_Income  interface{} `db:"Total_Income" json:"Total_Income"`
}
