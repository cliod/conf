package conf

import (
	"errors"
	"reflect"
	"strings"
)

type Config struct {
	store StoreVariable // variables of configuration file

	dir    string // dir of configuration file, default: conf/
	name   string // name of configuration file, default: app.yaml
	isInit bool   // is initialized
	cType  CType  // type of configuration file, default: yaml

	includes []interface{}
	profile  StoreVariable
	// todo includes, profiles(cover)
}

func New(params ...string) *Config {
	var (
		confDir  string
		confName string
		cType    CType
	)
	if len(params) > 0 {
		confDir = params[0]
		if len(params) > 1 {
			confName = params[1]
			if len(params) > 2 {
				ct, ok := CTypeNames[params[2]]
				if ok {
					cType = CType(ct.(int))
				}
			}
		}
	}
	if strings.Trim(confDir, " ") == "" || confDir[0] != '/' || confDir[1] != ':' {
		confDir = rootPath() + defPath
	}
	if strings.Trim(confName, " ") == "" {
		confName = defName
	}
	return newConf(confDir, confName, cType)
}

func newConf(dir, name string, ct CType) *Config {
	var store StoreVariable
	switch ct {
	default:
		fallthrough
	case YAML:
		store = new(Yaml)
	case PROPS:
		store = new(Props)
	case JSON:
		store = new(Json)
	}
	return &Config{
		store: store,
		dir:   dir,
		name:  name,
		cType: ct,
	}
}

func (c *Config) load() error {
	err := c.store.Load(c.dir + c.name)
	eLog(err)
	return err
}

func (c *Config) initialize() {
	if !c.isInit {
		err := c.load()
		if err != nil {
			panic(errors.New("no " + c.name + " in path: " + c.dir))
		}
		c.isInit = true
	}
}

func (c *Config) Variable() StoreVariable {
	c.initialize()
	return c.store
}

// Reload uses to reload configuration file.
func (c *Config) Reload() error {
	return c.load()
}

func (c *Config) Value(name string) Variable {
	return c.Variable().Value(name)
}

func (c *Config) Get(name string) interface{} {
	return c.Variable().Get(name)
}

func (c *Config) GetString(name string) string {
	return c.Variable().GetString(name)
}

func (c *Config) GetFloat(name string) float64 {
	return c.Variable().GetFloat(name)
}

func (c *Config) GetInt(name string) int {
	return c.Variable().GetInt(name)
}

func (c *Config) GetBool(name string) bool {
	return c.Variable().GetBool(name)
}

func (c *Config) Struct(name string, receiver interface{}) {
	ref := reflect.TypeOf(receiver)
	if kind := ref.Kind(); kind != reflect.Ptr {
		eLog(errors.New("the receiver must be a pointer type, kind: " + kind.String()))
		return
	}
	if kind := ref.Elem().Kind(); kind != reflect.Struct && kind != reflect.Map {
		eLog(errors.New("the receiver must be a struct/map pointer type, kind: " + kind.String()))
		return
	}
	c.Variable().Struct(name, receiver)
}
