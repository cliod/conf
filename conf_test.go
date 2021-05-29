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
	t.Log(conf.GetString("app"))

	t.Log(conf.GetInt("app.version"))
	t.Log(conf.GetFloat("app.version"))
	t.Log(conf.GetBool("app.version"))

	var info AppInfo
	conf.ToStruct("app", &info)
	t.Logf("%#v", info)
}
