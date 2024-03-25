package models

import (
	"time"
	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Comment struct {
	Id uint `gorm:"primaryKey" json:"id"`
	UserID uint `gorm:"not null" json:"user_id"`
	PhotoID uint `gorm:"not null" json:"photo_id"`
	Message string `gorm:"not null;type:varchar(200)" json:"message" validate:"required-Message is a must"`
	CreatedAt time.Time  `json:"created_at"`
	UpdateAt time.Time  `gorm:"autoUpdateTime" json:"update_at"`
}

func (c *Comment) BeforeCreate(tx *gorm.DB) (err error){
	_,err = govalidator.ValidateStruct(c)

	return
}