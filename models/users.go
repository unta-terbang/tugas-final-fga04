package models

import (
	"time"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
	"main.go/helpers"
)

type User struct {
	Id uint `gorm:"primaryKey" json:"id"`
	Username string `gorm:"uniqueIndex;not null;type:varchar(50)" json:"username" validate:"required-Username is required"`
	Email	string `gorm:"uniqueIndex;not null;type:varchar(150)" json:"email" validate:"required-Email is required,email-invalid email"`
	Password string `gorm:"not null" json:"password" validate:"required-Password is required,MinStringLength(6)-Password has to have a minimum length of 6 characters"`
	Age int `gorm:"not null" json:"age" validate:"required,min=8-Age must be at least 8"`
	ProfileImageUrl string `json:"profile_image_url"`
	CreatedAt time.Time  `json:"created_at"`
	UpdateAt time.Time  `gorm:"autoUpdateTime" json:"update_at"`
	SocialMedias []SocialMedia 		//foreignkeymany
	Photos []Photo 					//foreignkeymany
	Comments []Comment 				//foreignkeymany
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error){

	_, err = govalidator.ValidateStruct(u)
	if err != nil{
		return err
	}

	hashed,err := helpers.HashPwd(u.Password)
	if err != nil{
		return
	}

	u.Password = hashed

	return

}