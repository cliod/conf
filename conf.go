package conf

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Config struct {
	variable StoreVariable // variables of configuration file

	dir    string // dir of configuration file, default: conf/
	name   string // name of configuration file, default: app.yaml
	isInit bool   // is initialized
	cType  CType  // type of configuration file, default: yaml

	store    StoreVariable   // variables of main config
	includes []StoreVariable // expanded config (supplement)
	profile  StoreVariable   // personal config (covering)
}

func New(params ...interface{}) *Config {
	var (
		confDir  string
		confName string
		cType    CType
		store    StoreVariable
	)
	if len(params) > 0 {
		confDir, _ = params[0].(string)
		if len(params) > 1 {
			confName, _ = params[1].(string)
			if len(params) > 2 {
				ct, ok := CTypeMapper[params[2]]
				if ok {
					if str, ok := ct.(string); ok {
						ct, ok = CTypeMapper[str]
					}
					if cType, ok = ct.(CType); !ok {
						cType = CType(ct.(int))
					}
				}
				if len(params) > 3 {
					store = params[3].(StoreVariable)
				}
			} else {
				cType = getCType(confName)
			}
		}
	}
	if strings.Trim(confDir, " ") == "" || confDir[0] != '/' || confDir[1] != ':' {
		confDir = rootPath() + defPath
	}
	if strings.Trim(confName, " ") == "" {
		confName = defName
	}
	return newConf(confDir, confName, cType, store)
}

func newConf(dir, name string, ct CType, store StoreVariable) *Config {
	if store == nil {
		store = newStore(ct)
	}
	return &Config{
		store: store,
		dir:   dir,
		name:  name,
		cType: ct,
	}
}

func newStore(ct CType) StoreVariable {
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
	return store
}

func (c *Config) loadIncludes() error {
	// get all includes
	// new store
	// load file
	var paths []string
	rootInclude := c.store.Get("include")
	appInclude := c.store.Get("app.include")
	paths = append(paths, c.toSlice(rootInclude)...)
	paths = append(paths, c.toSlice(appInclude)...)
	for _, path := range paths {
		store := newStore(c.cType)
		c.includes = append(c.includes, store)
		err := store.Load(c.dir + path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) toSlice(path interface{}) []string {
	switch value := path.(type) {
	case []interface{}:
		var s []string
		for _, val := range value {
			s = append(s, newVal(val).String())
		}
		return s
	default:
		return []string{newVal(value).String()}
	}
}

func (c *Config) loadProfile() error {
	appProfile := c.store.Get("app.profile")
	c.profile = newStore(c.cType)
	return c.profile.Load(fmt.Sprintf("%s%s-%s.%s", c.dir, c.name[:strings.LastIndex(c.name, ".")], appProfile, strings.ToLower(c.cType.String())))
}

func (c *Config) load() {
	err := c.store.Load(c.dir + c.name)
	if err != nil {
		panic(errors.New("no " + c.name + " in path: " + c.dir))
	}
	err = c.loadIncludes()
	if err != nil {
		wLog(err)
	}
	err = c.loadProfile()
	if err != nil {
		wLog(err)
	}
}

func (c *Config) expand() {
	m := newMixture(c.store, c.includes...)
	m.replace(c.profile)
	c.variable = m
}

func (c *Config) initialize() {
	c.load()
	c.expand()
	c.isInit = true
}

func (c *Config) Variable() StoreVariable {
	if !c.isInit {
		c.initialize()
	}
	return c.variable
}

// Reload uses to reload configuration file.
func (c *Config) Reload() error {
	defer func() {
		err := c.loadIncludes()
		if err != nil {
			wLog(err)
		}
		err = c.loadProfile()
		if err != nil {
			wLog(err)
		}
	}()
	return c.store.Load(c.dir + c.name)
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
