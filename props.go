package conf

import (
	"github.com/magiconair/properties"
)

type Props struct {
	Store
	props *properties.Properties
}

func (p *Props) Load(filename string) error {
	props := properties.MustLoadFile(filename, properties.UTF8)
	p.props = props
	p.StoreVariable = p
	return nil
}

func (p *Props) Value(name string) Variable {
	return &Value{p.props.GetString(name, "")}
}
