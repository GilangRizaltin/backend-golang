package models

import "time"

type PromoModel struct {
	Id             int        `db:"No" valid:"-"`
	Promo_code     string     `db:"Promo_code" form:"Promo_code" json:"Promo_code" valid:"uppercase, alphanum"`
	Promo_type     string     `db:"Promo_type" form:"Promo_type" json:"Promo_type" valid:"in(Flat|Percent)"`
	Flat_amount    int        `db:"Flat_amount" form:"Flat_amount" json:"Flat_amount" valid:"numeric, optional"`
	Percent_amount float64    `db:"Percent_amount" form:"Percent_amount" json:"Percent_amount" valid:"float, optional"`
	Created_at     *time.Time `db:"Time_created" valid:"optional"`
	Ended_at       string     `db:"Ended_at" form:"Ended_at" json:"Ended_at" valid:"optional"`
	Duration       int        `db:"Duration" form:"Duration" json:"Duration" valid:"optional"`
}

type QueryParamsPromo struct {
	Promo_code string     `form:"promo-code" json:"promo-code" valid:"alphanum, optional"`
	Time_end   *time.Time `form:"time-end" json:"time-end" valid:"alpha, optional"`
	Page       int        `form:"page" json:"page" valid:"numeric, optional"`
}

type UpdatePromoModel struct {
	Id             int        `db:"No" valid:"-"`
	Promo_code     string     `db:"Promo_code" form:"Promo_code" json:"Promo_code" valid:"uppercase, alphanum, optional"`
	Promo_type     string     `db:"Promo_type" form:"Promo_type" json:"Promo_type" valid:"in(Flat|Percent), optional"`
	Flat_amount    int        `db:"Flat_amount" form:"Flat_amount" json:"Flat_amount" valid:"numeric, optional"`
	Percent_amount float64    `db:"Percent_amount" form:"Percent_amount" json:"Percent_amount" valid:"float, optional"`
	Created_at     *time.Time `db:"Time_created" valid:"optional"`
	Ended_at       string     `db:"Ended_at" form:"Ended_at" json:"Ended_at" valid:"optional"`
	Duration       int        `db:"Duration" form:"Duration" json:"Duration" valid:"optional"`
}
