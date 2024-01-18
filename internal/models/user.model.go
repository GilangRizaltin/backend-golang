package models

import (
	"mime/multipart"
	"time"
)

type UserModel struct {
	Id            int         `db:"No" json:"Id" valid:"optional"`
	Photo_profile interface{} `db:"Photo_profile" form:"_" json:"Photo_profile" valid:"-"`
	User_name     *string     `db:"User_name" form:"User_name" json:"User_name" valid:"alphanum,optional"`
	Full_name     string      `db:"Full_name" form:"Full_name" json:"Full_name" valid:"matches(^[a-zA-Z ]+$)"`
	Phone         *string     `db:"Phone" form:"Phone" json:"Phone" valid:"optional"`
	Address       *string     `db:"Address" form:"Address" json:"Address" valid:"optional"`
	Email         string      `db:"Email" form:"Email" json:"Email" valid:"email"`
	Password      string      `db:"Password" form:"Password" json:"Password" valid:"required"`
	User_type     string      `db:"User_type" form:"User_type" json:"User_type" valid:"in(Admin|Normal User)optional"`
	Otp           *int        `db:"Otp" form:"Otp" json:"Otp" valid:"-"`
	Created_at    *time.Time  `db:"created_at" json:"created_at" valid:"-"`
}

type QueryParamsUser struct {
	Userid    string `form:"user_id" json:"user_id" valid:"numeric,optional"`
	Username  string `form:"user_name" json:"user_name" valid:"alphanum,optional"`
	Fullname  string `form:"full_name" json:"full_name" valid:"matches(^[a-zA-Z ]+$), optional"`
	Email     string `form:"email" json:"email" valid:"email, optional"`
	Phone     string `form:"phone" json:"phone" valid:"numeric, optional"`
	SortOrder string `form:"sortOrder" json:"sortOrder" valid:"in(asc|desc), optional"`
	Page      int    `form:"page" json:"page" valid:"numeric, optional"`
}

type UserUpdateModel struct {
	Userid        int            `form:"user_id" json:"user_id" valid:"numeric,optional"`
	Photo_profile multipart.File `form:"_" json:"Photo_profile" valid:"optional"`
	User_name     string         `form:"User_name" json:"User_name" valid:"alphanum,optional"`
	Full_name     string         `form:"Full_name" json:"Full_name" valid:"matches(^[a-zA-Z ]+$), optional"`
	Phone         *string        `form:"Phone" json:"Phone" valid:"optional"`
	Address       *string        `form:"Address" json:"Address" valid:"optional"`
	LastPassword  string         `form:"Last_Password" json:"Last_Password" valid:"optional"`
	NewPassword   string         `form:"New_Password" json:"New_Password" valid:"optional"`
	User_type     string         `form:"User_type" json:"User_type" valid:"in(Admin|Normal User),optional"`
}
