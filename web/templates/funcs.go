package templates

type kvReturn struct {
	Key   interface{}
	Value interface{}
}

func keyValue(key, value any) kvReturn {
	return kvReturn{
		Key:   key,
		Value: value,
	}
}
