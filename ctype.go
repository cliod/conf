package conf

import (
	"fmt"
	"strconv"
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
