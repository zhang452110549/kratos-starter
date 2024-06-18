package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName string `gorm:"column:f_user_name;not null;type:varchar(255);comment:用户名" json:"user_name"`
	Password string `gorm:"column:f_password;not null;type:varchar(255);comment:密码" json:"password"`
}
