package models

import (
	"time"
	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
)

type SocialMedia struct {
	Id uint `gorm:"primaryKey" json:"id"`
	Name string `gorm:"not null;type:varchar(50)" json:"name" validate:"required-Name is a must"`
	SocialMediaUrl string `gorm:"not null" json:"social_media_url" validate:"required-SocialMediaUrl is a must"`
	CreatedAt time.Time  `json:"created_at"`
	UpdateAt time.Time  `gorm:"autoUpdateTime" json:"update_at"`
	UserID uint `gorm:"not null" json:"user_id"`
}

func (s *SocialMedia) BeforeCreate(tx *gorm.DB) (err error){
	_,err = govalidator.ValidateStruct(s)

	return
}