package model

import (
	"gorm.io/gorm"
)

type CfIp struct {
	gorm.Model
	IP    string `json:"ip" gorm:"type:varchar(15)"`
	CIDR  string `json:"cidr" gorm:"type:varchar(18)"`
	Ports []int  `json:"ports" gorm:"type:json;serializer:json"`
}
