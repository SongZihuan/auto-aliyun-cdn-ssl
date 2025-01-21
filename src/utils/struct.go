package utils

import "reflect"

func HasFieldByReflect(typ reflect.Type, fieldName string) bool {
	for i := 0; i < typ.NumField(); i++ {
		if typ.Field(i).Name == fieldName {
			return true
		}
	}
	return false
}
