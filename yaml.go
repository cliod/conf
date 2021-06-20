package conf

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"gopkg.in/yaml.v2"
)

type Yaml struct {
	data map[interface{}]interface{}
}

func (y *Yaml) Variable() Variable {
	return newVal(y.data)
}

func (y *Yaml) LoadBytes(data []byte) (err error) {
	y.data = make(map[interface{}]interface{})
	err = yaml.Unmarshal(data, &y.data)
	eLog(err)
	return err
}

func (y *Yaml) Load(filename string) error {
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
	return y.LoadBytes(bs)
}

func (y *Yaml) Value(name string) Variable {
	var val interface{}
	name = strings.Trim(name, " ")
	if name == "" {
		val = y.data
	}
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
			val = value
			break
		}
		if reflect.TypeOf(value).Kind() == reflect.Map {
			data = value.(map[interface{}]interface{})
		}
	}
	return newVal(val)
}

func (y *Yaml) Get(name string) interface{} {
	return y.Value(name).Value()
}

func (y *Yaml) GetString(name string) string {
	return y.Value(name).String()
}

func (y *Yaml) GetFloat(name string) float64 {
	return y.Value(name).Float()
}

func (y *Yaml) GetInt(name string) int {
	return y.Value(name).Int()
}

func (y *Yaml) GetBool(name string) bool {
	return y.Value(name).Bool()
}

func (y *Yaml) Struct(name string, receiver interface{}) {
	value := y.Value(name)
	switch val := value.Value().(type) {
	case string, float64, int64, int, bool:
		if strings.Contains(name, ".") {
			name = name[strings.LastIndex(name, ".")+1:]
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

func (y *Yaml) Convert(converter Converter) KeyVariable {
	return converter.Convert(y)
}
