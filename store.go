package conf

// KeyVariable Various types of values can be obtained according to the key of the string type
type KeyVariable interface {
	Value(string) Variable
	Get(string) interface{}
	GetString(string) string
	GetFloat(string) float64
	GetInt(string) int
	GetBool(string) bool
	Struct(string, interface{})
}

// StoreVariable KeyVariable of that can be store
type StoreVariable interface {
	KeyVariable
	Variable() Variable
	Load(filename string) error
}
