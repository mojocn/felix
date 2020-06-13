package model

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/libragen/felix/util"
)

type WslogHook struct {
	BaseModel
	Name     string    `json:"name" gorm:"type:varchar(50);unique_index"`
	UserId   uint      `json:"user_id"`
	Token    string    `json:"token"`
	ExpireAt time.Time `json:"expire_at"`
}

func (m *WslogHook) Create() error {
	m.Id = 0
	err := db.Create(m).Error
	if err != nil {
		return err
	}
	err = m.generateAesToken()
	if err != nil {
		return err
	}
	return db.Save(m).Error
}

//update renew token and other infos
func (m *WslogHook) Update() error {
	err := m.generateAesToken()
	if err != nil {
		return err
	}
	return db.Model(m).Update(m).Error
}

func (m WslogHook) All(q *PaginationQ) (list *[]WslogHook, total uint, err error) {
	list = &[]WslogHook{}
	total, err = crudAll(q, db.Model(m), list)
	return
}
func (m WslogHook) Delete() (err error) {
	if m.Id == 0 {
		return errors.New("resource must not be zero value")
	}
	return crudDelete(m)
}
func (m *WslogHook) One() error {
	return crudOne(m)
}

func (m *WslogHook) generateAesToken() error {
	tokenString, err := util.AesEncrypt([]byte(fmt.Sprintf("%d", m.Id)), AppSecret)
	if err != nil {
		return err
	}
	m.Token = tokenString
	return nil
}

const godToken = "felix_websocket_log_rock"

func WslogHookCheckToken(token string) (*WslogHook, error) {
	if token == godToken {
		gm := WslogHook{}
		gm.Id = 1
		gm.Name = "GodHook"
		gm.ExpireAt = time.Now().Add(time.Hour * 24 * 365 * 10)
		gm.UserId = 1
		gm.Token = token
		gm.CreatedAt = time.Now()
		gm.UpdatedAt = time.Now()
		return &gm, nil
	}
	msg, err := util.AesDecrypt(token, AppSecret)
	if err != nil {
		return nil, err
	}
	hookID, err := strconv.ParseUint(string(msg), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%s's is not a int of decryption", token)
	}
	m := WslogHook{}
	m.Id = uint(hookID)
	err = m.One()
	if err != nil {
		return nil, fmt.Errorf("token: %s  your hook is not vilad or has been baned", token)
	}
	//check expire_at in database table
	if m.ExpireAt.Before(time.Now()) {
		return nil, fmt.Errorf("token has been expired at %v, please update expire_at column in database table", m.ExpireAt)
	}
	return &m, nil
}
