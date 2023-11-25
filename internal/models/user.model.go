package models

type UserModel struct {
	Id            int         `db:"No"`
	Photo_profile interface{} `db:"Photo_profile" form:"Photo_profile" json:"Photo_profile" valid:"-"`
	User_name     *string     `db:"User_name" form:"User_name" json:"User_name" valid:"-"`
	Full_name     *string     `db:"Full_name" form:"Full_name" json:"Full_name" valid:"-"`
	Phone         *string     `db:"Phone" form:"Phone" json:"Phone" valid:"-"`
	Address       *string     `db:"Address" form:"Address" json:"Address" valid:"-"`
	Email         string      `db:"Email" form:"Email" json:"Email" valid:"numeric"`
	Password      string      `db:"Password" form:"Password" json:"Password" valid:"-"`
	User_type     string      `db:"User_type" form:"User_type" json:"User_type" valid:"-"`
	Otp           *int        `db:"Otp" form:"Otp" json:"Otp" valid:"-"`
	// Created_at    *time.Time  `db:"created_at"`
}

type QueryParamsUser struct {
	Userid    string `form:"User_id" json:"User_id" valid:"numeric,optional"`
	Username  string `form:"User_name" json:"User_name" valid:"alphanum,optional"`
	Fullname  string `form:"Full_name" json:"Full_name" valid:"alpha, optional"`
	Email     string `form:"Email" json:"Email" valid:"email, optional"`
	Phone     string `form:"Phone" json:"Phone" valid:"numeric, optional"`
	SortOrder string `form:"SortOrder" json:"SortOrder" valid:"in(asc|desc), optional"`
	Page      int    `form:"Page" json:"Page" valid:"numeric"`
}
