package model

import (
	"errors"
)

type WslogMsg struct {
	BaseModel
	HookId   uint     `gorm:"index" json:"hook_id"`
	UserId   uint     `gorm:"index" json:"user_id"`
	ToUid    uint     `gorm:"index" json:"to_uid"`
	SlackMsg SlackMsg `gorm:"type:json" json:"slack_msg"`
}

func (m *WslogMsg) Create() error {
	m.Id = 0
	return db.Create(m).Error
}

//update renew token and other infos
func (m *WslogMsg) Update() error {
	return db.Model(m).Update(m).Error
}

func (m WslogMsg) All(q *PaginationQ) (list *[]WslogMsg, total uint, err error) {
	list = &[]WslogMsg{}
	tx := db.Model(m).Order("id DESC")
	total, err = crudAll(q, tx, list)
	return
}
func (m WslogMsg) Delete() (err error) {
	if m.Id == 0 {
		return errors.New("resource must not be zero value")
	}
	return crudDelete(m)
}

func (m WslogMsg) Truncate() (err error) {
	return db.Exec("TRUNCATE TABLE wslog_msgs").Error
}
