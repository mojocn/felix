package model

type Config struct {
	BaseModel
	User     string    `json:"user"`
	Password string    `json:"password"`
	Pkey     string    `json:"pkey"`
	Hosts    []Machine `json:"hosts"`
}
