package conf

import (
	"errors"
	"reflect"
	"unicode"
)

func SetFieldValue(receiver interface{}, fieldName string, value interface{}, fs ...func(fieldType reflect.Type, value interface{})) error {
	for i, v := range fieldName {
		fieldName = string(unicode.ToUpper(v)) + fieldName[i+1:]
		break
	}
	structValue := reflect.Indirect(reflect.ValueOf(receiver))
	structFieldValue := structValue.FieldByName(fieldName)
	if !structFieldValue.IsValid() {
		return errors.New("no such field: " + fieldName)
	}
	if !structFieldValue.CanSet() {
		return errors.New("can not set field value: " + fieldName)
	}
	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)
	if structFieldType.Kind() == reflect.Struct && val.Kind() == reflect.Map {
		switch nVal := val.Interface().(type) {
		case map[interface{}]interface{}:
			for k, v := range nVal {
				err := SetFieldValue(structFieldValue.Addr().Interface(), newVal(k).String(), v)
				wLog(err)
			}
		case map[string]interface{}:
			for k, v := range nVal {
				err := SetFieldValue(structFieldValue.Addr().Interface(), k, v)
				wLog(err)
			}
		}
	} else {
		if len(fs) == 0 {
			if structFieldType != val.Type() {
				return errors.New("provided value type didn't match field type")
			}
			structFieldValue.Set(val)
		} else {
			for _, f := range fs {
				f(structFieldType, value)
			}
		}
	}
	return nil
}
