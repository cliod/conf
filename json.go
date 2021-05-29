package conf

import (
	"encoding/json"
	"io/ioutil"
)

type Json struct {
	Store
	data map[string]interface{}
}

func (j *Json) Load(filename string) error {
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	m := make(map[string]interface{})
	err = json.Unmarshal(bs, &m)
	eLog(err)
	j.data = m
	j.StoreVariable = j
	return err
}

func (j *Json) Value(name string) Variable {
	return &Value{j.data[name]}
}
