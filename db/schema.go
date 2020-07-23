package db

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr/v2"
	"github.com/pkg/errors"
)

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

	_, err = conn.Exec("CREATE DATABASE IF NOT EXISTS " + conn.QuoteIdent(dbName))
	if err != nil {
		return errors.WithStack(err)
	}

	return
}

// CreateDatabase 如果表不存在，则创建
func CreateTable(session *dbr.Session, table string, schema interface{}) (err error) {
	rows, err := session.Select("*").From(table).Where("1 != 1").Rows()
	if err != nil {
		createSQL := schema2CreateTableSQL(session, table, schema)
		_, err = session.InsertBySql(createSQL).Exec()

		return errors.WithStack(err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return errors.WithStack(err)
	}
	colTypes, _ := rows.ColumnTypes()

	quotedTable := session.QuoteIdent(table)
	tp := reflect.TypeOf(reflect.Indirect(reflect.ValueOf(schema)))

	for i := 0; i < tp.NumField(); i++ {
		field := tp.Field(i)
		columnName := dbr.NameMapping(field.Name)

		if field.Tag.Get("db") == "-" {
			continue
		}

		if index := findIndex(columns, columnName); index < 0 {
			sqlstr := fmt.Sprintf("ALTER TABLE %s ADD %s", quotedTable, field2SQL(session, field))
			if _, err = session.InsertBySql(sqlstr).Exec(); err != nil {
				return errors.WithStack(err)
			}
		} else if scanType := colTypes[i].ScanType(); scanType != nil {
			scanKind := scanType.Kind()
			fieldKind := fieldKind(field)

			if kindType(scanKind) != kindType(fieldKind) {
				sqlstr := fmt.Sprintf("ALTER TABLE %s MODIFY %s", quotedTable, field2SQL(session, field))
				if _, err = session.InsertBySql(sqlstr).Exec(); err != nil {
					return errors.WithStack(err)
				}
			}
		}
	}

	return
}

func schema2CreateTableSQL(d dbr.Dialect, table string, schema interface{}) string {
	sqlstr := "CREATE TABLE IF NOT EXISTS " + d.QuoteIdent(table)
	sqlstr += "(\n"

	tp := reflect.TypeOf(reflect.Indirect(reflect.ValueOf(schema)))
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

	defaultLen := ""

	if typ := tag.Get("type"); len(typ) > 0 {
		buf.WriteString(typ)
	} else {
		switch kindType(kind) {
		case kindType_Number:
			if kind == reflect.Float32 || kind == reflect.Float64 {
				buf.WriteString("DECIMAL")
				defaultLen = "(20,2)"
			} else {
				buf.WriteString("INT")
			}
		case kindType_String:
			buf.WriteString("VARCHAR")
			defaultLen = "(255)"
		}
	}

	if length := tag.Get("length"); len(length) > 0 {
		buf.WriteString("(" + length + ")")
	} else {
		buf.WriteString(defaultLen)
	}

	buf.WriteString(" ")

	if field.Type.Kind() == reflect.Ptr {
		buf.WriteString("NULL")
	} else {
		buf.WriteString("NOT NULL")
	}

	return buf.String()
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

const (
	kindType_Number = "number"
	kindType_String = "string"
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
	} else {
		return kindType_String
	}
}

func fieldKind(field reflect.StructField) reflect.Kind {
	fk := field.Type.Kind()

	if fk == reflect.Ptr {
		fk = field.Type.Elem().Kind()
	}

	return fk
}

func findIndex(list []string, e string) int {
	for i, s := range list {
		if s == e {
			return i
		}
	}
	return -1
}
