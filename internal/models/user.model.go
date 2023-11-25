package models

type UserModel struct {
	Id            int         `db:"No"`
	Photo_profile interface{} `db:"Photo_profile" form:"Photo_profile" json:"Photo_profile"`
	User_name     *string     `db:"User_name" form:"User_name" json:"User_name"`
	Full_name     *string     `db:"Full_name" form:"Full_name" json:"Full_name"`
	Phone         *string     `db:"Phone" form:"Phone" json:"Phone"`
	Address       *string     `db:"Address" form:"Address" json:"Address"`
	Email         string      `db:"Email" form:"Email" json:"Email"`
	Password      string      `db:"Password" form:"Password" json:"Password"`
	User_type     string      `db:"User_type" form:"User_type" json:"User_type"`
	Otp           *int        `db:"Otp" form:"Otp" json:"Otp"`
	// Created_at    *time.Time  `db:"created_at"`
	// Updated_at    *time.Time  `db:"updated_at"`
}

type QueryParamsUser struct {
	Userid    string `form:"User_id" json:"User_id"`
	Username  string `form:"User_name" json:"User_name"`
	Fullname  string `form:"Full_name" json:"Full_name"`
	Email     string `form:"Email" json:"Email"`
	Phone     string `form:"Phone" json:"Phone"`
	SortOrder string `form:"SortOrder" json:"SortOrder"`
	Page      int    `form:"Page" json:"Page"`
}
