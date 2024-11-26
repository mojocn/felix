package model

import "gorm.io/gorm"

type Meta struct {
	gorm.Model
	Config Config `json:"config" gorm:"type:json;serializer:json"`
}
