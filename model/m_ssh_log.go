package model

import (
	"errors"
	"time"
)

type SshLog struct {
	BaseModel
	UserId    uint      `gorm:"index" json:"user_id" form:"user_id"`
	MachineId uint      `gorm:"index" json:"machine_id" form:"machine_id"`
	SshUser   string    `json:"ssh_user" comment:"ssh账号"`
	ClientIp  string    `json:"client_ip" form:"client_ip"`
	StartedAt time.Time `json:"started_at" form:"started_at"`
	Status    uint      `json:"status" comment:"0-未标记 2-正常 4-警告 8-危险 16-致命"`
	Remark    string    `json:"remark"`
	Log       string    `gorm:"type:text" json:"log"`
	Machine   Machine   `gorm:"association_autoupdate:false;association_autocreate:false" json:"machine"`
	User      User      `gorm:"association_autoupdate:false;association_autocreate:false" json:"user"`
}

type SshLogQ struct {
	SshLog
	PaginationQ
	FromTime string `form:"from_time"`
	ToTime   string `form:"to_time"`
}

func (m SshLogQ) Search() (list *[]SshLog, total uint, err error) {
	list = &[]SshLog{}
	tx := db.Model(m.SshLog).Preload("User").Preload("Machine")
	if m.ClientIp != "" {
		tx = tx.Where("client_ip like ?", "%"+m.ClientIp+"%")
	}
	if m.FromTime != "" && m.ToTime != "" {
		tx = tx.Where("`created_at` BETWEEN ? AND ?", m.FromTime, m.ToTime)
	}
	total, err = crudAll(&m.PaginationQ, tx, list)
	return
}

func (m *SshLog) AfterFind() (err error) {
	return
}

//One
func (m *SshLog) One() error {
	return crudOne(m)
}

//All
func (m SshLog) All(q *PaginationQ) (list *[]SshLog, total uint, err error) {
	list = &[]SshLog{}
	tx := db.Model(m)
	total, err = crudAll(q, tx, list)
	return
}

//Update
func (m *SshLog) Update() (err error) {
	if m.Id < 1 {
		return errors.New("id必须大于0")
	}
	return db.Model(m).Update(m).Error
}

//Create
func (m *SshLog) Create() (err error) {
	m.Id = 0
	return db.Create(m).Error
}

//Delete
func (m *SshLog) Delete() (err error) {
	if m.Id < 1 {
		return errors.New("id必须大于0")
	}
	return crudDelete(m)
}
