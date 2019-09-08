package spiderhn

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/dejavuzhou/felix/model"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const backQuote = "_[BACKQUOTE]_"

const jekyllMarkdownTemplate = `---
layout: post
title: Hacknews{{.Day}}新闻
category: Hacknews
tags: hacknews
keywords: hacknews
coverage: hacknews-banner.jpg
---

Hacker News 是一家关于计算机黑客和创业公司的社会化新闻网站，由保罗·格雷厄姆的创业孵化器 Y Combinator 创建。
与其它社会化新闻网站不同的是 Hacker News 没有踩或反对一条提交新闻的选项（不过评论还是可以被有足够 Karma 的用户投反对票）；只可以赞或是完全不投票。简而言之，Hacker News 允许提交任何可以被理解为“任何满足人们求知欲”的新闻。

## HackNews Hack新闻

{{range .News}}
<li><a href="{{.Url}}" rel="nofollow noreferrer">{{.TitleEn}}</a></li>
- _[BACKQUOTE]_{{.TitleZh}}_[BACKQUOTE]_{{end}}


## HackShows Hacks展示
{{range .Shows}}
<li><a href="{{.Url}}" rel="nofollow noreferrer">{{.TitleEn}}</a></li>
- _[BACKQUOTE]_{{.TitleZh}}_[BACKQUOTE]_{{end}}

`
const hackNewsUrl = "https://news.ycombinator.com/news"

func downloadHtml(url string) (*goquery.Document, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("cookie", viper.GetString("spiderhn.cookie"))
	req.Header.Set("User-Agent", viper.GetString("spiderhn.userAgent"))
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, errors.New("the get request's response code is not 200")
	}
	defer res.Body.Close()
	return goquery.NewDocumentFromReader(res.Body)
}
func SpiderHackNews() error {
	doc, err := downloadHtml(hackNewsUrl)
	if err != nil {
		return err
	}
	doc.Find("a.storylink").Each(func(i int, s *goquery.Selection) {
		//todo 相对地址解析
		url, _ := s.Attr("href")
		titleEn := s.Text()
		titleEn = strings.ReplaceAll(titleEn, "[", "")
		titleEn = strings.ReplaceAll(titleEn, "]", "")
		newsItem := model.HackNew{TitleEn: titleEn, Url: url, Cate: "hacknews"}
		err = newsItem.CreateOrUpdate()
		if err != nil {
			logrus.WithError(err).Error("goquery Each save news to db failed")
		}
	})
	return nil
}

func ParsemarkdownHacknews() error {
	mdl := model.HackNew{}
	newsItems, err := mdl.TodayRowBy("hacknews")
	if err != nil {
		return err
	}
	showItems, err := mdl.TodayRowBy("hackshows")
	if err != nil {
		return err
	}
	if len(newsItems) < 1 {
		return nil
	}

	techMojotvCnSrcDir := viper.GetString("tech_mojotv_cn.srcDir")

	tplString := strings.Replace(jekyllMarkdownTemplate, backQuote, "`", -1)
	tmpl, err := template.New("awesome").Parse(tplString)
	if err != nil {
		return err
	}
	day := time.Now().Format("2006-01-02")
	mdFile := fmt.Sprintf("%s-hacknews.md", day)
	mdFilePath := filepath.Join(techMojotvCnSrcDir, "_posts", "hacknews", mdFile)
	mdFile = filepath.Clean(mdFilePath)
	file, err := os.Create(mdFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	data := struct {
		Day   string
		News  []model.HackNew
		Shows []model.HackNew
	}{day, newsItems, showItems}
	err = tmpl.Execute(file, data) //执行模板的merger操作
	return err
}

func SpiderHackShows() error {
	doc, err := downloadHtml("https://news.ycombinator.com/show")
	if err != nil {
		return err
	}
	doc.Find("a.storylink").Each(func(i int, s *goquery.Selection) {
		url, _ := s.Attr("href")
		if !strings.Contains(url, "http") {
			url = "https://news.ycombinator.com/" + url
		}
		titleEn := s.Text()
		titleEn = strings.ReplaceAll(titleEn, "[", "")
		titleEn = strings.ReplaceAll(titleEn, "]", "")
		titleEn = strings.Replace(titleEn, "Show HN:", "", -1)

		newsItem := model.HackNew{TitleEn: titleEn, Url: url, Cate: "hackshows"}
		err = newsItem.CreateOrUpdate()
		if err != nil {
			logrus.WithError(err).Error("goquery Each save shows to db failed")
		}
	})

	return nil
}
