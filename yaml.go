package conf

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

type Yaml struct {
	data map[interface{}]interface{}
}

// Keys returns root keys
func (y *Yaml) Keys() (keys []string) {
	for key := range y.data {
		var vKey = &Value{key}
		keys = append(keys, vKey.String())
	}
	return
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
		defer resp.Body.Close()
		bs, err = ioutil.ReadAll(resp.Body)
	} else {
		bs, err = ioutil.ReadFile(filename)
	}
	if err != nil {
		return err
	}
	y.data = make(map[interface{}]interface{})
	err = yaml.Unmarshal(bs, &y.data)
	eLog(err)
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
			return newVal(value)
		}
		if reflect.TypeOf(value).Kind() == reflect.Map {
			data = value.(map[interface{}]interface{})
		}
	}
	return nil
}

func (y *Yaml) GetValue(name string) interface{} {
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
	value := y.GetValue(name)
	switch val := value.(type) {
	case string:
		err := y.setField(receiver, name, value)
		wLog(err)
	case map[interface{}]interface{}:
		y.mapToStruct(val, receiver)
	}
}

func (y *Yaml) Convert(converter Converter) KindVariable {
	return converter.Convert(y)
}

func (y *Yaml) Props() *Props {
	return y.Convert(Y2P).(*Props)
}

func (y *Yaml) Json() *Json {
	return y.Convert(Y2J).(*Json)
}

func (y *Yaml) mapToStruct(m map[interface{}]interface{}, receiver interface{}) interface{} {
	for key, value := range m {
		switch k := key.(type) {
		case string:
			err := y.setField(receiver, k, value)
			wLog(err)
		default:
			err := y.setField(receiver, newVal(k).String(), value)
			wLog(err)
		}
	}
	return receiver
}

func (y *Yaml) setField(receiver interface{}, name string, value interface{}) error {
	ref := reflect.TypeOf(receiver)
	if ref.Elem().Kind() == reflect.Map {
		m, ok := receiver.(*map[interface{}]interface{})
		if ok {
			(*m)[name] = value
			return nil
		}
		m1, ok := receiver.(*map[string]interface{})
		if ok {
			(*m1)[name] = value
			return nil
		}
		return errors.New("the receiver has type error")
	} else {
		return SetFieldValue(receiver, name, value)
	}
}
