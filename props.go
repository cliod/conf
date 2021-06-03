package conf

import (
	"github.com/magiconair/properties"
	"strings"
)

type Props struct {
	props *properties.Properties
}

func (p *Props) loadMap(data map[string]string) (err error) {
	p.props = properties.LoadMap(data)
	return
}

func (p *Props) loadByte(data []byte) (err error) {
	p.props, err = properties.Load(data, properties.UTF8)
	return
}

func (p *Props) Keys() (keys []string) {
	keys = p.props.Keys()
	return
}

func (p *Props) Load(filename string) (err error) {
	loader := &properties.Loader{Encoding: properties.UTF8}
	p.props, err = loader.LoadAll([]string{filename})
	return
}

func (p *Props) Value(name string) Variable {
	return &Value{p.props.GetString(name, "")}
}

func (p *Props) GetValue(name string) interface{} {
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
	props := p.props.FilterPrefix(name)
	if props.Len() > 0 {
		for key, value := range props.Map() {
			key = strings.TrimPrefix(key, name+".")
			switch receiver.(type) {
			case *map[string]string:
				(*receiver.(*map[string]string))[key] = value
			case *map[string]interface{}:
				(*receiver.(*map[string]interface{}))[key] = value
			default:
				err := p.setField(receiver, key, value)
				wLog(err)
			}
		}
	}
}

func (p *Props) Convert(converter Converter) KindVariable {
	return converter.Convert(p)
}

func (p *Props) Yaml() *Yaml {
	return p.Convert(P2Y).(*Yaml)
}

func (p *Props) Json() *Json {
	return p.Convert(P2J).(*Json)
}

func (p *Props) setField(receiver interface{}, name string, value interface{}) error {
	return SetFieldValue(receiver, name, value)
}
