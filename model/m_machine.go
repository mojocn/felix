package model

import (
	"errors"
	"fmt"
	"time"
)

type Machine struct {
	BaseModel
	Name     string `json:"name" gorm:"type:varchar(50);unique_index"`
	Host     string `json:"host" gorm:"type:varchar(50)"`
	Ip       string `json:"ip" gorm:"type:varchar(80)"`
	Port     uint   `json:"port" gorm:"type:int(6)"`
	User     string `json:"user" gorm:"type:varchar(20)"`
	Password string `json:"password,omitempty"`
	Key      string `json:"key,omitempty"`
	Type     string `json:"type" gorm:"type:varchar(20)"`
}

func MachineAdd(name, addr, ip, user, password, key, auth string, port uint) error {
	ins := &Machine{Name: name, Ip: ip, Host: addr, User: user, Password: password, Key: key, Type: auth, Port: port}
	return db.Create(ins).Error
}

func MachineAll(search string) ([]Machine, error) {
	var hs []Machine
	query := db.Order("updated_at")
	if search != "" {
		query = db.Where("name like ?", "%"+search+"%")
	}
	err := query.Find(&hs).Error
	return hs, err
}

func MachineFind(idx uint) (*Machine, error) {
	ins := &Machine{}
	ins.Id = idx
	return ins, db.First(ins).Error
}

func MachineDelete(idx uint) error {
	ins := Machine{}
	ins.Id = idx
	return db.Where("id = ?", idx).Delete(&ins).Error
}
func MachineDeleteAll() error {
	ins := Machine{}
	return db.Delete(&ins).Error
}

func MachineUpdate(name, addr, user, password, pkey, t string, id, port uint) error {
	ins := Machine{Name: name, Host: addr, User: user, Password: password, Key: pkey, Type: t, Port: port}
	wh := &Machine{}
	wh.Id = id
	return db.Model(wh).Updates(ins).Error
}

func MachineDuplicate(idx uint) error {
	ins := &Machine{}
	ins.Id = idx
	err := db.First(ins).Error
	if err != nil {
		return err
	}
	ins.Id = 0
	ins.Name = fmt.Sprintf("%s_du", ins.Name)
	return db.Create(ins).Error
}

func (m *Machine) One() (err error) {
	err = crudOne(m)
	return
}

//All get all for pagination
func (m *Machine) All(q *PaginationQ) (list *[]Machine, total uint, err error) {
	list = &[]Machine{}
	tx := db.Model(m)
	total, err = crudAll(q, tx, list)
	return
}

//Update a row
func (m *Machine) Update() (err error) {
	return db.Model(m).Update(m).Error
}

//Create insert a row
func (m *Machine) Create() (err error) {
	m.Id = 0
	return db.Create(m).Error
}

//Delete destroy a row
func (m *Machine) Delete() (err error) {
	if m.Id == 0 {
		return errors.New("resource must not be zero value")
	}
	return crudDelete(m)
}

func (m *Machine) ChangeUpdateTime() (err error) {
	m.UpdatedAt = time.Now()
	return db.Save(m).Error
}
