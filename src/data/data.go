package data

type Data struct {
	Key   string
	Value interface{}
}

func NewData(key string, value interface{}) Data {
	return Data{
		Key:   key,
		Value: value,
	}
}
