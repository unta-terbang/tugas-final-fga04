package models

import (
	"time"

	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type Photo struct {
	Id uint `gorm:"primaryKey" json:"id"`
	Title string `gorm:"not null;type:varchar(100)" json:"title" validate:"required-Tittle is a must"`
	Caption string `gorm:"varchar(200)" json:"caption" validate:"required-Caption is a must"`
	PhotoUrl string `gorm:"not null" json:"photo_url"`
	UserID uint `gorm:"not null" json:"user_id"` //foreignkeyone
	CreatedAt time.Time  `json:"created_at"`
	UpdateAt time.Time  `gorm:"autoUpdateTime" json:"update_at"`
}

func (p *Photo) BeforeCreate(tx *gorm.DB) (err error){
	_,err = govalidator.ValidateStruct(p)

	return
}