package configs

import (
	"fmt"
	"reflect"
)

var tagToField = map[string]string{}

func init() {
	t := reflect.TypeOf(Title{})
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("gf")
		if tag != "" {
			tagToField[tag] = field.Name
		}
	}
}

func SetByGFTag(obj any, key, value string) error {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("obj must not be a nil pointer")
	}
	v = v.Elem()

	fieldName, ok := tagToField[key]
	if !ok {
		return fmt.Errorf("tag=%s not found", key)
	}

	f := v.FieldByName(fieldName)
	if !f.IsValid() || !f.CanSet() {
		return fmt.Errorf(" %s can't be set", fieldName)
	}

	if f.Kind() == reflect.String {
		f.SetString(value)
		return nil
	}
	return fmt.Errorf("unsupported type: %s", f.Kind())
}
