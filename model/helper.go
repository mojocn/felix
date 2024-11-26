package model

import (
	"errors"
	"gorm.io/gorm"
)

// PaginationQ gin handler query binding struct
type PaginationQ struct {
	Ok    bool        `json:"ok"`
	Size  int         `form:"size" json:"size"`
	Page  int         `form:"page" json:"page"`
	Data  interface{} `json:"data" comment:"muster be a pointer of slice gorm.Model"` // save pagination list
	Total int64       `json:"total"`
}

// SearchAll optimized pagination method for gorm
func (p *PaginationQ) SearchAll(queryTx *gorm.DB) (data *PaginationQ, err error) {
	//99999 magic number for get all list without pagination
	if p.Size == 9999 || p.Size == 99999 {
		err = queryTx.Find(p.Data).Error
		p.Ok = err == nil
		return p, err
	}

	if p.Size < 1 {
		p.Size = 10
	}
	if p.Page < 1 {
		p.Page = 1
	}
	offset := p.Size * (p.Page - 1)
	err = queryTx.Count(&p.Total).Error
	if err != nil {
		return p, err
	}
	err = queryTx.Limit(p.Size).Offset(offset).Find(p.Data).Error
	p.Ok = err == nil
	return p, err
}

func crudAll(p *PaginationQ, queryTx *gorm.DB, list interface{}) (int64, error) {
	if p.Size < 1 {
		p.Size = 10
	}
	if p.Page < 1 {
		p.Page = 1
	}

	var total int64
	err := queryTx.Count(&total).Error
	if err != nil {
		return 0, err
	}
	offset := p.Size * (p.Page - 1)
	err = queryTx.Limit(p.Size).Offset(offset).Find(list).Error
	if err != nil {
		return 0, err
	}
	return total, err
}

func crudOne(m interface{}) (err error) {
	if err := db.First(m).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("resource is not found")
	}
	return nil
}

func crudDelete(m interface{}) (err error) {
	//WARNING When delete a record, you need to ensure it’s primary field has value, and GORM will use the primary key to delete the record, if primary field’s blank, GORM will delete all records for the model
	//primary key must be not zero value
	db := db.Unscoped().Delete(m)
	if err = db.Error; err != nil {
		return
	}
	if db.RowsAffected != 1 {
		return errors.New("resource is not found to destroy")
	}
	return nil
}
