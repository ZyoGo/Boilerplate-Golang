package gin

import (
	"fmt"
)

var sensitiveFields = map[string]bool{
	"password": true,
}

type Document struct{}

type MapType map[string]interface{}
type ArrayType []interface{}

// ProcessMap processes a map and obfuscates sensitive fields, truncates large arrays, and recursively processes nested structures.
func (doc *Document) ProcessMap(data MapType) MapType {
	for key, value := range data {
		if value == nil {
			continue
		}
		switch v := value.(type) {
		case map[string]interface{}:
			data[key] = doc.ProcessMap(v)
		case []interface{}:
			if len(v) > 10 {
				data[key] = fmt.Sprintf(`{"count" : "%d"}`, len(v))
			} else {
				data[key] = doc.ProcessArray(v)
			}
		default:
			if sensitiveFields[key] {
				data[key] = "*******"
			} else {
				data[key] = v
			}
		}
	}
	return data
}

// ProcessArray processes an array and recursively processes nested maps or arrays.
func (doc *Document) ProcessArray(data ArrayType) ArrayType {
	for i, value := range data {
		switch v := value.(type) {
		case map[string]interface{}:
			data[i] = doc.ProcessMap(v)
		case []interface{}:
			data[i] = doc.ProcessArray(v)
		default:
			data[i] = v
		}
	}
	return data
}
