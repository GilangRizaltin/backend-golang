package models

type Auth struct {
	Id            int         `db:"No" valid:"-"`
	Photo_profile interface{} `db:"Photo_profile" json:"Photo_profile" valid:"-"`
	Full_name     *string     `db:"Full_name" form:"Full_name" json:"Full_name" valid:"required"`
	Email         string      `db:"Email" form:"Email" json:"Email" valid:"email, required"`
	Password      string      `db:"Password" form:"Password" json:"Password" valid:"required"`
	User_type     string      `db:"User_type" form:"User_type" json:"User_type" valid:"-"`
	Otp           int         `db:"Otp" form:"Otp" json:"Otp" valid:"-"`
	ActivateUser  bool        `db:"activated" form:"activated" json:"activated" valid:"-"`
}

type AuthLogin struct {
	Email    string `db:"email" form:"email" json:"email" valid:"email, required"`
	Password string `db:"password" form:"password" json:"password" valid:"required"`
}

type AuthRegister struct {
	Full_name string `db:"full_name" form:"full_name" json:"full_name" valid:"required"`
	Email     string `db:"email" form:"email" json:"email" valid:"email, required"`
	Password  string `db:"password" form:"password" json:"password" valid:"required"`
}
