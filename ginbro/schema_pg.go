package ginbro

import (
	"database/sql"
	"github.com/sirupsen/logrus"
)

func fetchSchemaPg(addr, user, password, database, charset string) ([]ColumnSchema, error) {
	db, err := newDb("postgres", addr, user, password, database, charset)
	defer db.Close()
	if err != nil {
		return nil, err
	}
	var list []ColumnSchema
	//https://stackoverflow.com/questions/343138/retrieving-comments-from-a-postgresql-db
	//SELECT c.table_schema,c.table_name,c.column_name,pgd.description
	//FROM pg_catalog.pg_statio_all_tables as st
	//  inner join pg_catalog.pg_description pgd on (pgd.objoid=st.relid)
	//  inner join information_schema.columns c on (pgd.objsubid=c.ordinal_position
	//    and  c.table_schema=st.schemaname and c.table_name=st.relname);
	//rawSql := "SELECT table_name, column_name, data_type, udt_name FROM information_schema.columns WHERE table_catalog = $1 AND table_schema = 'public'"
	rawSql := `SELECT c.table_name, c.column_name, c.data_type, c.udt_name, pgd.description
FROM information_schema.columns c
LEFT JOIN pg_catalog.pg_statio_all_tables AS st ON ( c.table_schema=st.schemaname AND c.table_name=st.relname)
LEFT JOIN pg_catalog.pg_description pgd ON (pgd.objoid=st.relid AND pgd.objsubid=c.ordinal_position)
WHERE table_catalog = $1 AND table_schema = 'public'`
	rows, err := db.Query(rawSql, database)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var tableName, columnName, columnType0, columnType, columnComment sql.NullString
		if err = rows.Scan(&tableName, &columnName, &columnType0, &columnType, &columnComment); err != nil {
			logrus.WithError(err).Errorf("get %s database column info failed", database)
		} else {
			c := ColumnSchema{
				tableName.String,
				columnName.String,
				columnType.String,
				columnType0.String,
				columnComment.String,
				""}
			list = append(list, c)
		}
	}
	return list, err
}
