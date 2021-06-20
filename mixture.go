package conf

import (
	"errors"
	"reflect"
	"strings"
)

type Mixture struct {
	data *Json
}

func newMixture(variable StoreVariable, variables ...StoreVariable) *Mixture {
	if variables == nil {
		variables = append([]StoreVariable{})
	}
	variables = append(variables, variable)
	m := make(map[string]interface{})
	for _, store := range variables {
		value := store.Variable().Value()
		push(m, value, false)
	}
	return &Mixture{&Json{m}}
}

func push(m map[string]interface{}, value interface{}, replace bool) {
	var data map[string]interface{}
	switch val := value.(type) {
	case map[string]interface{}:
		data = val
	case map[interface{}]interface{}:
		data = y2j.yaml2Json(val)
	case map[string]string:
		data = make(map[string]interface{})
		p2j.props2Json(val, data, 0)
	}
	for k, v := range data {
		val, exist := m[k]
		if replace {
			if val == nil || reflect.TypeOf(val).Kind() != reflect.Map {
				m[k] = v
			} else {
				push(val.(map[string]interface{}), v, true)
			}
		} else {
			if !exist {
				m[k] = v
			}
		}
	}
}

func (m *Mixture) replace(profile StoreVariable) {
	push(m.data.data, profile.Variable().Value(), true)
}

func (m *Mixture) Variable() Variable {
	return newVal(m.data.data)
}

func (m *Mixture) Value(name string) Variable {
	return m.data.Value(name)
}

func (m *Mixture) Get(name string) interface{} {
	return m.Value(name).Value()
}

func (m *Mixture) GetString(name string) string {
	return m.Value(name).String()
}

func (m *Mixture) GetFloat(name string) float64 {
	return m.Value(name).Float()
}

func (m *Mixture) GetInt(name string) int {
	return m.Value(name).Int()
}

func (m *Mixture) GetBool(name string) bool {
	return m.Value(name).Bool()
}

func (m *Mixture) Struct(name string, receiver interface{}) {
	m.data.Struct(name, receiver)
}

func (m *Mixture) Load(filename string) error {
	extName := filename[:strings.LastIndex(filename, ".")]
	switch strings.ToLower(extName) {
	case "json":
		return m.data.Load(filename)
	case "yaml", "yml":
		y := new(Yaml)
		err := y.Load(filename)
		if err != nil {
			return err
		}
		m.data = y.Convert(y2j).(*Json)
		return nil
	case "properties", "conf":
		p := new(Props)
		err := p.Load(filename)
		if err != nil {
			return err
		}
		m.data = p.Convert(p2j).(*Json)
		return nil
	default:
		return errors.New("file format does not match")
	}
}
