package models

type Auth struct {
	Id            int         `db:"No" valid:"-"`
	Photo_profile interface{} `db:"Photo_profile" json:"Photo_profile" valid:"-"`
	Full_name     *string     `db:"Full_name" form:"Full_name" json:"Full_name" valid:"required"`
	Email         string      `db:"Email" form:"Email" json:"Email" valid:"email, required"`
	Password      string      `db:"Password" form:"Password" json:"Password" valid:"-"`
	User_type     string      `db:"User_type" form:"User_type" json:"User_type" valid:"-"`
	Otp           *int        `db:"Otp" form:"Otp" json:"Otp" valid:"-"`
}

type AuthLogin struct {
	Id        int     `db:"No" valid:"-"`
	Full_name *string `db:"Full_name" form:"Full_name" json:"Full_name" valid:"-"`
	Email     string  `db:"Email" form:"Email" json:"Email" valid:"email, required"`
	Password  string  `db:"Password" form:"Password" json:"Password" valid:"required"`
	User_type string  `db:"User_type" form:"User_type" json:"User_type" valid:"-"`
	Otp       *int    `db:"Otp" form:"Otp" json:"Otp" valid:"-"`
}
