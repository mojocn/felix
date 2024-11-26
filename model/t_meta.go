package model

type Meta struct {
	ModelBase
	Config Config `json:"config" gorm:"type:json;serializer:json"`
}
