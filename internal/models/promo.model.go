package models

import "time"

type PromoModel struct {
	Id             string     `db:"No"`
	Promo_code     string     `db:"Promo_code" form:"Promo_code" json:"Promo_code"`
	Promo_type     string     `db:"Promo_type" form:"Promo_type" json:"Promo_type"`
	Flat_amount    int        `db:"Flat_amount" form:"Flat_amount" json:"Flat_amount"`
	Percent_amount float64    `db:"Percent_amount" form:"Percent_amount" json:"Percent_amount"`
	Created_at     *time.Time `db:"Time_created"`
	Ended_at       *time.Time `db:"Time_ended" form:"Time_ended" json:"Time_ended"`
}
