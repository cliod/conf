package conf

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type CType int

func (t CType) String() string {
	name, ok := CTypeNames[int(t)]
	if ok {
		return fmt.Sprint(name)
	}
	return "type" + strconv.Itoa(int(t))
}

const (
	YAML CType = iota
	PROPS
	JSON

	defPath = "conf/"
	defName = "app.yaml"
)

var (
	CTypeNames = map[interface{}]interface{}{
		0:       "YAML",
		"0":     0,
		"YAML":  0,
		1:       "PROPS",
		"1":     1,
		"PROPS": 1,
		2:       "JSON",
		"2":     2,
		"JSON":  2,
	}
)

type Config struct {
	store StoreVariable // variables of configuration file

	dir    string // dir of configuration file, default: conf/
	name   string // name of configuration file, default: app.yaml
	isInit bool   // is initialized
	cType  CType  // type of configuration file, default: yaml

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

func (c *Config) initialize() {
	if !c.isInit {
		err := c.Load(c.dir + c.name)
		if err != nil {
			panic(errors.New("no " + c.name + " in path: " + c.dir))
		}
		c.isInit = true
	}
}

func (c *Config) Variable() StoreVariable {
	return c.store
}

// Load supports reload configuration file
func (c *Config) Load(filename string) error {
	err := c.store.Load(filename)
	eLog(err)
	return err
}

func (c *Config) Value(name string) Variable {
	c.initialize()
	return c.store.Value(name)
}

func (c *Config) GetValue(name string) interface{} {
	c.initialize()
	return c.store.GetValue(name)
}

func (c *Config) GetString(name string) string {
	c.initialize()
	return c.store.GetString(name)
}

func (c *Config) GetFloat(name string) float64 {
	c.initialize()
	return c.store.GetFloat(name)
}

func (c *Config) GetInt(name string) int {
	c.initialize()
	return c.store.GetInt(name)
}

func (c *Config) GetBool(name string) bool {
	c.initialize()
	return c.store.GetBool(name)
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
	c.initialize()
	c.store.Struct(name, receiver)
}
