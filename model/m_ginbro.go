package model

type Ginbro struct {
	BaseModel
	IsSuccess  bool   `json:"is_success" form:"is_success"`
	AppSecret  string `json:"app_secret" form:"app_secret"`
	AppAddr    string `json:"app_addr" form:"app_addr"`
	AppDir     string `json:"app_dir" form:"app_dir"`
	AppPkg     string `json:"app_pkg" form:"app_pkg"`
	AuthTable  string `json:"auth_table" form:"auth_table"`
	AuthColumn string `json:"auth_column" form:"auth_column"`
	DbUser     string `json:"db_user" form:"db_user"`
	DbPassword string `json:"db_password" form:"db_password"`
	DbAddr     string `json:"db_addr" form:"db_addr"`
	DbName     string `json:"db_name" form:"db_name"`
	DbChar     string `json:"db_char" form:"db_char"`
	DbType     string `json:"db_type" form:"db_type"`
}

//CreateUserOfRole
func (m *Ginbro) Create() (err error) {
	m.Id = 0
	return db.Create(m).Error
}

//All
func (m Ginbro) All(q *PaginationQ) (list *[]Ginbro, total uint, err error) {
	list = &[]Ginbro{}
	total, err = crudAll(q, db.Model(m), list)
	return
}
