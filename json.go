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

// Keys returns root keys
func (j *Json) Keys() (keys []string) {
	for key := range j.data {
		keys = append(keys, key)
	}
	return
}

func (j *Json) Variable() Variable {
	return newVal(j.data)
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
		defer func() {
			err = resp.Body.Close()
			if err != nil {
				wLog(err)
			}
		}()
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
	var val interface{}
	name = strings.Trim(name, " ")
	if name == "" {
		val = j.data
	}
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
			val = value
			break
		}
		if reflect.TypeOf(value).Kind() == reflect.Map {
			data = value.(map[string]interface{})
		}
	}
	return newVal(val)
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
	value := j.Value(name)
	switch val := value.Value().(type) {
	case string, float64, int64, int, bool:
		if strings.Contains(name, ".") {
			name = name[strings.LastIndex(name, "."):]
		}
		err := SetFieldValue(receiver, name, val)
		wLog(err)
	}
	extVal, ok := value.(ExtVariable)
	if ok {
		extVal.Struct(receiver)
		return
	}
	bs, err := json.Marshal(value.Value())
	if err != nil {
		eLog(err)
		return
	}
	err = json.Unmarshal(bs, receiver)
	eLog(err)
}

func (j *Json) Convert(converter Converter) KindVariable {
	return converter.Convert(j)
}

func (j *Json) Props() *Props {
	return j.Convert(j2p).(*Props)
}

func (j *Json) Yaml() *Yaml {
	return j.Convert(j2y).(*Yaml)
}
