package conf

type KindVariable interface {
	Value(string) Variable
	Get(string) interface{}
	GetString(string) string
	GetFloat(string) float64
	GetInt(string) int
	GetBool(string) bool
	Struct(string, interface{})
}

type StoreVariable interface {
	KindVariable
	Load(filename string) error
}
