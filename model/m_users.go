package model

import (
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var _ = time.Thursday

//User
type User struct {
	BaseModel
	Username       string        `gorm:"column:username" form:"username" json:"username" comment:"昵称/登陆用户名" columnType:"varchar(50)" dataType:"varchar" columnKey:"UNI"`
	NickName       string        `gorm:"column:nick_name" form:"nick_name" json:"nick_name" comment:"真实姓名"`
	Email          string        `gorm:"column:email" form:"email" json:"email" comment:"邮箱" columnType:"varchar(255)" dataType:"varchar" columnKey:"UNI"`
	Mobile         string        `gorm:"column:mobile" form:"mobile" json:"mobile" comment:"手机号码" columnType:"varchar(11)" dataType:"varchar" columnKey:"UNI"`
	Password       string        `gorm:"column:password" form:"password" json:"password,omitempty" comment:"密码" columnType:"varchar(255)" dataType:"varchar" columnKey:""`
	RoleId         uint          `gorm:"column:role_id" form:"role_id" json:"role_id" comment:"角色ID:2-超级用户,4-普通用户" columnType:"int(10) unsigned" dataType:"int" columnKey:""`
	Status         uint          `gorm:"column:status" form:"status" json:"status" comment:"状态: 1-正常,2-禁用/删除" columnType:"int(10) unsigned" dataType:"int" columnKey:""`
	Avatar         string        `gorm:"column:avatar" form:"avatar" json:"avatar" comment:"用户头像" columnType:"varchar(255)" dataType:"varchar" columnKey:""`
	Remark         string        `gorm:"column:remark" form:"remark" json:"remark" comment:"备注" columnType:"varchar(255)" dataType:"varchar" columnKey:""`
	FriendIds      JsonArrayUint `gorm:"type:json" json:"friend_ids" comment:"json uint 数组"`
	Karma          uint          `json:"karma"`
	CommentIds     JsonArrayUint `gorm:"type:json" json:"comment_ids"`
	HashedPassword string        `gorm:"-" json:"-"`
}

func (m *User) AfterFind() (err error) {
	m.HashedPassword = m.Password
	m.Password = ""
	return
}

//One
func (m *User) One() error {
	return crudOne(m)
}

//All
func (m *User) All(q *PaginationQ) (list *[]User, total uint, err error) {
	list = &[]User{}
	total, err = crudAll(q, db.Model(m), list)
	return
}

//Update
func (m *User) Update() (err error) {
	m.makePassword()
	return db.Model(m).Update(m).Error
}

//Create
func (m *User) Create() (err error) {
	m.Id = 0
	m.makePassword()

	return db.Create(m).Error
}

//Delete
func (m *User) Delete() (err error) {
	if m.Id == 0 {
		return errors.New("resource must not be zero value")
	}
	return crudDelete(m)
}

//Login
func (m *User) Login(ip string) (*jwtObj, error) {
	m.Id = 0
	if m.Password == "" {
		return nil, errors.New("password is required")
	}
	inputPassword := m.Password

	err := db.Where("username = ? or email = ?", m.Username, m.Username).First(&m).Error
	if err != nil {
		return nil, err
	}
	//password is set to bcrypt check
	if err := bcrypt.CompareHashAndPassword([]byte(m.HashedPassword), []byte(inputPassword)); err != nil {
		return nil, err
	}
	m.Password = ""
	data, err := jwtGenerateToken(m)
	return data, err
}

func (m *User) makePassword() {
	if m.Password != "" {
		if bytes, err := bcrypt.GenerateFromPassword([]byte(m.Password), bcrypt.DefaultCost); err != nil {
			logrus.WithError(err).Error("bcrypt making password is failed")
		} else {
			m.Password = string(bytes)
		}
	}
}

func CreateGodUser(user, password string) error {
	m := &User{RoleId: 2, Username: user, Password: password, Email: "dejavuzhou@qq.com", Avatar: "https://tech.mojotv.cn/assets/image/logo01.png"}
	m.makePassword()
	return db.Where("username = ?", m.Username).FirstOrCreate(m).Error
}
