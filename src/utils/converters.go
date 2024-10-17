package utils

import (
	"reflect"
	"strconv"
)

func ConvertStructToMap(obj interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	val := reflect.ValueOf(obj)
	typ := reflect.TypeOf(obj)

	// Check if the input is a struct
	if val.Kind() != reflect.Struct {
		return result
	}

	// Iterate over struct fields
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Name
		result[fieldName] = field.Interface()
	}
	return result
}

func ConvertMapToStruct(m map[string]interface{}, obj interface{}) {
	val := reflect.ValueOf(obj).Elem()
	typ := reflect.TypeOf(obj).Elem()

	// Iterate over struct fields
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Name
		field.Set(reflect.ValueOf(m[fieldName]))
	}
}

func ConvertIntToString(i int) string {
	return strconv.Itoa(i)
}
