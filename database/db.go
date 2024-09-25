package database

import (
	"database/sql"
	"fmt"
	"reflect"

	_ "github.com/mattn/go-sqlite3"
)

// query: get data
// exec: insert, update, delete
// prepare: when need to execute statements several times

func Connect() *sql.DB {
	const file string = "./qawwali-glossary.db"
	db, err := sql.Open("sqlite3", file)

	if err != nil {
		panic(err)
	}

	return db
}

var data = [][]interface{}{
	{"1", "2"},
	{"4", "5"},
	{"7", "8"},
}

func QueryRows(db *sql.DB, query string, printDesc string) []string {

	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	var finalRes []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			panic(err)
		}
		finalRes = append(finalRes, tableName)
	}
	defer rows.Close()

	if printDesc != "" {
		fmt.Println(printDesc, ":", finalRes)
	}

	return finalRes

}

func RecursiveEntries(db *sql.DB, data []interface{}, table string, query string) []interface{} {
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	if len(data) == 0 {
		return nil
	}

	var failedEntries []interface{}
	// Use reflection to get field names from the first struct
	structType := reflect.TypeOf(data[0])
	numFields := structType.NumField()

	// Create column names and placeholders based on struct fields
	columns := make([]string, numFields)
	placeholders := make([]string, numFields)
	for i := 0; i < numFields; i++ {
		columns[i] = structType.Field(i).Name
		placeholders[i] = "?"
	}

	// Build the INSERT SQL statement
	columnList := sqlJoin(columns, ", ")
	placeholderList := sqlJoin(placeholders, ", ")
	statement := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) %s", table, columnList, placeholderList, query)
	stmt, err := tx.Prepare(statement)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	// Iterate over the structs and execute the insert
	for _, item := range data {
		// Use reflection to get the values of the struct fields
		structValue := reflect.ValueOf(item)
		values := make([]interface{}, numFields)
		for i := 0; i < numFields; i++ {
			values[i] = structValue.Field(i).Interface()
		}

		// Execute the insert statement
		_, err := stmt.Exec(values...)
		if err != nil {
			fmt.Println(err)
			failedEntries = append(failedEntries, item)
		}
	}
	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	fmt.Println(len(data)-len(failedEntries), "entries inserted", len(failedEntries), "entries failed")

	return failedEntries
}

// Utility function to join a slice of strings with a delimiter
func sqlJoin(elements []string, sep string) string {
	result := ""
	for i, element := range elements {
		if i > 0 {
			result += sep
		}
		result += element
	}
	return fmt.Sprintf("%s", result)
}
