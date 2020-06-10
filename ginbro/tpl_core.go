package ginbro

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const backQuote = "_[BACKQUOTE]_"

type tplNode struct {
	NameFormat string
	TplContent string
}

func (n *tplNode) ParseExecute(appDir, pathArg string, data interface{}) error {
	var p string
	if pathArg != "" {
		p = fmt.Sprintf(n.NameFormat, pathArg)
	} else {
		p = n.NameFormat
	}
	p = filepath.Clean(filepath.Join(appDir, p))
	err := os.MkdirAll(filepath.Dir(p), 0644)
	if err != nil {
		return err
	}
	tplFormat := strings.Replace(n.TplContent, backQuote, "`", -1)
	tmpl, err := template.New(p).Parse(tplFormat)
	file, err := os.Create(p)
	if err != nil {
		return err
	}
	defer file.Close()
	return tmpl.Execute(file, data)
}
