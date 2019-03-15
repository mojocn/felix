package ginbro

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/sirupsen/logrus"
	"strings"
)

type ColumnSchema struct {
	TableName      string
	ColumnName     string
	ColumnType     string
	ColumnDataType string
	ColumnComment  string
	ColumnKey      string
}

func (c *ColumnSchema) toProperty(authTable, passwordColumn string) Property {
	//very import the g
	p := Property{
		ColumnName:    c.ColumnName,
		ColumnComment: strings.Replace(c.ColumnComment, `"`, "", -1),
		ModelProp:     strcase.ToCamel(strings.ToLower(c.ColumnName)),
	}
	if c.ColumnName == passwordColumn && c.TableName == authTable {
		p.IsAuthColumn = true
	}
	modelType, swgType, swgFormat := transType(c.ColumnDataType, c.ColumnName, c.ColumnType, c.ColumnKey)
	p.ModelType = modelType
	p.SwaggerType = swgType
	p.SwaggerFormat = swgFormat
	p.ModelTag = fmt.Sprintf(`gorm:"column:%s" form:"%s" json:"%s" comment:"%s" columnType:"%s" dataType:"%s" columnKey:"%s"`, c.ColumnName, c.ColumnName, c.ColumnName, c.ColumnComment, c.ColumnType, c.ColumnDataType, c.ColumnKey)
	return p
}
func transType1(dataType, columnName, columnType, columnKey string) (string, string, string) {
	modelType, swgType, swgFormat := "NoneType", "", ""
	switch dataType {
	case "varchar", "longtext", "char", "enum", "set", "mediumtext", "json", "text", "tinytext":
		modelType = "string"
		swgFormat, swgType = "string", "string"
	case "bigint":
		modelType, swgFormat, swgType = "int", "int64", "integer"
		if strings.Contains(columnType, "unsigned") {
			modelType, swgFormat, swgType = "uint64", "int64", "integer"
		}
	case "int", "tinyint", "smallint", "mediumint":
		modelType, swgFormat, swgType = "int", "int64", "integer"
		if strings.Contains(columnType, "unsigned") {
			modelType, swgFormat, swgType = "uint", "int32", "integer"
		}
	case "decimal", "float":
		modelType, swgFormat, swgType = "float32", "float", "number"
	case "double":
		modelType, swgFormat, swgType = "float64", "float", "number"
	case "blob":
		modelType, swgFormat, swgType = "*[]byte", "binary", "string"
	case "time", "datetime", "timestamp":
		modelType, swgFormat, swgType = "*time.Time", "date-time", "string"
	}
	if (columnName == "ID" || columnName == "Id" || columnName == "iD" || columnName == "id") && columnKey == "PRI" {
		modelType, swgFormat, swgType = "uint", "int64", "integer"
	}
	return modelType, swgType, swgFormat
}
func transType(dataType, columnName, columnType, columnKey string) (string, string, string) {
	modelType, swgType, swgFormat := "", "", ""
	switch dataType {
	case "varchar", "longtext", "char", "enum", "set", "mediumtext", "json", "text", "tinytext", "date", "year":
		modelType, swgFormat, swgType = "string", "string", "string"
	case "bigint":
		modelType, swgFormat, swgType = "int", "int64", "integer"
		if strings.Contains(columnType, "unsigned") {
			modelType, swgFormat, swgType = "uint64", "int64", "integer"
		}
	case "int", "tinyint", "smallint", "mediumint":
		modelType, swgFormat, swgType = "int", "int64", "integer"
		if strings.Contains(columnType, "unsigned") {
			modelType, swgFormat, swgType = "uint", "int32", "integer"
		}
	case "decimal", "float":
		modelType, swgFormat, swgType = "float32", "float", "number"
	case "double":
		modelType, swgFormat, swgType = "float64", "float", "number"
	case "blob":
		modelType, swgFormat, swgType = "*[]byte", "binary", "string"
	case "time", "datetime", "timestamp":
		modelType, swgFormat, swgType = "*time.Time", "date-time", "string"
	default:
		//logrus.WithField("dataType",dataType).WithField("columnType",columnType).Info("PG")

	}
	if modelType == "" {
		switch columnType {
		case "varchar":
			modelType, swgFormat, swgType = "string", "string", "string"
		case "int4":
			modelType, swgFormat, swgType = "uint", "int32", "integer"
		case "time", "datetime", "timestamp":
			modelType, swgFormat, swgType = "*time.Time", "date-time", "string"
		case "bool":
			modelType, swgFormat, swgType = "bool", "boolean", "boolean"
		default:
			logrus.WithField("dataType", dataType).WithField("columnType", columnType).Info("PG")
			modelType, swgFormat, swgType = "NoneType", "NoneType", "NoneType"
		}

	}

	if (columnName == "ID" || columnName == "Id" || columnName == "iD" || columnName == "id") && columnKey == "PRI" {
		modelType, swgFormat, swgType = "uint", "int64", "integer"
	}
	return modelType, swgType, swgFormat
}
