package conf

import (
	"strings"
)

var (
	y2p *Yaml2PropsConverter
	y2j *Yaml2JsonConverter
	p2y *Props2YamlConverter
	p2j *Props2JsonConverter
	j2y *Json2YamlConverter
	j2p *Json2PropsConverter
	m2j *Mixture2JsonConverter
)

type Convertable interface {
	Convert(Converter) KindVariable
}

type Converter interface {
	Convert(KindVariable) KindVariable
}

type Mixture2JsonConverter struct {
}

func (m *Mixture2JsonConverter) Convert(variable KindVariable) KindVariable {
	mix, ok := variable.(*Mixture)
	if !ok {
		return new(Json)
	}
	return mix.data
}

type Yaml2PropsConverter struct {
}

func (conv *Yaml2PropsConverter) Convert(variable KindVariable) KindVariable {
	yaml, ok := variable.(*Yaml)
	props := new(Props)
	if !ok {
		return props
	}
	err := props.loadMap(conv.yaml2Props("", yaml.data))
	eLog(err)
	return props
}

func (conv *Yaml2PropsConverter) yaml2Props(preKey string, data map[interface{}]interface{}) (res map[string]string) {
	res = make(map[string]string)
	for key, value := range data {
		keyStr := newVal(key).String()
		if preKey != "" {
			keyStr = preKey + "." + keyStr
		}
		switch val := value.(type) {
		case map[interface{}]interface{}:
			subRes := conv.yaml2Props(keyStr, val)
			for k, v := range subRes {
				res[k] = v
			}
		default:
			res[keyStr] = newVal(value).String()
		}
	}
	return
}

type Props2YamlConverter struct {
}

func (conv *Props2YamlConverter) Convert(variable KindVariable) KindVariable {
	props, ok := variable.(*Props)
	yaml := new(Yaml)
	if !ok {
		return yaml
	}
	yaml.data = make(map[interface{}]interface{})
	conv.props2Yaml(props.props.Map(), yaml.data, 0)
	return yaml
}

func (conv *Props2YamlConverter) props2Yaml(data map[string]string, levelRes map[interface{}]interface{}, level int) {
	for key, value := range data {
		subKeys := strings.Split(key, ".")
		for index, subKey := range subKeys {
			if index != level {
				continue
			}
			subVal, exist := levelRes[subKey]
			if !exist {
				subVal = make(map[interface{}]interface{})
				levelRes[subKey] = subVal
			}
			switch val := subVal.(type) {
			case map[interface{}]interface{}:
				conv.props2Yaml(data, val, level+1)
			}
			if index+1 == len(subKeys) {
				levelRes[subKey] = value
			}
		}
	}
	return
}

type Yaml2JsonConverter struct {
}

func (conv *Yaml2JsonConverter) Convert(variable KindVariable) KindVariable {
	yaml, ok := variable.(*Yaml)
	json := new(Json)
	if !ok {
		return json
	}
	json.data = conv.yaml2Json(yaml.data)
	return json
}

func (conv *Yaml2JsonConverter) yaml2Json(data map[interface{}]interface{}) (res map[string]interface{}) {
	res = make(map[string]interface{})
	for key, value := range data {
		switch val := value.(type) {
		case map[interface{}]interface{}:
			res[newVal(key).String()] = conv.yaml2Json(val)
		default:
			res[newVal(key).String()] = value
		}
	}
	return
}

type Json2YamlConverter struct {
}

func (conv *Json2YamlConverter) Convert(variable KindVariable) KindVariable {
	json, ok := variable.(*Json)
	yaml := new(Yaml)
	if !ok {
		return yaml
	}
	yaml.data = conv.json2Yaml(json.data)
	return yaml
}

func (conv *Json2YamlConverter) json2Yaml(data map[string]interface{}) (res map[interface{}]interface{}) {
	res = make(map[interface{}]interface{})
	for key, value := range data {
		switch val := value.(type) {
		case map[string]interface{}:
			res[key] = conv.json2Yaml(val)
		default:
			res[key] = value
		}
	}
	return
}

type Json2PropsConverter struct {
}

func (conv *Json2PropsConverter) Convert(variable KindVariable) KindVariable {
	json, ok := variable.(*Json)
	props := new(Props)
	if !ok {
		return props
	}
	err := props.loadMap(conv.json2Props("", json.data))
	eLog(err)
	return props
}

func (conv *Json2PropsConverter) json2Props(preKey string, data map[string]interface{}) (res map[string]string) {
	res = make(map[string]string)
	for key, value := range data {
		keyStr := newVal(key).String()
		if preKey != "" {
			keyStr = preKey + "." + keyStr
		}
		switch val := value.(type) {
		case map[string]interface{}:
			subRes := conv.json2Props(keyStr, val)
			for k, v := range subRes {
				res[k] = v
			}
		default:
			res[keyStr] = newVal(value).String()
		}
	}
	return
}

type Props2JsonConverter struct {
}

func (conv *Props2JsonConverter) Convert(variable KindVariable) KindVariable {
	props, ok := variable.(*Props)
	json := new(Json)
	if !ok {
		return json
	}
	json.data = make(map[string]interface{})
	conv.props2Json(props.props.Map(), json.data, 0)
	return json
}

func (conv *Props2JsonConverter) props2Json(data map[string]string, levelRes map[string]interface{}, level int) {
	for key, value := range data {
		subKeys := strings.Split(key, ".")
		for index, subKey := range subKeys {
			if index != level {
				continue
			}
			subVal, exist := levelRes[subKey]
			if !exist {
				subVal = make(map[string]interface{})
				levelRes[subKey] = subVal
			}
			switch val := subVal.(type) {
			case map[string]interface{}:
				conv.props2Json(data, val, level+1)
			}
			if index+1 == len(subKeys) {
				levelRes[subKey] = value
			}
		}
	}
	return
}
