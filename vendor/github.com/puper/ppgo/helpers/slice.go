package helpers

import (
	"fmt"
	"reflect"
)

func SlicePluck(key string, in, out interface{}) error {
	inVal := reflect.Indirect(reflect.ValueOf(in))
	outVal := reflect.Indirect(reflect.ValueOf(out))
	if inVal.Kind() != reflect.Slice || outVal.Kind() != reflect.Slice {
		return fmt.Errorf("in, out not a slice")
	}
	for i := 0; i < inVal.Len(); i++ {
		outVal.Set(reflect.Append(outVal, reflect.Indirect(inVal.Index(i)).FieldByName(key)))
	}
	return nil
}

func SliceUnique(s interface{}, uf func(i int) interface{}) {
	rv := reflect.ValueOf(s)
	if rv.Kind() == reflect.Interface || rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Slice {
		panic("s not a slice")
	}
	uniqueMap := map[interface{}]bool{}
	j := 0
	for i := 0; i < rv.Len(); i++ {
		rv1 := rv.Index(i)
		var k interface{}
		if uf == nil {
			k = rv1.Interface()
		} else {
			k = uf(i)
		}
		if uniqueMap[k] {
			continue
		} else {
			uniqueMap[k] = true
			rv.Index(j).Set(rv1)
			j++
		}
	}
	rv.SetLen(j)
}

func SliceReverse(s interface{}) {
	rv := reflect.ValueOf(s)
	if rv.Kind() == reflect.Interface || rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Slice {
		panic("arg is not a slice")
	}
	l := rv.Len()
	max := l / 2
	for i := 0; i < max; i++ {
		tmp := rv.Index(i).Interface()
		rv.Index(i).Set(rv.Index(l - i - 1))
		rv.Index(l - i - 1).Set(reflect.ValueOf(tmp))
	}
}
