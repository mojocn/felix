package ginbro

import (
	"database/sql"
	"github.com/sirupsen/logrus"
)

func fetchSchemaMysql(addr, user, password, database, charset string) ([]ColumnSchema, error) {
	db, err := newDb("mysql", addr, user, password, database, charset)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var list []ColumnSchema
	rawSql := "SELECT `TABLE_NAME`, `COLUMN_NAME`,`DATA_TYPE`,`COLUMN_TYPE`,`COLUMN_COMMENT`,`COLUMN_KEY` FROM `INFORMATION_SCHEMA`.`COLUMNS` WHERE `TABLE_SCHEMA` = ?"
	rows, err := db.Query(rawSql, database)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var tableName, columnName, dataType, columnType, columnComment, columnKey sql.NullString
		if rows.Scan(&tableName, &columnName, &dataType, &columnType, &columnComment, &columnKey) == nil {
			c := ColumnSchema{tableName.String, columnName.String, columnType.String, dataType.String, columnComment.String, columnKey.String}
			list = append(list, c)
		} else {
			//create model and handler for every tableName
			logrus.Errorf("get %s database column info failed", database)
		}
	}

	return list, err
}
