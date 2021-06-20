package conf

import (
	"github.com/magiconair/properties"
	"strings"
)

type Props struct {
	props *properties.Properties
}

func (p *Props) LoadMap(data map[string]string) (err error) {
	p.props = properties.LoadMap(data)
	return
}

func (p *Props) LoadBytes(data []byte) (err error) {
	p.props, err = properties.Load(data, properties.UTF8)
	return
}

func (p *Props) Keys() (keys []string) {
	keys = p.props.Keys()
	return
}

func (p *Props) Variable() Variable {
	return newVal(p.props.Map())
}

func (p *Props) Load(filename string) (err error) {
	loader := &properties.Loader{Encoding: properties.UTF8}
	p.props, err = loader.LoadAll([]string{filename})
	return
}

func (p *Props) Value(name string) Variable {
	return &Value{p.props.GetString(name, "")}
}

func (p *Props) Get(name string) interface{} {
	return p.Value(name).Value()
}

func (p *Props) GetString(name string) string {
	return p.Value(name).String()
}

func (p *Props) GetFloat(name string) float64 {
	return p.Value(name).Float()
}

func (p *Props) GetInt(name string) int {
	return p.Value(name).Int()
}

func (p *Props) GetBool(name string) bool {
	return p.Value(name).Bool()
}

func (p *Props) Struct(name string, receiver interface{}) {
	var (
		m = make(map[string]string)
		j = make(map[string]interface{})
	)
	for key, value := range p.props.Map() {
		if key == name {
			continue
		}
		if strings.HasPrefix(key, name) {
			key = strings.TrimPrefix(key, name+".")
			m[key] = value
		}
	}
	p2j.props2Json(m, j, 0)
	for key, value := range j {
		if strings.Contains(key, ".") {
			continue
		}
		switch receiver.(type) {
		case *map[string]string:
			(*receiver.(*map[string]string))[key] = newVal(value).String()
		case *map[string]interface{}:
			(*receiver.(*map[string]interface{}))[key] = value
		default:
			err := SetFieldValue(receiver, key, value)
			wLog(err)
		}
	}
}

func (p *Props) Convert(converter Converter) KeyVariable {
	return converter.Convert(p)
}
