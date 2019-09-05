package ginbro

var tDocYaml = tplNode{
	NameFormat: "doc/swagger.yaml",
	TplContent: `swagger: "2.0"
info:
  description: "A GinBro RESTful APIs"
  version: "1.0.0"
  title: "GinBro RESTful APIs Application"
host: "{{.AppAddr}}"
basePath: "/api/v1"

schemes:
- "http"
paths:
  {{range .Resources}}
  {{if .IsAuthTable}}
  /login:
    post:
      tags:
      - "auth"
      summary: "login by {{.ResourceName}}"
      description: "login by {{.ResourceName}}"
      consumes:
      - "application/json"
      produces:
      - "application/json"
      parameters:
      - in: "body"
        name: "body"
        description: "create {{.ResourceName}}"
        required: true
        schema:
          $ref: "#/definitions/{{.ModelName}}"
      responses:
        200:
          description: "successful operation"
          schema:
            $ref: "#/definitions/{{.ModelName}}Auth"

  {{end}}
  /{{.ResourceName}}:
    get:
      tags:
      - "{{.ResourceName}}"
      summary: "get all {{.ResourceName}} by pagination"
      description: ""
      consumes:
      - "application/json"
      produces:
      - "application/json"
      parameters:
      - name: "where"
        in: "query"
        description: "column:value will use sql LIKE for search eg:id:67 will where id >67 eg2: name:eric => where name LIKE '%eric%'"
        required: false
        type: "array"
        items:
          type: "string"
      - name: "fields"
        in: "query"
        description: "{$tableColumn},{$tableColumn}... "
        required: false
        type: "string"
      - name: "order"
        in: "query"
        description: "eg: id desc, name desc"
        required: false
        type: "string"
      - name: "offset"
        in: "query"
        description: "sql offset eg: 10"
        required: false
        type: "integer"
      - name: "limit"
        in: "query"
        default: "2"
        description: "limit returning object count"
        required: false
        type: "integer"

      responses:
        200:
          description: "successful operation"
          schema:
            type: "array"
            items:
              $ref: "#/definitions/{{.ModelName}}Pagination"
    post:
      tags:
      - "{{.ResourceName}}"
      summary: "create {{.ResourceName}}"
      description: "create {{.ResourceName}}"
      consumes:
      - "application/json"
      produces:
      - "application/json"
      parameters:
      - in: "body"
        name: "body"
        description: "create {{.ResourceName}}"
        required: true
        schema:
          $ref: "#/definitions/{{.ModelName}}"
      responses:
        200:
          description: "successful operation"
          schema:
            $ref: "#/definitions/ApiResponse"

    patch:
      tags:
      - "{{.ResourceName}}"
      summary: "update {{.ResourceName}}"
      description: "update {{.ResourceName}}"
      consumes:
      - "application/json"
      produces:
      - "application/json"
      parameters:
      - in: "body"
        name: "body"
        description: "create {{.ResourceName}}"
        required: true
        schema:
          $ref: "#/definitions/{{.ModelName}}"
      responses:
        200:
          description: "successful operation"
          schema:
            $ref: "#/definitions/ApiResponse"

  /{{.ResourceName}}/{ID}:
    get:
      tags:
      - "{{.ResourceName}}"
      summary: "get a {{.ResourceName}} by ID"
      description: "get a {{.ResourceName}} by ID"
      consumes:
      - "application/json"
      produces:
      - "application/json"
      parameters:
      - name: "ID"
        in: "path"
        description: "ID of {{.ResourceName}} to return"
        required: true
        type: "integer"
        format: "int64"
      responses:
        200:
          description: "successful operation"
          schema:
            $ref: "#/definitions/{{.ModelName}}"
    delete:
      tags:
      - "{{.ResourceName}}"
      summary: "Destroy a {{.ResourceName}} by ID"
      description: "delete a {{.ResourceName}} by ID"
      consumes:
      - "application/json"
      produces:
      - "application/json"
      parameters:
      - name: "ID"
        in: "path"
        description: "ID of {{.ResourceName}} to delete"
        required: true
        type: "integer"
        format: "int64"
      responses:
        200:
          description: "successful operation"
          schema:
            $ref: "#/definitions/ApiResponse"
  {{end}}


definitions:
  {{range $table := .Resources}}
  {{ if $table.IsAuthTable }}
  {{$table.ModelName}}Auth:
    type: "object"
    properties:
      {{range $row := $table.Properties}}
      {{$row.ColumnName}}:
        type: "{{$row.SwaggerType}}"
        description: "{{$row.ColumnComment}}"
        format: "{{$row.SwaggerFormat}}"
        {{end}}
      token:
        type: "string"
        description: "jwt token"
        format: "string"
      expire:
        type: "string"
        description: "jwt token expire time"
        format: "date-time"
      expire_ts:
        type: "integer"
        description: "expire timestamp unix"
        format: "int64"
  {{end}}
  {{$table.ModelName}}:
    type: "object"
    properties:
    {{range $row := $table.Properties}}
      {{$row.ColumnName}}:
        type: "{{$row.SwaggerType}}"
        description: "{{$row.ColumnComment}}"
        format: "{{$row.SwaggerFormat}}"
      {{end}}
  {{$table.ModelName}}Pagination:
    type: "object"
    properties:
      code:
        type: "integer"
        description: "json repose code"
        format: "int32"
      total:
        type: "integer"
        description: "total numbers"
        format: "int32"
      offset:
        type: "integer"
        description: "offset"
        format: "int32"
      limit:
        type: "integer"
        description: "limit"
        format: "int32"
      list:
        type: "array"
        items:
          $ref: "#/definitions/{{$table.ModelName}}"
{{end}}
  ApiResponse:
    type: "object"
    properties:
      code:
        type: "integer"
        format: "int32"
      msg:
        type: "string"
externalDocs:
  description: "Find out more about Swagger"
  url: "http://swagger.io"
`,
}
var tReadme = tplNode{
	NameFormat: "readme.md",
	TplContent: `
# A GinBro RESTful APIs

## Recommend Go version > 1.12
- for Chinese users: set env GOPROXY=https://goproxy.io
- run: go tidy
    
## Usage
- [swagger DOC ](http://{{.AppAddr}}/doc)_[BACKQUOTE]_http://{{.AppAddr}}/swagger/_[BACKQUOTE]_
- [static ](http://{{.AppAddr}})_[BACKQUOTE]_http://{{.AppAddr}}_[BACKQUOTE]_
- [GinbroApp INFO ](http://{{.AppAddr}}/GinbroApp/info)_[BACKQUOTE]_http://{{.AppAddr}}/GinbroApp/info_[BACKQUOTE]_
- API baseURL : _[BACKQUOTE]_http://{{.AppAddr}}/api/v1_[BACKQUOTE]_

## Info
- table'schema which has no "ID","id","ID" or "iD" will not generate model or route.
- the column which type is json value must be a string which is able to decode to a JSON,when call POST or PATCH.
## Thanks
- [swagger Specification](https://swagger.io/specification/)
- [gin-gonic/gin](https://github.com/gin-gonic/gin)
- [GORM](http://gorm.io/)
- [viper](https://github.com/spf13/viper)
- [cobra](https://github.com/spf13/cobra#getting-started)
- [go-redis](https://go get github.com/go-redis/redis)
`,
}

var tModelObj = tplNode{
	NameFormat: "model/m_%s.go",
	TplContent: `
package model

import (
	"errors"
	"time"
	{{if .IsAuthTable}}"fmt"
	"{{.AppPkg}}/config"
	"golang.org/x/crypto/bcrypt"
	"github.com/sirupsen/logrus"
	{{end}}
)

var _ = time.Thursday
//{{.ModelName}}
type {{.ModelName}} struct {
	{{range .Properties}}
	 {{.ModelProp}}      {{.ModelType}}         _[BACKQUOTE]_{{.ModelTag}}_[BACKQUOTE]_{{end}}
}
//TableName
func (m *{{.ModelName}}) TableName() string {
	return "{{.TableName}}"
}
//One
func (m *{{.ModelName}}) One() (one *{{.ModelName}}, err error) {
	one = &{{.ModelName}}{}
	err = crudOne(m, one)
	return
}
//All
func (m *{{.ModelName}}) All(q *PaginationQuery) (list *[]{{.ModelName}}, total uint, err error) {
	list = &[]{{.ModelName}}{}
	total, err = crudAll(m, q, list)
	return
}
//Update
func (m *{{.ModelName}}) Update() (err error) {
	where := {{.ModelName}}{Id: m.Id}
	m.Id = 0
	{{if .IsAuthTable }}m.makePassword()
	{{end}}
	return crudUpdate(m, where)
}
//CreateUserOfRole
func (m *{{.ModelName}}) CreateUserOfRole() (err error) {
	m.Id = 0
    {{if .IsAuthTable }}m.makePassword()
    {{end}}
	return db.CreateUserOfRole(m).Error
}
//Delete
func (m *{{.ModelName}}) Delete() (err error) {
	if m.Id == 0 {
		return errors.New("resource must not be zero value")
	}
	return crudDelete(m)
}
{{if .IsAuthTable }}

//Login
func (m *{{.ModelName}}) Login(ip string) (*jwtObj, error) {
	m.Id = 0
	if m.{{.PasswordPropertyName}} == "" {
		return nil, errors.New("password is required")
	}
	inputPassword := m.{{.PasswordPropertyName}}
	m.{{.PasswordPropertyName}} = ""
	loginTryKey := "login:" + ip
	loginRetries, _ := mem.GetUint(loginTryKey)
	if loginRetries > uint(config.GetInt("GinbroApp.login_try")) {
		memExpire := config.GetInt("GinbroApp.mem_expire_min")
		return nil, fmt.Errorf("for too many wrong login retries the %s will ban for login in %d minitues", ip, memExpire)
	}
	//you can implement more detailed login retry rule
	//for i don't know what your login username i can't implement the ip+username rule in my boilerplate project
	// about username and ip retry rule

	err := db.Where(m).First(&m).Error
	if err != nil {
		//username fail ip retries add 5
		loginRetries = loginRetries + 5
		mem.Set(loginTryKey, loginRetries)
		return nil, err
	}
	//password is set to bcrypt check
	if err := bcrypt.CompareHashAndPassword([]byte(m.{{.PasswordPropertyName}}), []byte(inputPassword)); err != nil {
		// when password failed reties will add 1
		loginRetries = loginRetries + 1
		mem.Set(loginTryKey, loginRetries)
		return nil, err
	}
    m.{{.PasswordPropertyName}} = ""
	key := fmt.Sprintf("login:%d", m.Id)

	//save login user  into the memory store

    data ,err := jwtGenerateToken(m)
    mem.Set(key, data)
    return data,err
}

func (m *{{.ModelName}}) makePassword() {
	if m.{{.PasswordPropertyName}} != "" {
		if bytes, err := bcrypt.GenerateFromPassword([]byte(m.{{.PasswordPropertyName}}), bcrypt.DefaultCost); err != nil {
			logrus.WithError(err).Error("bcrypt making password is failed")
		} else {
			m.{{.PasswordPropertyName}} = string(bytes)
		}
	}
}

{{end}}

`,
}

var tModelJwt = tplNode{
	NameFormat: "model/m_jwt.go",
	TplContent: `
package model

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"{{.AppPkg}}/config"
	"time"
)

func jwtGenerateToken(m *{{.ModelName}}) (*jwtObj, error) {
	m.{{.PasswordPropertyName}} = ""
	expireAfterTime := time.Hour * time.Duration(config.GetInt("GinbroApp.jwt_expire_hour"))
	iss := config.GetString("GinbroApp.name")
	appSecret := config.GetString("GinbroApp.secret")
	expireTime := time.Now().Add(expireAfterTime)
	stdClaims := jwt.StandardClaims{
		ExpiresAt: expireTime.Unix(),
		IssuedAt:  time.Now().Unix(),
		Id:        fmt.Sprintf("%d", m.Id),
		Issuer:    iss,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, stdClaims)
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(appSecret))
	if err != nil {
		logrus.WithError(err).Fatal("config is wrong, can not generate jwt")
	}
	data := &jwtObj{     {{.ModelName}}: *m, Token: tokenString, Expire: expireTime, ExpireTs: expireTime.Unix()}
	return data, err
}

type jwtObj struct {
	{{.ModelName}}
	Token    string    _[BACKQUOTE]_json:"token"_[BACKQUOTE]_
	Expire   time.Time _[BACKQUOTE]_json:"expire"_[BACKQUOTE]_
	ExpireTs int64     _[BACKQUOTE]_json:"expire_ts"_[BACKQUOTE]_
}
//JwtParseUser
func JwtParseUser(tokenString string) (*{{.ModelName}}, error) {
	if tokenString == "" {
		return nil, errors.New("no token is found in Authorization Bearer")
	}
	claims := jwt.StandardClaims{}
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		secret := config.GetString("GinbroApp.secret")
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims.VerifyExpiresAt(time.Now().Unix(), true) == false {
		return nil, errors.New("token is expired")
	}
	appName := config.GetString("GinbroApp.name")
	if !claims.VerifyIssuer(appName, true) {
		return nil, errors.New("token's issuer is wrong,greetings Hacker")
	}
	key := fmt.Sprintf("login:%s", claims.Id)
	jwtObj, err := mem.GetJwtObj(key)
	if err != nil {
		return nil, err
	}
	return &jwtObj.{{.ModelName}}, err
}
//GetJwtObj
func (s *memoryStore) GetJwtObj(id string) (value *jwtObj, err error) {
	vv, err := s.Get(id, false)
	if err != nil {
		return nil, err
	}
	value, ok := vv.(*jwtObj)
	if ok {
		return value, nil
	}
	return nil, errors.New("mem:has value of this id, but is not type of *jwtObj")
}

`,
}

var tHandlersObj = tplNode{
	NameFormat: "handlers/h_%s.go",
	TplContent: `
package handlers

import (
	"{{.AppPkg}}/model"
	"github.com/gin-gonic/gin"
)

func init() {
	groupApi.GET("{{.ResourceName}}",{{if .IsAuthTable}}jwtMiddleware,{{end}} {{.HandlerName}}All)
	{{if .HasId}}groupApi.GET("{{.ResourceName}}/:id", {{if .IsAuthTable}}jwtMiddleware,{{end}} {{.HandlerName}}One){{end}}
	groupApi.POST("{{.ResourceName}}", {{if .IsAuthTable}}jwtMiddleware,{{end}} {{.HandlerName}}CreateUserOfRole)
	groupApi.PATCH("{{.ResourceName}}", {{if .IsAuthTable}}jwtMiddleware,{{end}} {{.HandlerName}}Update)
	{{if .HasId}}groupApi.DELETE("{{.ResourceName}}/:id", {{if .IsAuthTable}}jwtMiddleware,{{end}} {{.HandlerName}}Delete){{end}}
}
//All
func {{.HandlerName}}All(c *gin.Context) {
	mdl := model.{{.ModelName}}{}
	query := &model.PaginationQuery{}
	err := c.ShouldBindQuery(query)
	if handleError(c, err) {
		return
	}
	list, total, err := mdl.All(query)
	if handleError(c, err) {
		return
	}
	jsonPagination(c, list, total, query)
}
{{if .HasId}}
//One
func {{.HandlerName}}One(c *gin.Context) {
	var mdl model.{{.ModelName}}
	id, err := parseParamID(c)
	if handleError(c, err) {
		return
	}
	mdl.Id = id
	data, err := mdl.One()
	if handleError(c, err) {
		return
	}
	jsonData(c, data)
}
{{end}}
//CreateUserOfRole
func {{.HandlerName}}CreateUserOfRole(c *gin.Context) {
	var mdl model.{{.ModelName}}
	err := c.ShouldBind(&mdl)
	if handleError(c, err) {
		return
	}
	err = mdl.CreateUserOfRole()
	if handleError(c, err) {
		return
	}
	jsonData(c, mdl)
}
//Update
func {{.HandlerName}}Update(c *gin.Context) {
	var mdl model.{{.ModelName}}
	err := c.ShouldBind(&mdl)
	if handleError(c, err) {
		return
	}
	err = mdl.Update()
	if handleError(c, err) {
		return
	}
	jsonSuccess(c)
}
{{if .HasId}}
//Delete
func {{.HandlerName}}Delete(c *gin.Context) {
	var mdl model.{{.ModelName}}
	id, err := parseParamID(c)
	if handleError(c, err) {
		return
	}
	mdl.Id = id
	err = mdl.Delete()
	if handleError(c, err) {
		return
	}
	jsonSuccess(c)
}
{{end}}
`,
}

var tStaticIndex = tplNode{
	"static/index.html",
	`
<!DOCTYPE html><html lang="en">
<head><meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
  <meta name="viewport" content="width=device-width">
  <title>Golang+gin+gorm+sql created by felix ginbro</title>
<body>
  <h2 style="text-align: center">put front end files into this folder</h2>
</body>
</html>
`,
}

var tHandlersObjBare = tplNode{
	NameFormat: "handlers/h_%s.go",
	TplContent: `
package handlers

import (
	"github.com/gin-gonic/gin"
)

func init() {
	groupApi.GET("{{.ResourceName}}",{{.HandlerName}}All)
	groupApi.GET("{{.ResourceName}}/:id", {{.HandlerName}}One)
	groupApi.POST("{{.ResourceName}}", {{.HandlerName}}CreateUserOfRole)
	groupApi.PATCH("{{.ResourceName}}", {{.HandlerName}}Update)
	groupApi.DELETE("{{.ResourceName}}/:id", {{.HandlerName}}Delete)
}
//All
func {{.HandlerName}}All(c *gin.Context) {

}
//One
func {{.HandlerName}}One(c *gin.Context) {


}
//CreateUserOfRole
func {{.HandlerName}}CreateUserOfRole(c *gin.Context) {

}
//Update
func {{.HandlerName}}Update(c *gin.Context) {

}
//Delete
func {{.HandlerName}}Delete(c *gin.Context) {

}

`,
}

var tConfigToml = tplNode{
	NameFormat: "config.toml",
	TplContent: `
[GinbroApp]
    name = "ginBro"
    addr ="{{.AppAddr}}" # eg1: www.mojotv.cn     eg:localhost:3333 eg1:127.0.0.1:88
    secret = "{{.AppSecret}}"
    env = "local" # only allows local/dev/test/prod
    log_level = "error" # only allows debug info warn error fatal panic
    enable_not_found = true # if true and static_path is not empty string, all not found route will serve static/index.html
    enable_swagger = true
    enable_cors = true  # true will case 403 error in swaggerUI  may cause api perform decrease
    enable_sql_log = true # show gorm sql in terminal
    enable_https = false # if addr is a domain enable_https will works
    enable_cron = false # is enable buildin schedule job
    time_zone = "Asia/Shanghai"
    api_prefix = "v1" #  api_prefix could be empty string,            the api uri will be api/v1/resource
    static_path = "./static/"  # path must be an absolute path or relative to the go-build-executable file, may cause api perform decrease
    mem_expire_min = 60 # memory cache expire in 60 minutes
    mem_max_count = 1024000 # memory cache maxium store count
    login_try = 100 # after 100 times login failure the IP will be ban for mem_expire_min(default 600min), wrong username costs 5 times, wrong password costs 1 time,
    jwt_expire_hour = 24 # jwt expire in 24 hours
[db]
    type = "{{.DbType}}"
    addr = "{{.DbAddr}}"
    user = "{{.DbUser}}"
    password = "{{.DbPassword}}"
    database = "{{.DbName}}"
    charset = "{{.DbChar}}"
[redis]
    addr = "" # 127.0.0.1:6379 empty string will not init the redis db in model package
    password = ""
    db_idx = 0


# the init config has not impelement yet
[init]
    user_email= "admin@ginbro.com" # if not exist, create a user with the bcrypt password, if the value is empty will do nothing
    user_password = "123123" # print the bcrypted password in console for you to paste into mysql auth_table.password column
`,
}

var tModelDbMem = tplNode{
	NameFormat: "model/db_mem.go",
	TplContent: `
package model

import (
	"container/list"
	"errors"
	"github.com/sirupsen/logrus"
	"{{.AppPkg}}/config"
	"sync"
	"time"
)

var mem *memoryStore

func init() {
	mem = new(memoryStore)
	mem.digitsById = make(map[string]interface{})
	mem.idByTime = list.New()
	maxCount := config.GetInt("GinbroApp.mem_max_count")
	if maxCount <= 1024 {
		maxCount = 1024
	}
	mem.collectNum = maxCount
	expireIn := config.GetInt("GinbroApp.mem_expire_min")
	if expireIn <= 0 {
		expireIn = 30
	}
	mem.expiration = time.Minute * time.Duration(expireIn)
}

type idByTimeValue struct {
	timestamp time.Time
	id        string
}

// memoryStore is an internal store for captcha ids and their values.
type memoryStore struct {
	sync.RWMutex
	digitsById map[string]interface{}
	idByTime   *list.List
	// Number of items stored since last collection.
	numStored int
	// Number of saved items that triggers collection.
	collectNum int
	// Expiration time of captchas.
	expiration time.Duration
}

func (s *memoryStore) Set(id string, value interface{}) {
	s.Lock()
	s.digitsById[id] = value
	s.idByTime.PushBack(idByTimeValue{time.Now(), id})
	s.numStored++
	s.Unlock()
	if s.numStored > s.collectNum {
		go s.collect()
	}
}

func (s *memoryStore) Get(id string, clear bool) (value interface{}, err error) {
	if !clear {
		// When we don't need to clear captcha, acquire read lock.
		s.RLock()
		defer s.RUnlock()
	} else {
		s.Lock()
		defer s.Unlock()
	}
	value, ok := s.digitsById[id]
	if !ok {
		return nil, errors.New("value not found")
	}
	if clear {
		delete(s.digitsById, id)
	}
	return
}

func (s *memoryStore) collect() {
	logrus.Warn("memory store collect function has been called some value will be lost")
	now := time.Now()
	s.Lock()
	defer s.Unlock()
	s.numStored = 0
	for e := s.idByTime.Front(); e != nil; {
		ev, ok := e.Value.(idByTimeValue)
		if !ok {
			return
		}
		if ev.timestamp.Add(s.expiration).Before(now) {
			delete(s.digitsById, ev.id)
			next := e.Next()
			s.idByTime.Remove(e)
			e = next
		} else {
			return
		}
	}
}

func (s *memoryStore) GetUint(id string) (value uint, err error) {
	vv, err := s.Get(id, false)
	if err != nil {
		return 0, err
	}
	value, ok := vv.(uint)
	if ok {
		return value, nil
	}
	return 0, errors.New("mem:has value of this id, but is not type of uint")
}

`,
}

var tModelHelper = tplNode{
	NameFormat: "model/helper.go",
	TplContent: `
package model

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"reflect"
	"strconv"
	"strings"
)

//PaginationQuery gin handler query binding struct
type PaginationQuery struct {
	Where  string _[BACKQUOTE]_form:"where"_[BACKQUOTE]_
	Fields string _[BACKQUOTE]_form:"fields"_[BACKQUOTE]_
	Order  string _[BACKQUOTE]_form:"order"_[BACKQUOTE]_
	Offset uint   _[BACKQUOTE]_form:"offset"_[BACKQUOTE]_
	Limit  uint   _[BACKQUOTE]_form:"limit"_[BACKQUOTE]_
}

//String to string
func (pq *PaginationQuery) String() string {
	return fmt.Sprintf("w=%v_f=%s_o=%s_of=%d_l=%d", pq.Where, pq.Fields, pq.Order, pq.Offset, pq.Limit)
}

func crudAll(m interface{}, q *PaginationQuery, list interface{}) (total uint, err error) {
	var tx *gorm.DB
	total, tx = getResourceCount(m, q)
	if q.Fields != "" {
		columns := strings.Split(q.Fields, ",")
		if len(columns) > 0 {
			tx = tx.Select(q.Fields)
		}
	}
	if q.Order != "" {
		tx = tx.Order(q.Order)
	}
	if q.Offset > 0 {
		tx = tx.Offset(q.Offset)
	}
	if q.Limit <= 0 {
		q.Limit = 15
	}
	err = tx.Limit(q.Limit).Find(list).Error
	return
}

func crudOne(m interface{}, one interface{}) (err error) {
	if db.Where(m).First(one).RecordNotFound() {
		return errors.New("resource is not found")
	}
	return nil
}

func crudUpdate(m interface{}, where interface{}) (err error) {
	db := db.Model(where).Updates(m)
	if err = db.Error; err != nil {
		return
	}
	if db.RowsAffected != 1 {
		return errors.New("id is invalid and resource is not found")
	}
	return nil
}

func crudDelete(m interface{}) (err error) {
	//WARNING When delete a record, you need to ensure it’s primary field has value, and GORM will use the primary key to delete the record, if primary field’s blank, GORM will delete all records for the model
	//primary key must be not zero value
	db := db.Delete(m)
	if err = db.Error; err != nil {
		return
	}
	if db.RowsAffected != 1 {
		return errors.New("resource is not found to destroy")
	}
	return nil
}
func getResourceCount(m interface{}, q *PaginationQuery) (uint, *gorm.DB) {
	var tx = db.Model(m)
	conditions := strings.Split(q.Where, ",")
	for _, val := range conditions {
		w := strings.SplitN(val, ":", 2)
		if len(w) == 2 {
			bindKey, bindValue := w[0], w[1]
			if intV, err := strconv.ParseInt(bindValue, 10, 64); err == nil {
				// bind value is int
				field := fmt.Sprintf("_[BACKQUOTE]_%s_[BACKQUOTE]_ > ?", bindKey)
				tx = tx.Where(field, intV)
			} else if fV, err := strconv.ParseFloat(bindValue, 64); err == nil {
				// bind value is float
				field := fmt.Sprintf("_[BACKQUOTE]_%s_[BACKQUOTE]_ > ?", bindKey)
				tx = tx.Where(field, fV)
			} else if bindValue != "" {
				// bind value is string
				field := fmt.Sprintf("_[BACKQUOTE]_%s_[BACKQUOTE]_ LIKE ?", bindKey)
				sV := fmt.Sprintf("%%%s%%", bindValue)
				tx = tx.Where(field, sV)
			}
		}
	}
	modelName := getType(m)
	rKey := redisPrefix + modelName + q.String() + "_count"
	v, err := mem.GetUint(rKey)
	if err != nil {
		var count uint
		tx.Count(&count)
		mem.Set(rKey, count)
		return count, tx
	}
	return v, tx
}

func getType(v interface{}) string {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	}
	return t.Name()
}

`,
}

var tModelDbSql = tplNode{
	NameFormat: "model/db_sql.go",
	TplContent: `
package model

import (
	"fmt"
    _ "github.com/jinzhu/gorm/dialects/mssql"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    _ "github.com/jinzhu/gorm/dialects/postgres"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"{{.AppPkg}}/config"
	"strings"
	"errors"
)

var db *gorm.DB

func init() {
	if gormDB, err := createDatabase(); err == nil {
		db = gormDB
	} else {
		logrus.WithError(err).Fatalln("create database connection failed")
	}
	//enable Gorm mysql log
	if flag := config.GetBool("GinbroApp.enable_sql_log"); flag {
		db.LogMode(flag)
		//f, err := os.OpenFile("mysql_gorm.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		//if err != nil {
		//	logrus.WithError(err).Fatalln("could not create mysql gorm log file")
		//}
		//logger :=  New(f,"", Ldate)
		//db.SetLogger(logger)
	}
	//db.AutoMigrate()

}

//Close clear db collection
func Close() {
	if db != nil {
		db.Close()
	}
	if redisDB != nil {
		redisDB.Close()
	}
}

func createDatabase() (*gorm.DB,error) {
	dbType := config.GetString("db.type")
	dbAddr := config.GetString("db.addr")
	dbName := config.GetString("db.database")
	dbUser := config.GetString("db.user")
	dbPassword := config.GetString("db.password")
	dbCharset := config.GetString("db.charset")
	conn := ""
	switch dbType {
	case "mysql":
		conn = fmt.Sprintf("%s:%s@(%s)/%s?charset=%s&parseTime=True&loc=Local", dbUser,dbPassword, dbAddr, dbName,dbCharset)
	case "sqlite":
		conn = dbAddr
	case "mssql":
		return nil,errors.New("TODO:suport sqlServer")
	case "postgres":
		hostPort := strings.Split(dbAddr, ":")
		if len(hostPort) == 2{
			return nil,errors.New("db.addr must be like this host:ip")
		}
		conn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", hostPort[0], hostPort[1], dbUser, dbName, dbPassword)
	default:
		return nil,fmt.Errorf("database type %s is not supported by felix ginrbo",dbType)
	}
	return gorm.Open(dbType,conn)
}
`,
}

var tMain = tplNode{
	NameFormat: "main.go",
	TplContent: `
package main

import (
	"{{.AppPkg}}/handlers"
	"{{.AppPkg}}/tasks"
	"{{.AppPkg}}/config"
)

func main() {
	if config.GetBool("GinbroApp.enable_cron") {
		go tasks.RunTasks()
	}
	defer handlers.Close()
	handlers.ServerRun()
}

`,
}

var tHandlerMiddlewareJwt = tplNode{
	NameFormat: "handlers/h_middleware_jwt.go",
	TplContent: `
package handlers

import (
	"{{.AppPkg}}/model"
	"net/http"
	"github.com/gin-gonic/gin"
	"strings"
)

var jwtMiddleware = jwtCheck()

const tokenPrefix = "Bearer "
const bearerLength = len(tokenPrefix)

func jwtCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		hToken := c.GetHeader("Authorization")
		if len(hToken) < bearerLength {
			c.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{"msg": "header Authorization has not Bearer token"})
			return
		}
		token := strings.TrimSpace(hToken[bearerLength:])
		user, err := model.JwtParseUser(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusPreconditionFailed, gin.H{"msg": err.Error()})
			return
		}
		//store the user Model in the context
		c.Set("user", user)
		c.Next()
		// after request
	}
}

`,
}

var tHandlerLogin = tplNode{
	NameFormat: "handlers/h_login.go",
	TplContent: `
package handlers

import (
	"{{.AppPkg}}/model"
	"github.com/gin-gonic/gin"
)

func init() {
	groupApi.POST("login", login)
}

func login(c *gin.Context) {
	var mdl model.{{.ModelName }}
	err := c.ShouldBind(&mdl)
	if handleError(c, err) {
		return
	}
	ip := c.ClientIP()
	data, err := mdl.Login(ip)
	if handleError(c, err) {
		return
	}
	jsonData(c, data)
}
`,
}

var tHandlerHelper = tplNode{
	NameFormat: "handlers/helper.go",
	TplContent: `
package handlers

import (
	"errors"
	"{{.AppPkg}}/model"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func jsonError(c *gin.Context, msg string) {
	c.AbortWithStatusJSON(200, gin.H{"code": 0, "msg": msg})
}
func jsonData(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{"code": 1, "data": data})
}
func jsonPagination(c *gin.Context, list interface{}, total uint, query *model.PaginationQuery) {
	c.JSON(200, gin.H{"code": 1, "data": list, "total": total, "offset": query.Offset, "limit": query.Limit})
}
func jsonSuccess(c *gin.Context) {
	c.JSON(200, gin.H{"code": 1, "msg": "success"})
}

func handleError(c *gin.Context, err error) bool {
	if err != nil {
		jsonError(c, err.Error())
		return true
	}
	return false
}

func parseParamID(c *gin.Context) (uint, error) {
	id := c.Param("id")
	parseId, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return 0, errors.New("id must be an unsigned int")
	}
	return uint(parseId), nil
}

func enableCorsMiddleware() {
		r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))
}

`,
}

var tHandlerGin = tplNode{
	NameFormat: "handlers/gin.go",
	TplContent: `
package handlers

import (
	"fmt"
	"{{.AppPkg}}/model"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"{{.AppPkg}}/config"
	"log"
	"path"
)

var r = gin.Default()
var groupApi *gin.RouterGroup

//in the same package init executes in file'name alphabet order
func init() {
	if config.GetBool("GinbroApp.enable_cors") {
		enableCorsMiddleware()
	}
	if sp := config.GetString("GinbroApp.static_path"); sp != "" {
		r.Use(static.Serve("/", static.LocalFile(sp, true)))
		if config.GetBool("GinbroApp.enable_not_found") {
			r.NoRoute(func(c *gin.Context) {
				file := path.Join(sp, "index.html")
				c.File(file)
			})
		}
	}

	if config.GetBool("GinbroApp.enable_swagger") && config.GetString("GinbroApp.env") != "prod" {
		//add edit your own swagger.doc.yml file in ./swagger/doc.yml
		//generateSwaggerDocJson()
		r.Static("doc", "./doc")
	}
	prefix := config.GetString("GinbroApp.api_prefix")
	api := "api"
	if prefix != "" {
		api = fmt.Sprintf("%s/%s", api, prefix)
	}
	groupApi = r.Group(api)

	if config.GetString("GinbroApp.env") != "prod" {
		r.GET("/GinbroApp/info", func(c *gin.Context) {
			c.JSON(200, config.GetStringMapString("GinbroApp"))
		})
	}

}

//ServerRun start the gin server
func ServerRun() {

	addr := config.GetString("GinbroApp.addr")
	if config.GetBool("GinbroApp.enable_https") {
		log.Fatal(autotls.Run(r, addr))
	} else {
		log.Printf("visit http://%s/doc for RESTful APIs Document", addr)
		log.Printf("visit http://%s/ for front-end static html files", addr)
		log.Printf("visit http://%s/GinbroApp/info for GinbroApp info only on not-prod mode", addr)
		r.Run(addr)
	}
}

//Close gin GinbroApp
func Close() {
	model.Close()
}

`,
}

var tConfig = tplNode{
	NameFormat: "config/viper.go",
	TplContent: `
package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config") // name of config file (without extension)
	//viper.AddConfigPath("/etc/appname/")   // path to look for the config file in
	//viper.AddConfigPath("$HOME/.appname")  // call multiple times to add many search paths
	viper.AddConfigPath(".")    // optionally look for config in the working directory
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		logrus.WithError(err).Error("application configuration'initialization is failed ")
	}
}

// GetString returns the value associated with the key as a string.
func GetString(key string) string {
	return viper.GetString(key)
}

// GetInt returns the value associated with the key as an integer.
func GetInt(key string) int {
	return viper.GetInt(key)
}

// GetBool returns the value associated with the key as a boolean.
func GetBool(key string) bool {
	return viper.GetBool(key)
}

func GetStringMapString(key string) map[string]string {
	return viper.GetStringMapString(key)
}

`,
}

var tGitIgnore = tplNode{
	NameFormat: ".gitignore",
	TplContent: `
*.exe
*.log
*.toml
.idea/*
.idea
.vscode/*
.vscode
`,
}
var tTaskCore = tplNode{
	NameFormat: "tasks/core.go",
	TplContent: `
package tasks

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Time location, default set by the time.Local (*time.Location)
var loc = time.Local

// Change the time location
func ChangeLoc(newLocation *time.Location) {
	loc = newLocation
}

// Max number of jobs, hack it if you need.
const MAXJOBNUM = 10000

type Job struct {
	// pause interval * unit bettween runs
	interval uint64
	
	// the job jobFunc to run, func[jobFunc]
	jobFunc string
	// time units, ,e.g. 'minutes', 'hours'...
	unit string
	// optional time at which this job runs
	atTime string
	
	// datetime of last run
	lastRun time.Time
	// datetime of next run
	nextRun time.Time
	// cache the period between last an next run
	period time.Duration
	
	// Specific day of the week to start on
	startDay time.Weekday
	
	// Map for the function task store
	funcs map[string]interface{}
	
	// Map for function and  params of function
	fparams map[string]([]interface{})
}

// CreateUserOfRole a new job with the time interval.
func NewJob(intervel uint64) *Job {
	return &Job{
		intervel,
		"", "", "",
		time.Unix(0, 0),
		time.Unix(0, 0), 0,
		time.Sunday,
		make(map[string]interface{}),
		make(map[string]([]interface{})),
	}
}

// True if the job should be run now
func (j *Job) shouldRun() bool {
	return time.Now().After(j.nextRun)
}

//Run the job and immediately reschedule it
func (j *Job) run() (result []reflect.Value, err error) {
	f := reflect.ValueOf(j.funcs[j.jobFunc])
	params := j.fparams[j.jobFunc]
	if len(params) != f.Type().NumIn() {
		err = errors.New("the number of param is not adapted")
		return
	}
	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	result = f.Call(in)
	j.lastRun = time.Now()
	j.scheduleNextRun()
	return
}

// for given function fn, get the name of function.
func getFunctionName(fn interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf((fn)).Pointer()).Name()
}

// Specifies the jobFunc that should be called every time the job runs
//
func (j *Job) Do(jobFun interface{}, params ...interface{}) {
	typ := reflect.TypeOf(jobFun)
	if typ.Kind() != reflect.Func {
		panic("only function can be schedule into the job queue.")
	}
	
	fname := getFunctionName(jobFun)
	j.funcs[fname] = jobFun
	j.fparams[fname] = params
	j.jobFunc = fname
	//schedule the next run
	j.scheduleNextRun()
}

func formatTime(t string) (hour, min int, err error) {
	var er = errors.New("time format error")
	ts := strings.Split(t, ":")
	if len(ts) != 2 {
		err = er
		return
	}
	
	hour, err = strconv.Atoi(ts[0])
	if err != nil {
		return
	}
	min, err = strconv.Atoi(ts[1])
	if err != nil {
		return
	}
	
	if hour < 0 || hour > 23 || min < 0 || min > 59 {
		err = er
		return
	}
	return hour, min, nil
}

//	s.Every(1).Day().At("10:30").Do(task)
//	s.Every(1).Monday().At("10:30").Do(task)
func (j *Job) At(t string) *Job {
	hour, min, err := formatTime(t)
	if err != nil {
		panic(err)
	}
	
	// time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	mock := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), int(hour), int(min), 0, 0, loc)
	
	if j.unit == "days" {
		if time.Now().After(mock) {
			j.lastRun = mock
		} else {
			j.lastRun = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-1, hour, min, 0, 0, loc)
		}
	} else if j.unit == "weeks" {
		if j.startDay != time.Now().Weekday() || (time.Now().After(mock) && j.startDay == time.Now().Weekday()) {
			i := mock.Weekday() - j.startDay
			if i < 0 {
				i = 7 + i
			}
			j.lastRun = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-int(i), hour, min, 0, 0, loc)
		} else {
			j.lastRun = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-7, hour, min, 0, 0, loc)
		}
	}
	return j
}

//Compute the instant when this job should run next
func (j *Job) scheduleNextRun() {
	if j.lastRun == time.Unix(0, 0) {
		if j.unit == "weeks" {
			i := time.Now().Weekday() - j.startDay
			if i < 0 {
				i = 7 + i
			}
			j.lastRun = time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-int(i), 0, 0, 0, 0, loc)
			
		} else {
			j.lastRun = time.Now()
		}
	}
	
	if j.period != 0 {
		// translate all the units to the Seconds
		j.nextRun = j.lastRun.Add(j.period * time.Second)
	} else {
		switch j.unit {
		case "minutes":
			j.period = time.Duration(j.interval * 60)
			break
		case "hours":
			j.period = time.Duration(j.interval * 60 * 60)
			break
		case "days":
			j.period = time.Duration(j.interval * 60 * 60 * 24)
			break
		case "weeks":
			j.period = time.Duration(j.interval * 60 * 60 * 24 * 7)
			break
		case "seconds":
			j.period = time.Duration(j.interval)
		}
		j.nextRun = j.lastRun.Add(j.period * time.Second)
	}
}

// NextScheduledTime returns the time of when this job is to run next
func (j *Job) NextScheduledTime() time.Time {
	return j.nextRun
}

// the follow functions set the job's unit with seconds,minutes,hours...

// Set the unit with second
func (j *Job) Second() (job *Job) {
	if j.interval != 1 {
		panic("")
	}
	job = j.Seconds()
	return
}

// Set the unit with seconds
func (j *Job) Seconds() (job *Job) {
	j.unit = "seconds"
	return j
}

// Set the unit  with minute, which interval is 1
func (j *Job) Minute() (job *Job) {
	if j.interval != 1 {
		panic("")
	}
	job = j.Minutes()
	return
}

//set the unit with minute
func (j *Job) Minutes() (job *Job) {
	j.unit = "minutes"
	return j
}

//set the unit with hour, which interval is 1
func (j *Job) Hour() (job *Job) {
	if j.interval != 1 {
		panic("")
	}
	job = j.Hours()
	return
}

// Set the unit with hours
func (j *Job) Hours() (job *Job) {
	j.unit = "hours"
	return j
}

// Set the job's unit with day, which interval is 1
func (j *Job) Day() (job *Job) {
	if j.interval != 1 {
		panic("")
	}
	job = j.Days()
	return
}

// Set the job's unit with days
func (j *Job) Days() *Job {
	j.unit = "days"
	return j
}

// s.Every(1).Monday().Do(task)
// Set the start day with Monday
func (j *Job) Monday() (job *Job) {
	if j.interval != 1 {
		panic("")
	}
	j.startDay = 1
	job = j.Weeks()
	return
}

// Set the start day with Tuesday
func (j *Job) Tuesday() (job *Job) {
	if j.interval != 1 {
		panic("")
	}
	j.startDay = 2
	job = j.Weeks()
	return
}

// Set the start day woth Wednesday
func (j *Job) Wednesday() (job *Job) {
	if j.interval != 1 {
		panic("")
	}
	j.startDay = 3
	job = j.Weeks()
	return
}

// Set the start day with thursday
func (j *Job) Thursday() (job *Job) {
	if j.interval != 1 {
		panic("")
	}
	j.startDay = 4
	job = j.Weeks()
	return
}

// Set the start day with friday
func (j *Job) Friday() (job *Job) {
	if j.interval != 1 {
		panic("")
	}
	j.startDay = 5
	job = j.Weeks()
	return
}

// Set the start day with saturday
func (j *Job) Saturday() (job *Job) {
	if j.interval != 1 {
		panic("")
	}
	j.startDay = 6
	job = j.Weeks()
	return
}

// Set the start day with sunday
func (j *Job) Sunday() (job *Job) {
	if j.interval != 1 {
		panic("")
	}
	j.startDay = 0
	job = j.Weeks()
	return
}

//Set the units as weeks
func (j *Job) Weeks() *Job {
	j.unit = "weeks"
	return j
}

// Class Scheduler, the only data member is the list of jobs.
type Scheduler struct {
	// Array store jobs
	jobs [MAXJOBNUM]*Job
	
	// Size of jobs which jobs holding.
	size int
}

// Scheduler implements the sort.Interface{} for sorting jobs, by the time nextRun

func (s *Scheduler) Len() int {
	return s.size
}

func (s *Scheduler) Swap(i, j int) {
	s.jobs[i], s.jobs[j] = s.jobs[j], s.jobs[i]
}

func (s *Scheduler) Less(i, j int) bool {
	return s.jobs[j].nextRun.After(s.jobs[i].nextRun)
}

// CreateUserOfRole a new scheduler
func NewScheduler() *Scheduler {
	return &Scheduler{[MAXJOBNUM]*Job{}, 0}
}

// Get the current runnable jobs, which shouldRun is True
func (s *Scheduler) getRunnableJobs() (running_jobs [MAXJOBNUM]*Job, n int) {
	runnableJobs := [MAXJOBNUM]*Job{}
	n = 0
	sort.Sort(s)
	for i := 0; i < s.size; i++ {
		if s.jobs[i].shouldRun() {
			
			runnableJobs[n] = s.jobs[i]
			//fmt.Println(runnableJobs)
			n++
		} else {
			break
		}
	}
	return runnableJobs, n
}

// Datetime when the next job should run.
func (s *Scheduler) NextRun() (*Job, time.Time) {
	if s.size <= 0 {
		return nil, time.Now()
	}
	sort.Sort(s)
	return s.jobs[0], s.jobs[0].nextRun
}

// Schedule a new periodic job
func (s *Scheduler) Every(interval uint64) *Job {
	job := NewJob(interval)
	s.jobs[s.size] = job
	s.size++
	return job
}

// Run all the jobs that are scheduled to run.
func (s *Scheduler) RunPending() {
	runnableJobs, n := s.getRunnableJobs()
	
	if n != 0 {
		for i := 0; i < n; i++ {
			runnableJobs[i].run()
		}
	}
}

// Run all jobs regardless if they are scheduled to run or not
func (s *Scheduler) RunAll() {
	for i := 0; i < s.size; i++ {
		s.jobs[i].run()
	}
}

// Run all jobs with delay seconds
func (s *Scheduler) RunAllwithDelay(d int) {
	for i := 0; i < s.size; i++ {
		s.jobs[i].run()
		time.Sleep(time.Duration(d))
	}
}

// Remove specific job j
func (s *Scheduler) Remove(j interface{}) {
	i := 0
	for ; i < s.size; i++ {
		if s.jobs[i].jobFunc == getFunctionName(j) {
			break
		}
	}
	
	for j := (i + 1); j < s.size; j++ {
		s.jobs[i] = s.jobs[j]
		i++
	}
	s.size = s.size - 1
}

// Delete all scheduled jobs
func (s *Scheduler) Clear() {
	for i := 0; i < s.size; i++ {
		s.jobs[i] = nil
	}
	s.size = 0
}

// Start all the pending jobs
// Add seconds ticker
func (s *Scheduler) Start() chan bool {
	stopped := make(chan bool, 1)
	ticker := time.NewTicker(1 * time.Second)
	
	go func() {
		for {
			select {
			case <-ticker.C:
				s.RunPending()
			case <-stopped:
				return
			}
		}
	}()
	
	return stopped
}

// The following methods are shortcuts for not having to
// create a Schduler instance

var defaultScheduler = NewScheduler()
var jobs = defaultScheduler.jobs

// Schedule a new periodic job
func Every(interval uint64) *Job {
	return defaultScheduler.Every(interval)
}

// Run all jobs that are scheduled to run
//
// Please note that it is *intended behavior that run_pending()
// does not run missed jobs*. For example, if you've registered a job
// that should run every minute and you only call run_pending()
// in one hour increments then your job won't be run 60 times in
// between but only once.
func RunPending() {
	defaultScheduler.RunPending()
}

// Run all jobs regardless if they are scheduled to run or not.
func RunAll() {
	defaultScheduler.RunAll()
}

// Run all the jobs with a delay in seconds
//
// A delay of 'delay' seconds is added between each job. This can help
// to distribute the system load generated by the jobs more evenly over
// time.
func RunAllwithDelay(d int) {
	defaultScheduler.RunAllwithDelay(d)
}

// Run all jobs that are scheduled to run
func Start() chan bool {
	return defaultScheduler.Start()
}

// Clear
func Clear() {
	defaultScheduler.Clear()
}

// Remove
func Remove(j interface{}) {
	defaultScheduler.Remove(j)
}

// NextRun gets the next running time
func NextRun() (job *Job, time time.Time) {
	return defaultScheduler.NextRun()
}

func RunTasks() {
	// Do jobs with params
	Every(1).Second().Do(taskWithParams, 1, "hello")
	
	// Do jobs without params
	Every(1).Second().Do(task)
	Every(2).Seconds().Do(task)
	Every(1).Minute().Do(task)
	Every(2).Minutes().Do(task)
	Every(1).Hour().Do(task)
	Every(2).Hours().Do(task)
	Every(1).Day().Do(task)
	Every(2).Days().Do(task)
	
	// Do jobs on specific weekday
	Every(1).Monday().Do(task)
	Every(1).Thursday().Do(task)
	
	// function At() take a string like 'hour:min'
	Every(1).Day().At("10:30").Do(task)
	Every(1).Monday().At("18:30").Do(task)
	
	// remove, clear and next_run
	_, time := NextRun()
	fmt.Println(time)
	
	// Remove(task)
	// Clear()
	
	// function Start start all the pending jobs
	<-Start()
	
	// also , you can create a your new scheduler,
	// to run two scheduler concurrently
	s := NewScheduler()
	s.Every(3).Seconds().Do(task)
	<-s.Start()
}


`,
}
var tTaskExample = tplNode{
	NameFormat: "tasks/example.go",
	TplContent: `
package tasks

import "fmt"

//defining schedule task function here
//then add the function in manger.go
func task() {
	fmt.Println("task one is called")
}
func taskWithParams(a int, b string) {
	fmt.Println(a, b)
}

`,
}

var tDocOauth2 = tplNode{
	"doc/oauth-redirect.html",
	`
<!doctype html>
<html lang="en-US">
<body onload="run()">
</body>
</html>
<script>
    'use strict';

    function run() {
        var oauth2 = window.opener.swaggerUIRedirectOauth2;
        var sentState = oauth2.state;
        var redirectUrl = oauth2.redirectUrl;
        var isValid, qp, arr;

        if (/code|token|error/.test(window.location.hash)) {
            qp = window.location.hash.substring(1);
        } else {
            qp = location.search.substring(1);
        }

        arr = qp.split("&")
        arr.forEach(function (v, i, _arr) {
            _arr[i] = '"' + v.replace('=', '":"') + '"';
        })
        qp = qp ? JSON.parse('{' + arr.join() + '}',
                function (key, value) {
                    return key === "" ? value : decodeURIComponent(value)
                }
        ) : {}

        isValid = qp.state === sentState

        if ((
                oauth2.auth.schema.get("flow") === "accessCode" ||
                oauth2.auth.schema.get("flow") === "authorizationCode"
        ) && !oauth2.auth.code) {
            if (!isValid) {
                oauth2.errCb({
                    authId: oauth2.auth.name,
                    source: "auth",
                    level: "warning",
                    message: "Authorization may be unsafe, passed state was changed in server Passed state wasn't returned from auth server"
                });
            }

            if (qp.code) {
                delete oauth2.state;
                oauth2.auth.code = qp.code;
                oauth2.callback({auth: oauth2.auth, redirectUrl: redirectUrl});
            } else {
                let oauthErrorMsg
                if (qp.error) {
                    oauthErrorMsg = "[" + qp.error + "]: " +
                            (qp.error_description ? qp.error_description + ". " : "no accessCode received from the server. ") +
                            (qp.error_uri ? "More info: " + qp.error_uri : "");
                }

                oauth2.errCb({
                    authId: oauth2.auth.name,
                    source: "auth",
                    level: "error",
                    message: oauthErrorMsg || "[Authorization failed]: no accessCode received from the server"
                });
            }
        } else {
            oauth2.callback({auth: oauth2.auth, token: qp, isValid: isValid, redirectUrl: redirectUrl});
        }
        window.close();
    }
</script>

`,
}
var tDocIndex = tplNode{
	"doc/index.html",
	`
<!-- HTML for static distribution bundle build -->
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="//unpkg.com/swagger-ui-dist@3/swagger-ui.css">
    <link rel="icon" type="image/png" href="//unpkg.com/swagger-ui-dist@3/favicon-32x32.png" sizes="32x32"/>
    <link rel="icon" type="image/png" href="//unpkg.com/swagger-ui-dist@3/favicon-16x16.png" sizes="16x16"/>
    <style>
        html {
            box-sizing: border-box;
            overflow: -moz-scrollbars-vertical;
            overflow-y: scroll;
        }

        *,
        *:before,
        *:after {
            box-sizing: inherit;
        }

        body {
            margin: 0;
            background: #fafafa;
        }
    </style>
</head>

<body>
<div id="swagger-ui"></div>

<script src="//unpkg.com/swagger-ui-dist@3/swagger-ui-bundle.js"></script>
<script src="//unpkg.com/swagger-ui-dist@3/swagger-ui-standalone-preset.js"></script>
<script>
    window.onload = function () {
        // Build a system
        const ui = SwaggerUIBundle({
            url: "swagger.yaml",
            dom_id: '#swagger-ui',
            deepLinking: true,
            presets: [
                SwaggerUIBundle.presets.apis,
                SwaggerUIStandalonePreset
            ],
            plugins: [
                SwaggerUIBundle.plugins.DownloadUrl
            ],
            layout: "StandaloneLayout"
        })
        window.ui = ui
    }
</script>
</body>
</html>

`,
}

//"tpl/config.toml": "config.toml",

var tMod = tplNode{
	NameFormat: "go.mod",
	TplContent: `module {{.AppPkg}}

go 1.12
`,
}

var tModelDbRedis = tplNode{
	NameFormat: "model/db_redis.go",
	TplContent: `
package model

import (
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

var redisDB *redis.Client
const redisPrefix = "ginbro:"

func CreateRedis(redisAddr, redisPassword string, idx int) {
	//initializing redis client
	redisDB = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword, // no password set
		DB:       idx,           // use default DB
	})
	if pong, err := redisDB.Ping().Result(); err != nil || pong != "PONG" {
		logrus.WithError(err).Fatal("could not connect to the redis server")
	}

}

`,
}

var parseOneList = []tplNode{tDocIndex, tDocOauth2, tStaticIndex, tTaskCore, tTaskExample, tMod, tDocYaml, tConfig, tReadme, tGitIgnore, tHandlerGin, tHandlerHelper, tMain, tModelDbMem, tModelDbRedis, tModelDbSql, tModelHelper, tConfigToml}

var parseObjList = []tplNode{tModelObj, tHandlersObj}
