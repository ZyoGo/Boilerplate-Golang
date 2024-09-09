package gin

import (
	"fmt"
	"reflect"
)

var credFields = map[string]bool{
	"password": true,
}

type Document struct{}

type MapType map[string]interface{}
type ArrayType []interface{}

func (doc *Document) throughMap(docMap MapType) MapType {
	for k, v := range docMap {
		if v == nil {
			continue
		}
		vt := reflect.TypeOf(v)
		switch vt.Kind() {
		case reflect.Map:
			if mv, ok := v.(map[string]interface{}); ok {
				docMap[k] = doc.throughMap(mv)
			} else {
				panic("error.")
			}
		case reflect.Array, reflect.Slice:
			if mv, ok := v.([]interface{}); ok {
				if len(mv) > 10 {
					docMap[k] = fmt.Sprintf(`{"count" : "%d"}`, len(mv))
				} else {
					docMap[k] = doc.throughArray(mv)
				}
			} else {
				panic("error.")
			}
		default:
			if credFields[k] {
				docMap[k] = "*******"
			} else {
				docMap[k] = v
			}
		}
	}
	return docMap
}

func (doc *Document) throughArray(arrayType ArrayType) ArrayType {
	for k, v := range arrayType {
		vt := reflect.TypeOf(v)
		switch vt.Kind() {
		case reflect.Map:
			if mv, ok := v.(map[string]interface{}); ok {
				arrayType[k] = doc.throughMap(mv)
			} else {
				panic("error.")
			}
		case reflect.Array, reflect.Slice:
			if mv, ok := v.([]interface{}); ok {
				arrayType[k] = doc.throughArray(mv)
			} else {
				panic("error.")
			}
		default:
			arrayType[k] = v
		}
	}
	return arrayType
}
