package conf

import (
	"testing"
)

type AppInfo struct {
	name    string // must be exported
	Author  string
	Version float64
}

func TestConfig(t *testing.T) {
	t.Log(GetString("app"))

	t.Log(GetInt("app.version"))
	t.Log(GetFloat("app.version"))
	t.Log(GetBool("app.version"))

	var (
		info AppInfo
		m    = make(map[string]interface{})
	)
	ToStruct("app", &info)
	ToStruct("app", &m)
	t.Logf("%#v", info)
	t.Logf("%#v", m)

	props := New("", "app.properties", PROPS.String())
	val := props.Value("app.version")
	t.Log(val.String())
	info = AppInfo{}
	props.Struct("app", &info)
	t.Logf("%#v", info)

	json := New("", "app.json", JSON.String())
	val = json.Value("app.version")
	t.Log(val.String())
	info = AppInfo{}
	json.Struct("app", &info)
	t.Logf("%#v", info)
}

func TestValue(t *testing.T) {
	var converter Converter = new(Yaml2PropsConverter)
	var variable = converter.Convert(Conf().store)
	props := variable.(*Props)
	t.Log("=========== yaml -> props ============")
	t.Logf("%#v", props.props.Keys())

	converter = new(Props2YamlConverter)
	variable = converter.Convert(variable)
	yaml := variable.(*Yaml)
	t.Log("=========== props -> yaml ============")
	t.Logf("%#v", yaml.data)

	converter = new(Yaml2JsonConverter)
	variable = converter.Convert(variable)
	json := variable.(*Json)
	t.Log("=========== yaml -> json ============")
	t.Logf("%#v", json.data)

	converter = new(Json2YamlConverter)
	variable = converter.Convert(variable)
	yaml = variable.(*Yaml)
	t.Log("=========== json -> yaml ============")
	t.Logf("%#v", yaml.data)

	converter = new(Yaml2PropsConverter)
	variable = converter.Convert(variable)
	props = variable.(*Props)

	converter = new(Props2JsonConverter)
	variable = converter.Convert(variable)
	json = variable.(*Json)
	t.Log("=========== props -> json ============")
	t.Logf("%#v", json.data)

	converter = new(Json2PropsConverter)
	variable = converter.Convert(variable)
	props = variable.(*Props)
	t.Log("=========== json -> props ============")
	t.Logf("%#v", props.props.Map())
}
