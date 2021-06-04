package conf

import "sync"

var (
	once   sync.Once
	config *Config
)

func Conf(params ...string) *Config {
	once.Do(func() {
		config = New(params...)
	})
	return config
}

func Reload() error {
	return Conf().Reload()
}

func Get(name string) interface{} {
	return Conf().Get(name)
}

func GetString(name string) string {
	return Conf().GetString(name)
}

func GetFloat(name string) float64 {
	return Conf().GetFloat(name)
}

func GetInt(name string) int {
	return Conf().GetInt(name)
}

func GetBool(name string) bool {
	return Conf().GetBool(name)
}

func ToStruct(name string, receiver interface{}) {
	Conf().Struct(name, receiver)
}
