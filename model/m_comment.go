package model

import (
	"errors"
	"fmt"
)

type Comment struct {
	BaseModel
	PageUrl     string        `json:"page_url" gorm:"index" form:"page_url"`
	ParentPath  string        `json:"parent_path" form:"parent_path" gorm:"default:'0';index" comment:"父级评论的id路径使用like %查询"`
	UserId      uint          `json:"user_id"`
	Content     string        `json:"content"`
	LikeUids    JsonArrayUint `gorm:"type:json" json:"like_uids"`
	DislikeUids JsonArrayUint `gorm:"type:json" json:"dislike_uids"`
	ThankUids   JsonArrayUint `gorm:"type:json" json:"thank_uids"`
	AtUids      JsonArrayUint `gorm:"type:json" json:"at_uids" comment:"at用户IDs json uint array"`
	User        User          `json:"user"`
}

func (m *Comment) AfterFind() (err error) {
	return
}

//One
func (m *Comment) One() error {
	return crudOne(m)
}

//All
func (m *Comment) All(q *PaginationQ) (list *[]Comment, total uint, err error) {
	list = &[]Comment{}
	tx := db.Model(m).Order("CreatedAt").Where("page_url = ?", m.PageUrl).Preload("User")
	total, err = crudAll(q, tx, list)
	return
}

func (m *Comment) Action(id, uid uint, action string) (err error) {
	err = db.First(m, id).Error
	if err != nil {
		return
	}
	//if m.UserId == uid {
	//	return fmt.Errorf("you cann't take action on your own thread(%d)", id)
	//}
	switch action {
	case "like":
		_, t := m.LikeUids.AppendOrRemove(uid)
		m.LikeUids = t
	case "dislike":
		_, t := m.DislikeUids.AppendOrRemove(uid)
		m.DislikeUids = t
	case "thank":
		_, t := m.ThankUids.AppendOrRemove(uid)
		m.ThankUids = t
	default:
		return fmt.Errorf("this is has no action for %s", action)
	}
	return db.Model(m).Update(m).Error
}

//Create
func (m *Comment) Create() (err error) {
	m.Id = 0
	m.LikeUids = nil
	m.ThankUids = nil
	m.DislikeUids = nil
	return db.Create(m).Error
}

//Delete
func (m *Comment) Delete() (err error) {
	if m.Id == 0 {
		return errors.New("resource must not be zero value")
	}
	return crudDelete(m)
}
