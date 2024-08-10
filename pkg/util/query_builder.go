package util

import (
	"fmt"
	"reflect"
	"strings"
)

func buildSQLQuery(dataStruct interface{}, i int, skipID bool) (string, string, []interface{}, int) {
	v := reflect.ValueOf(dataStruct)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()

	if t.Kind() != reflect.Struct {
		panic("BuildSQLQuery expects a struct or a pointer to a struct")
	}

	var fields []string
	var placeholders []string
	var args []interface{}

	j := 0
	skip := 0
	for j < t.NumField() {
		field := t.Field(j)
		dbTag := field.Tag.Get("db")
		if dbTag == "" || (dbTag == "id" && skipID) {
			skip++
			j++
			continue
		}
		fields = append(fields, dbTag)
		placeholders = append(placeholders, fmt.Sprintf("$%d", j+i-skip))
		args = append(args, v.Field(j).Interface())
		j++
	}
	fieldsStr := strings.Join(fields, ", ")
	placeholdersStr := strings.Join(placeholders, ", ")

	return fieldsStr, placeholdersStr, args, i + j
}

// BuildSQLQuery builds a SQL query from a struct that implements db tags
// The function returns the fields and placeholders for the query and the number of fields
// so the SQL query can be build dynamically
func BuildSQLInsertQuery(dataStruct interface{}, i int) (string, string, []interface{}, int) {
	return buildSQLQuery(dataStruct, i, true)
}

// BuildSQLSelectQuery builds a SQL select query from a struct that implements db tags
// The function returns the fields and placeholders for the query and the number of fields
// so the SQL query can be build dynamically
func BuildSQLSelectQuery(dataStruct interface{}, i int) (string, string, []interface{}, int) {
	return buildSQLQuery(dataStruct, i, false)
}

// BuildSQLUpdateQuery builds a SQL update query from a struct that implements db tags
// The function returns the fields and placeholders for the query and the number of fields
// so the SQL query can be build dynamically
func BuildSQLUpdateQuery(dataStruct interface{}, i int) (string, []interface{}, int) {
	v := reflect.ValueOf(dataStruct)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()

	var fields []string
	var args []interface{}

	j := 0
	skip := 0
	for j < t.NumField() {
		field := t.Field(j)
		dbTag := field.Tag.Get("db")
		// For update queries we don't want to update the ID
		if dbTag == "" || dbTag == "id" {
			skip++
			j++
			continue
		}
		fields = append(fields, fmt.Sprintf("%s = $%d", dbTag, j+i-skip))
		args = append(args, v.Field(j).Interface())
		j++
	}

	fieldsStr := strings.Join(fields, ", ")

	return fieldsStr, args, i + j
}

func BuildSQLCreateTableQuery(tableName string, dataStruct interface{}) (string, error) {
	v := reflect.ValueOf(dataStruct)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()

	if t.Kind() != reflect.Struct {
		return "", fmt.Errorf("BuildSQLCreateTableQuery expects a struct or a pointer to a struct")
	}

	var columns []string

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag == "" {
			continue
		}

		// Determine the SQL type based on the Go type
		var sqlType string
		c := field.Type.Kind()
		switch c {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			sqlType = "INTEGER"
		case reflect.Float32, reflect.Float64:
			sqlType = "FLOAT"
		case reflect.String:
			sqlType = "TEXT"
		case reflect.Bool:
			sqlType = "BOOLEAN"
		case reflect.Struct:
			switch field.Type.String() {
			case "sql.NullInt64":
				sqlType = "INTEGER NULL"
			default:
				return "", fmt.Errorf("unsupported field type: %s", field.Type.String())
			}
		case reflect.Ptr:
			// Handle pointers by determining the underlying type
			switch field.Type.Elem().Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				sqlType = "INTEGER"
			case reflect.Float32, reflect.Float64:
				sqlType = "FLOAT"
			case reflect.String:
				sqlType = "TEXT"
			case reflect.Bool:
				sqlType = "BOOLEAN"
			default:
				return "", fmt.Errorf("unsupported field type: %s", field.Type.String())
			}
		default:
			return "", fmt.Errorf("unsupported field type: %s", field.Type.String())
		}
		if dbTag == "id" {
			sqlType = "SERIAL PRIMARY KEY"
		}

		columns = append(columns, fmt.Sprintf("%s %s", dbTag, sqlType))
	}

	columnsStr := strings.Join(columns, ", ")
	createTableQuery := fmt.Sprintf("CREATE TABLE %s (%s);", tableName, columnsStr)

	return createTableQuery, nil
}
