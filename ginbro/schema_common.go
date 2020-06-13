package ginbro

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/libragen/felix/model"
)

func newDb(dbType, addr, user, password, database, charset string) (*sql.DB, error) {
	hostPort := strings.Split(addr, ":")
	if len(hostPort) != 2 && dbType != "sqlite" {
		return nil, fmt.Errorf("flag addr [%s] is a wrong format string, must be like host:port", addr)
	}
	host := hostPort[0]
	port := hostPort[1]

	var dbConn string
	switch dbType {
	case "mysql":
		dbConn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s", user, password, host, port, database, charset)
	case "postgres":
		dbConn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", host, port, user, database, password)
	case "mssql":
		dbConn = fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", user, password, host, port, database)
	case "sqlite3":
		dbConn = addr
	default:
		return nil, fmt.Errorf("felix rest doesn't support %s database", dbType)
	}
	return sql.Open(dbType, dbConn)
}
func FetchDbColumn(gb model.Ginbro) ([]ColumnSchema, error) {
	switch gb.DbType {
	case "mysql":
		return fetchSchemaMysql(gb.DbAddr, gb.DbUser, gb.DbPassword, gb.DbName, gb.DbChar)
	case "postgres":
		return fetchSchemaPg(gb.DbAddr, gb.DbUser, gb.DbPassword, gb.DbName, gb.DbChar)
	case "mssql":
		//TODO:: mssql
		return nil, fmt.Errorf("to do support and test mssql database %s", "!!!sorry!!!")
		//returdbConn = fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", user, password, host, port, database)
	case "sqlite", "sqlite3":

		return fetchSchemaSQLite(gb.DbAddr, gb.DbUser, gb.DbPassword, gb.DbName, gb.DbChar)
	default:
		return nil, fmt.Errorf("felix rest doesn't support %s type DB", gb.DbType)
	}
}
