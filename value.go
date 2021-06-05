package conf

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

type Variable interface {
	Value() interface{}
	String() string
	Float() float64
	Int() int
	Bool() bool
}

type ExtVariable interface {
	Variable
	Struct(interface{})
}

type Value struct {
	value interface{}
}

func newVal(value interface{}) *Value {
	return &Value{value}
}

func (v *Value) Value() interface{} {
	return v.value
}

func (v *Value) String() string {
	if v.value == nil {
		return ""
	}
	switch val := v.value.(type) {
	case string:
		return val
	default:
		return fmt.Sprint(val)
	}
}

func (v *Value) Float() float64 {
	if v.value == nil {
		return 0
	}
	switch val := v.value.(type) {
	case string:
		digital, err := strconv.ParseFloat(val, 10)
		if err != nil {
			log.Panicln("[CONF ERROR]: ", err)
		}
		return digital
	case float64:
		return val
	case int:
		return float64(val)
	case int64:
		return float64(val)
	case bool:
		if val {
			return 1
		}
		return 0
	default:
		return 0
	}
}

func (v *Value) Int() int {
	if v.value == nil {
		return 0
	}
	switch val := v.value.(type) {
	case string:
		digital, err := strconv.Atoi(val)
		if err != nil {
			log.Panicln("[CONF ERROR]: ", err)
		}
		return digital
	case float64:
		return int(val)
	case int:
		return val
	case int64:
		return int(val)
	case bool:
		if val {
			return 1
		}
		return 0
	default:
		return 0
	}
}

func (v *Value) Bool() bool {
	if v.value == nil {
		return false
	}
	switch val := v.value.(type) {
	case string:
		b, err := v.parseBool(val)
		if err != nil {
			log.Panicln("[CONF ERROR]: ", err)
		}
		return b
	case float64:
		return val > 0.0
	case int:
		return val > 0
	case int64:
		return val > 0
	case bool:
		return val
	default:
		return false
	}
}

func (v *Value) parseBool(val string) (bool, error) {
	switch val {
	case "1", "t", "T", "true", "TRUE", "True", "on", "ON", "Y", "y", "YES", "yes", "Yes":
		return true, nil
	case "0", "f", "F", "false", "FALSE", "False", "off", "OFF", "n", "N", "NO", "no", "No":
		return false, nil
	}
	return false, &strconv.NumError{Func: "ParseBool", Num: val, Err: strconv.ErrSyntax}
}

func (v *Value) Struct(receiver interface{}) {
	switch val := v.value.(type) {
	case string:
		err := json.Unmarshal([]byte(val), receiver)
		wLog(err)
	case map[interface{}]interface{}:
		switch receiver.(type) {
		case *map[string]string:
			for key, value := range val {
				(*receiver.(*map[string]string))[newVal(key).String()] = newVal(value).String()
			}
		case *map[string]interface{}:
			v.mapToBMap(val, receiver.(*map[string]interface{}))
		default:
			v.mapToStruct(val, receiver)
		}
	case map[string]interface{}:
		switch receiver.(type) {
		case *map[string]string:
			for key, value := range val {
				(*receiver.(*map[string]string))[newVal(key).String()] = newVal(value).String()
			}
		case *map[string]interface{}:
			for key, value := range val {
				(*receiver.(*map[string]interface{}))[newVal(key).String()] = value
			}
		default:
			for key, value := range val {
				err := SetFieldValue(receiver, key, value)
				wLog(err)
			}
		}
	}
}

func (v *Value) mapToBMap(val map[interface{}]interface{}, receiver *map[string]interface{}) {
	for key, value := range val {
		if sv, ok := value.(map[interface{}]interface{}); ok {
			m := make(map[string]interface{})
			v.mapToBMap(sv, &m)
			(*receiver)[newVal(key).String()] = m
		} else {
			(*receiver)[newVal(key).String()] = value
		}
	}
}

func (v *Value) mapToStruct(m map[interface{}]interface{}, receiver interface{}) interface{} {
	for key, value := range m {
		switch k := key.(type) {
		case string:
			err := SetFieldValue(receiver, k, value)
			wLog(err)
		default:
			err := SetFieldValue(receiver, newVal(k).String(), value)
			wLog(err)
		}
	}
	return receiver
}
