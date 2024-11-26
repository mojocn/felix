package model

import (
	"gorm.io/gorm"
	"time"
)

type ModelBase struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"  gorm:"index"`
}

type CfIp struct {
	ModelBase
	IP    string `json:"ip" gorm:"type:varchar(15)"`
	Cidr  string `json:"cidr" gorm:"type:varchar(18)"`
	Ports []int  `json:"ports" gorm:"type:json;serializer:json"`
}
