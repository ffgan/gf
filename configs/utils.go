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
		return fmt.Errorf("obj 必须是非空指针")
	}
	v = v.Elem()

	fieldName, ok := tagToField[key]
	if !ok {
		return fmt.Errorf("未找到 tag=%s 对应字段", key)
	}

	f := v.FieldByName(fieldName)
	if !f.IsValid() || !f.CanSet() {
		return fmt.Errorf("字段 %s 无法设置", fieldName)
	}

	if f.Kind() == reflect.String {
		f.SetString(value)
		return nil
	}
	return fmt.Errorf("不支持类型: %s", f.Kind())
}
