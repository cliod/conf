package conf

import (
	"fmt"
	"strconv"
	"strings"
)

type CType int

func (t CType) String() string {
	name, ok := CTypeMapper[int(t)]
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
	CTypeMapper = map[interface{}]interface{}{
		YAML:    "YAML",
		JSON:    "JSON",
		PROPS:   "PROPS",
		0:       "YAML",
		1:       "PROPS",
		2:       "JSON",
		"0":     0,
		"1":     1,
		"2":     2,
		"YAML":  0,
		"PROPS": 1,
		"JSON":  2,
	}
)

func getCType(filename string) CType {
	extName := filename[strings.LastIndex(filename, ".")+1:]
	code, ok := CTypeMapper[strings.ToUpper(extName)].(int)
	if ok {
		return CType(code)
	}
	return YAML
}
