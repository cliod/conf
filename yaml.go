package conf

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"reflect"
	"strings"
	"unicode"
)

type Yaml struct {
	Store
	data map[interface{}]interface{}
}

func (y *Yaml) Load(filename string) error {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	m := make(map[interface{}]interface{})
	err = yaml.Unmarshal(bs, &m)
	eLog(err)
	y.data = m
	y.StoreVariable = y
	return err
}

func (y *Yaml) Value(name string) Variable {
	subKeys := strings.Split(name, ".")
	data := y.data
	for index, key := range subKeys {
		if strings.Trim(key, " ") == "" {
			continue
		}
		value, ok := data[key]
		if !ok {
			break
		}
		if (index + 1) == len(subKeys) {
			return &Value{value}
		}
		if reflect.TypeOf(value).Kind() == reflect.Map {
			data = value.(map[interface{}]interface{})
		}
	}
	return nil
}

func (y *Yaml) Struct(name string, receiver interface{}) {
	value := y.GetValue(name)
	switch value.(type) {
	case string:
		err := y.setField(receiver, name, value)
		wLog(err)
	case map[interface{}]interface{}:
		y.mapToStruct(value.(map[interface{}]interface{}), receiver)
	}
}

func (y *Yaml) mapToStruct(m map[interface{}]interface{}, receiver interface{}) interface{} {
	for key, value := range m {
		switch key.(type) {
		case string:
			err := y.setField(receiver, key.(string), value)
			wLog(err)
		}
	}
	return receiver
}

func (y *Yaml) setField(receiver interface{}, name string, value interface{}) error {

	for i, v := range name {
		name = string(unicode.ToUpper(v)) + name[i+1:]
		break
	}

	structValue := reflect.Indirect(reflect.ValueOf(receiver))
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return errors.New("no such field: " + name)
	}

	if !structFieldValue.CanSet() {
		return errors.New("can not set field value: " + name)
	}

	structFieldType := structFieldValue.Type()
	val := reflect.ValueOf(value)

	if structFieldType.Kind() == reflect.Struct && val.Kind() == reflect.Map {
		vint := val.Interface()
		switch vint.(type) {
		case map[interface{}]interface{}:
			for k, v := range vint.(map[interface{}]interface{}) {
				err := y.setField(structFieldValue.Addr().Interface(), k.(string), v)
				wLog(err)
			}
		case map[string]interface{}:
			for k, v := range vint.(map[string]interface{}) {
				err := y.setField(structFieldValue.Addr().Interface(), k, v)
				wLog(err)
			}
		}
	} else {
		if structFieldType != val.Type() {
			return errors.New("provided value type didn't match field type")
		}
		structFieldValue.Set(val)
	}
	return nil
}
