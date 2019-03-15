package ginbro

import (
	"database/sql"
	"fmt"
)

func fetchSchemaSQLite(addr, user, password, database, charset string) ([]ColumnSchema, error) {
	db, err := newDb("sqlite3", addr, user, password, database, charset)
	defer db.Close()
	if err != nil {
		return nil, err
	}
	var list []ColumnSchema
	//TODO::找不到table info BUG
	rawSql := `SELECT tbl_name FROM sqlite_master WHERE type = 'table'`
	rows, err := db.Query(rawSql, database)
	defer rows.Close()

	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var tableName sql.NullString
		if err = rows.Scan(&tableName); err != nil {
			return nil, fmt.Errorf("get %s slqite tables failed, [ERROR] %s", database, err)
		} else {
			cols, err := db.Query("PRAGMA table_info('?')", tableName.String)
			if err != nil {
				return nil, err
			}
			for cols.Next() {
				var cCid, columnName, columnType, cNotnull, cDfltValue, cPk sql.NullString
				if err = rows.Scan(&cCid, &columnName, &columnType, &cNotnull, &cDfltValue, &cPk); err != nil {
					return nil, fmt.Errorf("[Error]%s [table]%s", err, tableName.String)
				} else {
					node := ColumnSchema{tableName.String,
						columnName.String,
						columnType.String,
						"",
						"",
						cPk.String}
					if cPk.String == "1" {
						node.ColumnKey = "PRI"
					}
					list = append(list, node)
				}
			}
			cols.Close()
		}
	}
	return list, nil
}
