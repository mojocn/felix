package model

import (
	"fmt"
	"time"
)

type Task struct {
	BaseModel
	Content  string    `json:"content" gorm:"type:varchar(255)"`
	Deadline time.Time `json:"deadline"`
	Category string    `json:"category" gorm:"type:varchar(20)"`
	Status   string    `json:"type" gorm:"type:varchar(20)"`
}

const TimeLayout = "2006-01-02T15:04"

func TaskAdd(content, cate, deadLine string) error {

	ins := &Task{Content: content, Category: cate, Status: "TODO"}
	if deadLine == "" {
		ins.Deadline = time.Now().AddDate(0, 0, 7)
	} else {
		t, err := time.Parse(TimeLayout, deadLine)
		if err != nil {
			return fmt.Errorf("unknown time [%s] to parse by time laytou [%s]", deadLine, TimeLayout)
		}
		ins.Deadline = t
	}
	return db.Create(ins).Error
}

func TaskAll(search string) ([]Task, error) {
	var ts []Task
	query := db.Order("created_at desc")
	if search != "" {
		query = db.Where("content like ?", "%"+search+"%")
	}
	err := query.Find(&ts).Error
	return ts, err
}

func TaskRm(id uint) error {
	ins := Task{}
	ins.Id = id
	return db.Where("id = ?", id).Delete(&ins).Error
}

func TaskUpdate(id uint, status string) error {
	ins := Task{Status: status}
	wh := &Task{}
	wh.Id = id
	return db.Model(wh).Updates(ins).Error
}
