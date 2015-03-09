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
