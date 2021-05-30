package conf

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

type Json struct {
	data map[string]interface{}
}

func (j *Json) Load(filename string) error {
	var (
		bs  []byte
		err error
	)
	if strings.HasPrefix(strings.ToLower(filename), "http") {
		var resp *http.Response
		resp, err = http.Get(filename)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		bs, err = ioutil.ReadAll(resp.Body)
	} else {
		bs, err = ioutil.ReadFile(filename)
	}
	if err != nil {
		return err
	}
	j.data = make(map[string]interface{})
	err = json.Unmarshal(bs, &j.data)
	eLog(err)
	return err
}

func (j *Json) Value(name string) Variable {
	subKeys := strings.Split(name, ".")
	data := j.data
	for index, key := range subKeys {
		if strings.Trim(key, " ") == "" {
			continue
		}
		value, ok := data[key]
		if !ok {
			break
		}
		if (index + 1) == len(subKeys) {
			return newVal(value)
		}
		if reflect.TypeOf(value).Kind() == reflect.Map {
			data = value.(map[string]interface{})
		}
	}
	return nil
}

func (j *Json) GetValue(name string) interface{} {
	return j.Value(name).Value()
}

func (j *Json) GetString(name string) string {
	return j.Value(name).String()
}

func (j *Json) GetFloat(name string) float64 {
	return j.Value(name).Float()
}

func (j *Json) GetInt(name string) int {
	return j.Value(name).Int()
}

func (j *Json) GetBool(name string) bool {
	return j.Value(name).Bool()
}

func (j *Json) Struct(name string, receiver interface{}) {
	j.Value(name).Struct(receiver)
}

func (j *Json) Convert(converter Converter) KindVariable {
	return converter.Convert(j)
}

func (j *Json) Props() *Props {
	return j.Convert(J2P).(*Props)
}

func (j *Json) Yaml() *Yaml {
	return j.Convert(J2Y).(*Yaml)
}
