package conf_test

import (
	"github.com/cliod/conf"
	"testing"
)

type AppInfo struct {
	name    string // must be exported
	Author  string
	Version float64
}

func TestConfig(t *testing.T) {
	t.Log(conf.Reload())

	t.Log(conf.GetString("app"))

	t.Log(conf.GetInt("app.version"))
	t.Log(conf.GetFloat("app.version"))
	t.Log(conf.GetBool("app.version"))

	var (
		info AppInfo
		m    = make(map[string]interface{})
	)
	conf.ToStruct("app", &info)
	conf.ToStruct("app", &m)
	t.Logf("%#v", info)
	t.Logf("%#v", m)

}

func TestValue(t *testing.T) {
	var converter conf.Converter = new(conf.Yaml2PropsConverter)
	t.Log(conf.Conf().GetString("app"))
	var variable = converter.Convert(conf.Conf().Variable())
	props := variable.(*conf.Props)
	t.Log("=========== yaml -> props ============")
	t.Logf("%#v", props.Keys())

	converter = new(conf.Props2YamlConverter)
	variable = converter.Convert(variable)
	yaml := variable.(*conf.Yaml)
	t.Log("=========== props -> yaml ============")
	t.Logf("%#v", yaml.Variable())

	converter = new(conf.Yaml2JsonConverter)
	variable = converter.Convert(variable)
	json := variable.(*conf.Json)
	t.Log("=========== yaml -> json ============")
	t.Logf("%#v", json.Variable())

	converter = new(conf.Json2YamlConverter)
	variable = converter.Convert(variable)
	yaml = variable.(*conf.Yaml)
	t.Log("=========== json -> yaml ============")
	t.Logf("%#v", yaml.Variable())

	converter = new(conf.Yaml2PropsConverter)
	variable = converter.Convert(variable)
	props = variable.(*conf.Props)

	converter = new(conf.Props2JsonConverter)
	variable = converter.Convert(variable)
	json = variable.(*conf.Json)
	t.Log("=========== props -> json ============")
	t.Logf("%#v", json.Variable())

	converter = new(conf.Json2PropsConverter)
	variable = converter.Convert(variable)
	props = variable.(*conf.Props)
	t.Log("=========== json -> props ============")
	t.Logf("%#v", props.Variable())
}

func TestJson(t *testing.T) {
	json := conf.New("", "app.json", conf.JSON.String())

	t.Log(json.GetString("app"))
	t.Log(json.GetInt("app.version"))
	t.Log(json.GetFloat("app.version"))
	t.Log(json.GetBool("app.version"))

	info := AppInfo{}
	m := make(map[string]interface{})
	json.Struct("app", &info)
	json.Struct("app", &m)
	t.Logf("%#v", info)
	t.Logf("%#v", m)
}

func TestYaml(t *testing.T) {
	yaml := conf.New()

	t.Log(yaml.GetString("app"))
	t.Log(yaml.GetInt("app.version"))
	t.Log(yaml.GetFloat("app.version"))
	t.Log(yaml.GetBool("app.version"))

	info := AppInfo{}
	m := make(map[string]interface{})
	yaml.Struct("app", &info)
	yaml.Struct("app", &m)
	t.Logf("%#v", info)
	t.Logf("%#v", m)

	info.Version = 2
	yaml.Struct("app.version", &info)
	t.Logf("%#v", info)
}

func TestProps(t *testing.T) {
	props := conf.New("", "app.properties", conf.PROPS.String())

	t.Log(props.GetString("app"))
	t.Log(props.GetInt("app.version"))
	t.Log(props.GetFloat("app.version"))
	t.Log(props.GetBool("app.version"))

	info := AppInfo{}
	m := make(map[string]interface{})
	props.Struct("app", &info)
	props.Struct("app", &m)
	t.Logf("%#v", info)
	t.Logf("%#v", m)
}

func TestInclude(t *testing.T) {
	t.Log(conf.Get("include"))
	t.Log(conf.Get("app.include"))
}
