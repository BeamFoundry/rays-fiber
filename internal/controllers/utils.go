package controllers

import (
	"reflect"

	. "github.com/BeamFoundry/rays/pkg/core"
	. "github.com/BeamFoundry/rays/pkg/util"
)

func CopyPrivateToPublic(source, destination any) {
	st := reflect.TypeOf(source).Elem()
	dt := reflect.TypeOf(destination).Elem()
	for i := range st.NumField() {
		capitalized := CaseUtils(st.Field(i).Name).Capitalize().String()
		if dst, ok := dt.FieldByName(capitalized); ok {
			field := reflect.ValueOf(source).Elem().Field(i)
			if field.Kind() == reflect.Pointer {
				field = field.Elem()
			}
			if field.Type() != dst.Type {
				continue
			}
			value := UnprotectValue(field)
			reflect.ValueOf(destination).Elem().FieldByName(capitalized).Set(value)
		}
	}
}

func toAny[T any](s []T) []any {
	anySlice := make([]any, len(s))
	for i, v := range s {
		anySlice[i] = v
	}
	return anySlice
}
