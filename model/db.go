package model

import (
	"log"
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mitchellh/go-homedir"
)

var db *gorm.DB
var dbPath string

func init() {
	rand.Seed(time.Now().Unix())
	dir, err := homedir.Dir()
	if err != nil {
		log.Fatal("get home dir failed:", err)
	}
	dbPath = path.Join(dir, ".felix/sqlite.db")
}

func CreateSQLiteDb(verbose bool) {
	//log.Println("SQLite3 in:", dbPath)
	//sqlite, err := gorm.Open("sqlite3", dbPath)
	//if err != nil {
	//	logrus.WithError(err).Fatalf("master fail to open its sqlite db in %s. please install master first.", dbPath)
	//	return
	//}
	//
	//db = sqlite
	////TODO::optimize
	////db.DropTable("term_logs")
	//db.AutoMigrate(Machine{}, Task{}, User{}, Ginbro{}, SshLog{}, WslogHook{}, WslogMsg{})
	//db.LogMode(verbose)
}

func FlushSqliteDb() error {
	db.Close()
	return os.RemoveAll(dbPath)
}
