package handler

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"reflect"

	config "db"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var ctx context.Context

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func TrimSuffixer(s, suffix string) string {
	s = s[:len(s)]
	return s
}

func fetcher(rows *sql.Rows, cols []string, values []sql.RawBytes, scanArgs []interface{}) []map[string]string {

	mapper := make(map[string]string)
	list := make([]map[string]string, len(cols))

	for rows.Next() {
		err := rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
		}

		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			mapper[cols[i]] = value
			if i+1 == len(cols) {
				list = append(list, mapper)
			}
			fmt.Println(cols[i], ": ", value)
		}

		fmt.Println("----------------------")
	}

	if err := rows.Err(); err != nil {
		check(err)
	}

	return list
}

func Fetch(dbc *config.Dbconfig, query string) []map[string]string {

	db := config.Connect(dbc)

	//query := "SELECT * FROM " + table

	rows, err := db.Query(query) //set query with value 'select * from table'
	check(err)

	columns, err := rows.Columns()
	check(err)

	values := make([]sql.RawBytes, len(columns))

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	list := fetcher(rows, columns, values, scanArgs)

	return list
}

func FetchById(dbc *config.Dbconfig, query string) map[string]string {

	list := Fetch(dbc, query)

	return list[0]
}

func Insert(dbc *config.Dbconfig, table string, values []string) (int64, error) {

	db := config.Connect(dbc)

	args := reflect.ValueOf(values).Interface()
	n := len(values)
	var q string

	query := "INSERT INTO " + table + " VALUES("

	for i := 0; i < n; i++ {
		query += " ?, "

		if i+1 == n {
			q = TrimSuffixer(query, ",")
			q += ")"
		}
	}

	result, err := db.Exec(query, args)

	check(err)

	return result.LastInsertId()

}

func Update(dbc *config.Dbconfig, table string, fields []string, values []string) (int64, error) {

	db := config.Connect(dbc)

	args := reflect.ValueOf(values).Interface()

	var val string
	var q string

	query := "UPDATE " + table + " SET "

	for i := 0; i < len(values); i++ {
		query += fields[i] + " = " + reflect.ValueOf(values[i]).String() + ","

		if i+1 == len(values) {

			q = TrimSuffixer(query, ",")

			if strings.Contains(fields[i], "*") {
				val = strings.TrimSuffix(fields[i], "*")
				q += " WHERE " + val + " = " + reflect.ValueOf(values[i]).String() + ")"
			}
		}
	}

	result, err := db.Exec(query, args)

	check(err)

	return result.RowsAffected()

}

func Delete(dbc *config.Dbconfig, table string, field string, id string) (int64, error) {

	db := config.Connect(dbc)

	query := "DELETE FROM " + table + " WHERE " + field + " = " + id

	result, err := db.Exec(query, reflect.ValueOf(id).Interface())

	check(err)

	return result.LastInsertId()
}
