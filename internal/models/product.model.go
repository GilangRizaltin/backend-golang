package models

import "time"

type ProductModel struct {
	Id              int         `db:"No"`
	Product_photo_1 interface{} `db:"Product_photo_1" form:"Product_photo_1" json:"Product_photo_1"`
	Product_photo_2 interface{} `db:"Product_photo_2" form:"Product_photo_2" json:"Product_photo_2"`
	Product_photo_3 interface{} `db:"Product_photo_3" form:"Product_photo_3" json:"Product_photo_3"`
	Product_photo_4 interface{} `db:"Product_photo_4" form:"Product_photo_4" json:"Product_photo_4"`
	Product_name    string      `db:"Product" form:"Product" json:"Product"`
	Category        string      `db:"Categories" form:"Categories" json:"Categories"`
	Price_default   int         `db:"Price" form:"Price" json:"Price"`
	Description     string      `db:"Description" form:"Description" json:"Description"`
	Created_at      *time.Time  `db:"Created_at"`
	// Updated_at      *time.Time  `db:"Updated_at"`
}
