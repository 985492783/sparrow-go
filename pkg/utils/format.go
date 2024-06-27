package utils

import (
	"encoding/json"
	"reflect"
)

type Type string

const (
	String Type = "string"
	Json   Type = "json"
	Float  Type = "float"
	Bool   Type = "bool"
	Int    Type = "int"
)

func IsTypeOf(data any, tp string) bool {
	t := reflect.TypeOf(data)
	switch Type(tp) {
	case String:
		return t.Kind() == reflect.String
	case Json:
		return t.Kind() == reflect.String && isValidJSON(data.(string))
	case Float:
		return t.Kind() == reflect.Float32 || t.Kind() == reflect.Float64
	case Bool:
		return t.Kind() == reflect.Bool
	case Int:
		return t.Kind() == reflect.Int || t.Kind() == reflect.Int64
	}
	return false
}

func isValidJSON(s string) bool {
	var js interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}
