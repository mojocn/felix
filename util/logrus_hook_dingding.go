package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/sirupsen/logrus"
)

var allLvls = []logrus.Level{
	logrus.DebugLevel,
	logrus.InfoLevel,
	logrus.WarnLevel,
	logrus.ErrorLevel,
	logrus.FatalLevel,
	logrus.PanicLevel,
}

func NewDingHook(url, app string, thresholdLevel logrus.Level) *dingHook {
	temp := []logrus.Level{}
	for _, v := range allLvls {
		if v <= thresholdLevel {
			temp = append(temp, v)
		}
	}
	hook := &dingHook{apiUrl: url, levels: temp, appName: app}
	hook.jsonBodies = make(chan []byte)
	hook.closeChan = make(chan bool)
	//开启chan 队列 执行post dingding hook api
	go hook.startDingHookQueueJob()
	//todo 完美ding logrus ding ding hook
	return hook
}

func (dh *dingHook) startDingHookQueueJob() {
	for {
		select {
		case <-dh.closeChan:
			return
		case bs := <-dh.jsonBodies:
			res, err := http.Post(dh.apiUrl, "application/json", bytes.NewBuffer(bs))
			if err != nil {
				log.Println(err)
			}
			if res != nil && res.StatusCode != 200 {
				log.Println("dingHook go chan http post error", res.StatusCode)
			}
		}
	}

}

type dingHook struct {
	apiUrl     string
	levels     []logrus.Level
	appName    string
	jsonBodies chan []byte
	closeChan  chan bool
}

// Levels sets which levels to sent to slack
func (dh *dingHook) Levels() []logrus.Level {
	return dh.levels
}

//Fire2 这个异步有可能导致 最后一条消息丢失,main goroutine 提前结束到导致 子线程http post 没有发送
func (dh *dingHook) Fire2(e *logrus.Entry) error {
	msg, err := e.String()
	if err != nil {
		return err
	}
	dm := dingMsg{Msgtype: "text"}
	dm.Text.Content = fmt.Sprintf("%s \n %s", dh.appName, msg)
	bs, err := json.Marshal(dm)
	if err != nil {
		return err
	}
	dh.jsonBodies <- bs
	return nil
}
func (dh *dingHook) Fire(e *logrus.Entry) error {
	markdownstring := ``
	markdownstring += fmt.Sprintf("## %s %s %s\n", dh.appName, e.Level.String(), e.Time.Format("06-01-02T15:04"))
	markdownstring += fmt.Sprintf("### %s\n", e.Message)
	for k, v := range e.Data {
		switch v := v.(type) {
		case error:
			markdownstring += fmt.Sprintf("> %s \n", v)
		default:
			markdownstring += fmt.Sprintf("## ***%s*** %s\n", k, v)
		}
	}
	if e.HasCaller() {
		funcVal := e.Caller.Function
		fileVal := fmt.Sprintf("%s:%d", e.Caller.File, e.Caller.Line)
		if funcVal != "" {
			markdownstring += fmt.Sprintf("> %s \n", funcVal)
		}
		if fileVal != "" {
			markdownstring += fmt.Sprintf("> %s \n", fileVal)
		}
	}
	dm := dingMsg{Msgtype: "markdown"}
	dm.Markdown.Title = dh.appName
	dm.Markdown.Text = markdownstring

	bs, err := json.Marshal(dm)
	if err != nil {
		return err
	}
	res, err := http.Post(dh.apiUrl, "application/json", bytes.NewBuffer(bs))
	if err != nil {
		return err
	}
	if res != nil && res.StatusCode != 200 {
		return fmt.Errorf("dingHook go chan http post error %d", res.StatusCode)
	}
	return nil
}

type dingMsg struct {
	Msgtype string `json:"msgtype"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
	Markdown struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	} `json:"markdown"`
}
