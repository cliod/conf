package conf

type StoreVariable interface {
	Load(filename string) error
	Value(string) Variable
	GetValue(string) interface{}
	GetString(string) string
	GetFloat(string) float64
	GetInt(string) int
	GetBool(string) bool
	Struct(string, interface{})
}

type Store struct {
	StoreVariable
}

func (s *Store) GetValue(name string) interface{} {
	return s.Value(name).Value()
}

func (s *Store) GetString(name string) string {
	return s.Value(name).String()
}

func (s *Store) GetFloat(name string) float64 {
	return s.Value(name).Float()
}

func (s *Store) GetInt(name string) int {
	return s.Value(name).Int()
}

func (s *Store) GetBool(name string) bool {
	return s.Value(name).Bool()
}

func (s *Store) Struct(name string, receiver interface{}) {
	s.Value(name).Struct(receiver)
}
