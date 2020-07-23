package db

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr/v2"
	"github.com/gocraft/dbr/v2/dialect"
	"github.com/pkg/errors"
)

type indexInfo struct {
	Name    string
	Type    string
	Columns []string
}

type tagPair struct {
	Key string
	Val string
}

// CreateDatabase 如果数据库不存在，则创建
func CreateDatabase(dbConfig Config) (err error) {
	dbName, dsnWithoutDBName, err := splitDSN(dbConfig.Driver, dbConfig.DSN)
	if err != nil || len(dbName) == 0 || len(dsnWithoutDBName) == 0 {
		return
	}

	conn, err := dbr.Open(dbConfig.Driver, dsnWithoutDBName, nil)
	if err != nil {
		return errors.WithStack(err)
	}
	defer conn.Close()

	sqlstr := "CREATE DATABASE IF NOT EXISTS " + conn.QuoteIdent(dbName)

	if dbConfig.Driver == "mysql" {
		sqlstr += " DEFAULT CHARSET utf8mb4 COLLATE utf8_general_ci"
	}

	_, err = conn.Exec(sqlstr)
	return errors.WithStack(err)
}

// CreateDatabase 如果表不存在，则创建
func CreateTable(session *dbr.Session, table string, schema interface{}) (err error) {
	baseDialect := getBaseDialect(session)
	rows, err := session.Select("*").From(table).Where("1 != 1").Rows()
	if err != nil {
		createSQL := schema2CreateTableSQL(baseDialect, table, schema)
		_, err = session.InsertBySql(createSQL).Exec()
		if err != nil {
			return errors.WithStack(err)
		}
		for _, indexSQL := range listTableIndexSQL(baseDialect, table, schema) {
			if _, err = session.InsertBySql(indexSQL).Exec(); err != nil {
				return errors.WithStack(err)
			}
		}
		return
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return errors.WithStack(err)
	}
	colTypes, _ := rows.ColumnTypes()
	quotedTable := session.QuoteIdent(table)
	tp := reflect.Indirect(reflect.ValueOf(schema)).Type()

	for i := 0; i < tp.NumField(); i++ {
		field := tp.Field(i)
		columnName := dbr.NameMapping(field.Name)

		if field.Tag.Get("db") == "-" {
			continue
		}

		if index := findIndex(columns, columnName); index < 0 {
			sqlstr := fmt.Sprintf("ALTER TABLE %s ADD %s", quotedTable, field2SQL(baseDialect, field))
			if _, err = session.InsertBySql(sqlstr).Exec(); err != nil {
				return errors.WithStack(err)
			}
		} else if scanType := colTypes[i].ScanType(); scanType != nil {
			scanKind := scanType.Kind()
			fieldKind := fieldKind(field)

			if kindType(scanKind) != kindType(fieldKind) {
				// TODO: 暂时不开启修改字段
				// sqlstr := fmt.Sprintf("ALTER TABLE %s MODIFY %s", quotedTable, field2SQL(baseDialect, field))
				// if _, err = tx.InsertBySql(sqlstr).Exec(); err != nil {
				// 	return errors.WithStack(err)
				// }
			}
		}
	}

	return
}

func schema2CreateTableSQL(d dbr.Dialect, table string, schema interface{}) string {
	sqlstr := "CREATE TABLE IF NOT EXISTS " + d.QuoteIdent(table)
	sqlstr += "(\n"

	tp := reflect.Indirect(reflect.ValueOf(schema)).Type()
	columns := []string{}

	for i := 0; i < tp.NumField(); i++ {
		field := tp.Field(i)

		if field.Tag.Get("db") == "-" {
			continue
		}
		columns = append(columns, field2SQL(d, field))
	}

	sqlstr += strings.Join(columns, ",\n") + "\n)"

	return sqlstr
}

func field2SQL(d dbr.Dialect, field reflect.StructField) string {
	buf := strings.Builder{}
	columnName := dbr.NameMapping(field.Name)
	tag := field.Tag
	kind := fieldKind(field)

	buf.WriteString(d.QuoteIdent(columnName))
	buf.WriteString(" ")

	if strings.ToLower(field.Name) == "id" && kindType(kind) == kindType_Number {
		switch d {
		case dialect.MySQL, dialect.PostgreSQL:
			buf.WriteString("SERIAL PRIMARY KEY")
		case dialect.SQLite3:
			buf.WriteString("INTEGER PRIMARY KEY AUTOINCREMENT")
		}
		return buf.String()
	}

	sqlTag := tag.Get("sql")
	pairs := parseTag2Map(sqlTag)
	defaultLen, defaultVal := "", ""

	if typ := pairs["type"]; len(typ) > 0 {
		buf.WriteString(strings.ToUpper(typ))
	} else {
		switch kindType(kind) {
		case kindType_Number:
			if kind == reflect.Float32 || kind == reflect.Float64 {
				buf.WriteString("DECIMAL")
				defaultLen = "(20,2)"
			} else {
				buf.WriteString("INT")
				defaultLen = "(11)"
			}
			defaultVal = "'0'"
		case kindType_String:
			buf.WriteString("VARCHAR")
			defaultLen, defaultVal = "(255)", "''"
		case kindType_Object:
			switch fieldType(field) {
			case reflect.TypeOf(dbr.NullTime{}), reflect.TypeOf(time.Time{}):
				if d == dialect.PostgreSQL {
					buf.WriteString("TIMESTAMP")
				} else {
					buf.WriteString("DATETIME")
				}
			default:
				buf.WriteString("TEXT")
			}
		}
	}

	if length := pairs["length"]; len(length) > 0 {
		buf.WriteString("(" + length + ")")
	} else if len(defaultLen) > 0 {
		buf.WriteString(defaultLen)
	}

	buf.WriteString(" ")

	if field.Type.Kind() == reflect.Ptr {
		buf.WriteString("NULL")
		defaultVal = ""
	} else {
		buf.WriteString("NOT NULL")
	}

	if def, ok := pairs["default"]; ok {
		buf.WriteString(" DEFAULT " + def)
	} else if len(defaultVal) > 0 {
		buf.WriteString(" DEFAULT " + defaultVal)
	}

	return buf.String()
}

func listTableIndexSQL(d dbr.Dialect, table string, schema interface{}) (sqls []string) {
	for _, index := range listIndex(schema) {
		indexType := "INDEX"

		switch strings.ToLower(index.Type) {
		case "unique":
			indexType = "UNIQUE INDEX"
		case "primary":
			indexType = "PRIMARY KEY"
		}

		columns := []string{}

		for _, col := range index.Columns {
			columns = append(columns, d.QuoteIdent(col))
		}

		indexSQL := fmt.Sprintf("CREATE %s %s ON %s (%s)",
			indexType, d.QuoteIdent(index.Name), d.QuoteIdent(table), strings.Join(columns, ","))
		sqls = append(sqls, indexSQL)
	}

	return
}

func getBaseDialect(session *dbr.Session) dbr.Dialect {
	baseDialect := session.Dialect

	if ld, ok := baseDialect.(*localDialect); ok {
		baseDialect = ld.Dialect
	}

	return baseDialect
}

func splitDSN(driver, dsn string) (dbName, dsnWithoutDBName string, err error) {
	switch driver {
	case "mysql":
		cfg, err1 := mysql.ParseDSN(dsn)
		if err1 != nil {
			err = errors.WithStack(err1)
			return
		}

		dbName = cfg.DBName
		cfg.DBName = ""
		dsnWithoutDBName = cfg.FormatDSN()

	case "postgres":
		dsnURL, err1 := url.Parse(dsn)
		if err1 != nil {
			err = errors.WithStack(err1)
			return
		}

		if len(dsnURL.Path) > 0 {
			dbName = strings.Trim(dsnURL.Path, "/")
			dsnURL.Path = "/"
			dsnWithoutDBName = dsnURL.String()
		}
	}

	return
}

func parseTag2Slice(sqlTag string) (pairs []tagPair) {
	if len(sqlTag) > 0 {
		for _, tagField := range strings.Fields(sqlTag) {
			parts := strings.SplitN(tagField, ":", 2)
			if len(parts) != 2 {
				panic(fmt.Sprintf("unknonw tag: %s", sqlTag))
			}

			pairs = append(pairs, tagPair{
				Key: strings.TrimSpace(parts[0]),
				Val: strings.TrimSpace(parts[1]),
			})
		}
	}

	return
}

func parseTag2Map(sqlTag string) map[string]string {
	pairs := make(map[string]string)

	for _, pair := range parseTag2Slice(sqlTag) {
		pairs[pair.Key] = pair.Val
	}

	return pairs
}

func listIndex(schema interface{}) (pairs map[string]*indexInfo) {
	tp := reflect.Indirect(reflect.ValueOf(schema)).Type()
	pairs = make(map[string]*indexInfo)

	for i := 0; i < tp.NumField(); i++ {
		field := tp.Field(i)
		slice := parseTag2Slice(field.Tag.Get("sql"))
		column := dbr.NameMapping(field.Name)

		for _, pair := range slice {
			if pair.Key == "index" {
				parts := strings.SplitN(pair.Val, ",", 2)
				indexName := parts[0]
				indexType := ""
				if len(parts) > 1 {
					indexType = parts[1]
				}
				if index, ok := pairs[indexName]; ok {
					index.Columns = append(index.Columns, column)
				} else {
					pairs[indexName] = &indexInfo{
						Name:    indexName,
						Type:    indexType,
						Columns: []string{column},
					}
				}
			}
		}
	}
	return
}

const (
	kindType_Number = "number"
	kindType_String = "string"
	kindType_Object = "object"
)

func kindType(kind reflect.Kind) string {
	if kind == reflect.Bool ||
		kind == reflect.Int ||
		kind == reflect.Int8 ||
		kind == reflect.Int16 ||
		kind == reflect.Int32 ||
		kind == reflect.Int64 ||
		kind == reflect.Uint ||
		kind == reflect.Uint8 ||
		kind == reflect.Uint16 ||
		kind == reflect.Uint32 ||
		kind == reflect.Uint64 ||
		kind == reflect.Float32 ||
		kind == reflect.Float64 {
		return kindType_Number
	} else if kind == reflect.String {
		return kindType_String
	} else {
		return kindType_Object
	}
}

func fieldType(field reflect.StructField) reflect.Type {
	fk := field.Type.Kind()

	if fk == reflect.Ptr {
		return field.Type.Elem()
	}

	return field.Type
}

func fieldKind(field reflect.StructField) reflect.Kind {
	return fieldType(field).Kind()
}

func findIndex(list []string, e string) int {
	for i, s := range list {
		if s == e {
			return i
		}
	}
	return -1
}
