package model

import "gorm.io/gorm"

type Proxy struct {
	gorm.Model
	URI string `json:"uri" gorm:"text"`
}
