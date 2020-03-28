package helpers

import (
	"errors"
	"reflect"
)

func FillAttr(data interface{},
	setAttrFunc func(interface{}, interface{}),
	getAttrFunc func(interface{}) interface{},
	getDataFunc func(interface{}) (interface{}, error)) error {
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	if dataVal.Kind() != reflect.Slice {
		return errors.New("data not a slice")
	}
	attrs := make([]interface{}, 0, dataVal.Len())
	dataMap := make(map[interface{}][]int)
	for i := 0; i < dataVal.Len(); i++ {
		v := dataVal.Index(i).Interface()
		attr := getAttrFunc(v)
		attrs = append(attrs, attr)
		dataMap[attr] = append(dataMap[attr], i)
	}
	vals, err := getDataFunc(attrs)
	if err != nil {
		return err
	}
	valsVal := reflect.Indirect(reflect.ValueOf(vals))
	for _, k := range valsVal.MapKeys() {
		if dataMap[k.Interface()] != nil {
			for _, i := range dataMap[k.Interface()] {
				setAttrFunc(dataVal.Index(i).Interface(), valsVal.MapIndex(k).Interface())
			}
		}
	}
	return nil
}

func PickKeyValue(in, out interface{}, key, value string) error {
	inVal := reflect.Indirect(reflect.ValueOf(in))
	if inVal.Kind() != reflect.Slice {
		return errors.New("in not a slice")
	}
	outVal := reflect.Indirect(reflect.ValueOf(out))
	if outVal.Kind() != reflect.Map || outVal.IsNil() {
		return errors.New("out not a map")
	}
	for i := 0; i < inVal.Len(); i++ {
		v := reflect.Indirect(inVal.Index(i))
		outVal.SetMapIndex(v.FieldByName(key), v.FieldByName(value))
	}
	return nil
}
