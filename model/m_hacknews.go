package model

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type HackNewQ struct {
	PaginationQ
	HackNew
}

func (cq *HackNewQ) SearchAll() (data *PaginationQ, err error) {
	page := cq.PaginationQ
	page.Data = &[]Comment{} //make sure page.Data is not nil and is a slice gorm.Model

	m := cq.HackNew
	tx := db.Model(cq.HackNew).Preload("User")
	//customize search column
	if m.TitleEn != "" {
		tx = tx.Where("`title_en` = ?", m.TitleEn)
	}
	if m.TitleZh != "" {
		tx = tx.Where("`title_zh` = ?", m.TitleZh)
	}
	return page.SearchAll(tx)
}

type HackNew struct {
	BaseModel
	TitleZh string `json:"title_zh" form:"title_zh"`
	TitleEn string `json:"title_en" form:"title_en"`
	Url     string `gorm:"index" json:"url" form:"url"`
	Cate    string `json:"cate" comment:"news or show" form:"cate"`
}

//Delete
func (m HackNew) Delete(ids []uint) (err error) {
	if len(ids) > 0 {
		err = db.Where("`id` in (?)", ids).Delete(m).Error
	}
	return
}

//Update a row
func (m *HackNew) Update() (err error) {
	return db.Model(m).Update(m).Error
}

func (m *HackNew) translateEn2ch() (err error) {
	m.TitleZh, err = TranslateEn2Ch(m.TitleEn)
	return
}
func (m *HackNew) CreateOrUpdate() (err error) {
	_, err = url.Parse(m.Url)
	if err != nil {
		return err
	}
	row := HackNew{}
	if db.Where("url = ?", m.Url).First(&row).RecordNotFound() {
		err = m.translateEn2ch()
		if err != nil {
			logrus.WithError(err).WithField("en", m.TitleEn).Error("HackNew.CreateOrUpdate 翻译失败")
		}
		return db.Create(m).Error
	}
	if row.TitleZh == "" {
		err = m.translateEn2ch()
		if err != nil {
			logrus.WithError(err).WithField("en", m.TitleEn).Error("HackNew.CreateOrUpdate 翻译失败2")
		}
	}
	return db.Model(row).Where("url = ?", m.Url).Updates(*m).Error
}

func (m HackNew) TodayRowBy(cate string) (list []HackNew, err error) {
	list = []HackNew{}
	today := time.Now().Format("2006-01-02")
	err = db.Model(m).Where("date(`updated_at`) = ? AND `cate` = ?", today, cate).Find(&list).Error
	return
}

/*
	有道智云
	文档页面  http://ai.youdao.com/docs/doc-trans-api.s#p01
*/

func TranslateCh2En(text string) (string, error) {
	res, err := youdaoTranslateApi(text, "zh-CHS", "EN")
	if err != nil {
		return "", err
	}
	if len(res.Translation) > 0 {
		return res.Translation[0], nil
	}
	return "", nil

}

func TranslateEn2Ch(text string) (string, error) {
	res, err := youdaoTranslateApi(text, "EN", "zh-CHS")
	if err != nil {
		return "", err
	}
	if len(res.Translation) > 0 {
		return res.Translation[0], nil
	}
	return "", nil
}

func youdaoTranslateApi(q, from, to string) (obj *respObj, err error) {
	salt := fmt.Sprintf("%d", rand.Intn(99999))
	appKey := viper.GetString("spiderhn.youdaoAppKey")
	appSecret := viper.GetString("spiderhn.youdaoAppSecret")
	appHost := viper.GetString("spiderhn.youdaoAppHost")
	ts := fmt.Sprintf("%d", time.Now().Unix())
	sign := generateSign(appKey, q, salt, ts, appSecret)
	data := url.Values{
		"q":        {q},
		"to":       {to},
		"from":     {from},
		"appKey":   {appKey},
		"salt":     {salt},
		"sign":     {sign},
		"curtime":  {ts},
		"signType": {"v3"},
	}
	resp, err := http.PostForm(appHost, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	obj = &respObj{}
	err = json.NewDecoder(resp.Body).Decode(obj)
	if err != nil {
		return nil, err
	}
	if obj.HasError() {
		return nil, fmt.Errorf("错误码:%s", obj.ErrorCode)
	}
	time.Sleep(time.Millisecond * 100)
	return obj, nil
}

type respObj struct {
	TSpeakURL    string   `json:"tSpeakUrl"`
	ReturnPhrase []string `json:"returnPhrase"`
	Query        string   `json:"query"`
	Translation  []string `json:"translation"`
	ErrorCode    string   `json:"errorCode"`
	L            string   `json:"l"`
}

func (o respObj) HasError() bool {
	code, err := strconv.ParseInt(o.ErrorCode, 10, 32)
	if err != nil {
		logrus.WithError(err).Error("HasError parse int error")
		return true
	}
	if code > 0 {
		return true
	}
	return false
}

func generateSign(appKey, q, salt, curTime, appSecret string) string {
	input := ""
	a := []rune(q)
	inputLen := len(a)
	if inputLen < 20 {
		input = string(a)
	} else {
		input = fmt.Sprintf("%s%d%s", string(a[:10]), inputLen, string(a[inputLen-10:]))
	}
	temp := appKey + input + salt + curTime + appSecret
	sum := sha256.Sum256([]byte(temp))
	return fmt.Sprintf("%x", sum)
}
