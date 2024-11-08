package ginbro

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mojocn/felix/model"
	"github.com/mojocn/felix/util"
	"github.com/sirupsen/logrus"
)

type GinbroApp struct {
	model.Ginbro
	Resources []Resource `json:"-"`
	Files     []string
}

func Run(gc model.Ginbro) (*GinbroApp, error) {
	gc.IsSuccess = false
	if gc.AppPkg == "" {
		return nil, errors.New("app package name can't be empty string")
	}
	app, err := newGinbroApp(gc)
	if err != nil {
		return nil, err
	}
	err = app.generateCodeBase()
	if err != nil {
		return nil, err
	}
	//go fmt codebase
	//https://cloud.tencent.com/developer/article/1417112
	err = app.goFmtCodeBase()
	if err != nil {
		logrus.WithError(err).Error("go fmt code base failed")
	}
	app.IsSuccess = true
	//todo: remove ignore go fmt error
	return app, err
}

//generate model files from db
func RunModels(gc model.Ginbro) (*GinbroApp, error) {
	gc.IsSuccess = false
	if gc.AppPkg == "" {
		return nil, errors.New("app package name can't be empty string")
	}
	app, err := newGinbroApp(gc)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func newGinbroApp(gb model.Ginbro) (*GinbroApp, error) {
	cols, err := FetchDbColumn(gb)
	if err != nil {
		return nil, err
	}
	resources, err := transformToResources(cols, gb.AuthTable, gb.AuthColumn)
	if err != nil {
		return nil, err
	}
	if len(gb.AppSecret) < 32 {
		gb.AppSecret = util.RandomString(32)
	}
	return &GinbroApp{
		Ginbro:    gb,
		Resources: resources,
	}, nil
}
func (app *GinbroApp) generateCodeBase() error {

	for _, tplNode := range parseOneList {
		err := tplNode.ParseExecute(app.AppDir, "", app)
		if err != nil {
			return fmt.Errorf("parse [%s] template failed with error : %s", tplNode.NameFormat, err)
		}
	}

	for _, resource := range app.Resources {
		resource.AppPkg = app.AppPkg
		tableName := resource.TableName
		//generate model from resource
		for _, tplNode := range parseObjList {
			err := tplNode.ParseExecute(app.AppDir, tableName, resource)
			if err != nil {
				return fmt.Errorf("parse [%s] template failed with error : %s", tplNode.NameFormat, err)
			}
		}
		if resource.IsAuthTable {
			err := tHandlerLogin.ParseExecute(app.AppDir, "", resource)
			if err != nil {
				return fmt.Errorf("parse [%s] template failed with error : %s", tModelJwt.NameFormat, err)
			}
			err = tHandlerMiddlewareJwt.ParseExecute(app.AppDir, "", resource)
			if err != nil {
				return fmt.Errorf("parse [%s] template failed with error : %s", tHandlerMiddlewareJwt.NameFormat, err)
			}
			err = tModelJwt.ParseExecute(app.AppDir, "", resource)
			if err != nil {
				return fmt.Errorf("parse [%s] template failed with error : %s", tModelJwt.NameFormat, err)
			}
		}
	}
	return nil
}

func (app *GinbroApp) goFmtCodeBase() error {
	cmd := exec.Command("go", "fmt", "./...")
	cmd.Dir = app.AppDir
	cmd.Env = append(os.Environ(), "GOPROXY=https://goproxy.io")
	bb, err := cmd.CombinedOutput()
	if err != nil {
		//print gin-goinc/autols failure
		// fix it :::  https://github.com/gin-gonic/gin/issues/1673
		return fmt.Errorf("%s   %s", string(bb), err)
	}
	return nil
}
func (app *GinbroApp) ListAppFileTree() error {
	return filepath.Walk(app.AppDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fmt.Println(path)
		if !info.IsDir() {
			app.Files = append(app.Files, path)
		}
		return nil
	})
}
