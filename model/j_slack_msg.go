package model

import (
	"database/sql/driver"
	"encoding/json"
)

type SlackMsg struct {
	Text        string `json:"text"`
	Username    string `json:"username"`
	IconURL     string `json:"icon_url"`
	IconEmoji   string `json:"icon_emoji"`
	Channel     string `json:"channel"`
	UnfurlLinks bool   `json:"unfurl_links"`
	Attachments []struct {
		Title    string `json:"title"`
		Fallback string `json:"fallback"`
		Text     string `json:"text"`
		Pretext  string `json:"pretext"`
		Color    string `json:"color"`
		Fields   []struct {
			Title string `json:"title"`
			Value string `json:"value"`
			Short bool   `json:"short"`
		} `json:"fields"`
	} `json:"attachments"`
}

func (o SlackMsg) Value() (driver.Value, error) {
	b, err := json.Marshal(o)
	return string(b), err
}

func (o *SlackMsg) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), o)
}
