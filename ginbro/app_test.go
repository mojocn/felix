package ginbro

import (
	"os"
	"testing"

	"github.com/libragen/felix/model"
)

func TestRun(t *testing.T) {

	gc := model.Ginbro{
		AppSecret:  "sdfsadfewdddcd",
		AppDir:     "/TEST0",
		AppAddr:    "127.0.0.1:4444",
		AppPkg:     "ginBRO",
		AuthTable:  "users",
		AuthColumn: "password",
		DbUser:     "venom",
		DbPassword: os.Getenv("TESTPASSWORD"),
		DbAddr:     os.Getenv("TESTDBADDR"),
		DbType:     "mysql",
		DbName:     "venom",
		DbChar:     "utf8",
	}

	_, err := Run(gc)
	if err != nil {
		t.Fatal(err)
	}

}
