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
	Struct(interface{})
}

type Value struct {
	value interface{}
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

func (v *Value) Struct(receiver interface{}) {
	bs, err := json.Marshal(v.value)
	if err != nil {
		eLog(err)
		return
	}
	err = json.Unmarshal(bs, receiver)
	eLog(err)
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
