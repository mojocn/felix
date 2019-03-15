package model

import "time"

type BaseModel struct {
	Id        uint      `gorm:"primary_key" json:"id" form:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
