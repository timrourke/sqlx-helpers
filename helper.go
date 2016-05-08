package helper

import (
	"fmt"
	"reflect"
	"strings"
)

// Easy creation of wheres for sqlx named queries
// eg.
// where id in (:id)
func CreateWhere(model interface{}, pairs map[string]interface{}) (string, map[string]interface{}) {
	parts := []string{}
	expanded := map[string]interface{}{}

	modelType := reflect.TypeOf(model)

	for key, values := range pairs {
		col := strings.ToLower(key)
		field, ok := modelType.FieldByName(key)
		if ok {
			col = field.Tag.Get("db")
			if col == "" {
				col = strings.ToLower(key)
			}
		}

		switch reflect.TypeOf(values).Kind() {
		case reflect.Slice:
			slc := reflect.ValueOf(values)
			subs := []string{}
			for p := 0; p < slc.Len(); p++ {
				idxKey := fmt.Sprintf("%v%d", key, p)
				expanded[idxKey] = slc.Index(p).Interface()
				subs = append(subs, ":"+idxKey)
			}
			parts = append(parts, col+" in ("+strings.Join(subs, ",")+")")
		default:
			parts = append(parts, col+" = :"+key)
		}
	}

	return "where " + strings.Join(parts, " and "), expanded
}

// Given a struct, create an insert statement for it.
// eg.
// (id, name, stuff) values (?, ?, ?)
func CreateInsert(model interface{}, exclude ...string) string {
	return createInsertOrUpdate(model, true, exclude)
}

// Given a struct, create an update statement for it.
// eg.
// id = ?, name = ?, stuff = ?
func CreateUpdate(model interface{}, exclude ...string) string {
	return createInsertOrUpdate(model, false, exclude)
}

func createInsertOrUpdate(model interface{}, insert bool, exclude []string) string {
	modelType := reflect.TypeOf(model)
	modelField := modelType.Field(0)
	tableName := modelField.Tag.Get("tableName")

	columns := []string{}
	names := []string{}

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		col := field.Tag.Get("db")
		if col == "-" || contains(exclude, col) {
			continue
		} else if col == "" {
			columns = append(columns, strings.ToLower(field.Name))
		} else {
			columns = append(columns, col)
		}
		names = append(names, "?")
	}
	if insert {
		return "INSERT INTO " + tableName + " (" + strings.Join(columns, ", ") + ") values (" + strings.Join(names, ", ") + ")"
	}

	update := []string{}
	for i := 0; i < len(columns); i++ {
		update = append(update, columns[i]+"="+names[i])
	}
	return "UPDATE " + tableName + " SET " + strings.Join(update, ", ")
}

func contains(s []string, e string) bool {
	for _, ex := range s {
		if e == ex {
			return true
		}
	}
	return false
}
