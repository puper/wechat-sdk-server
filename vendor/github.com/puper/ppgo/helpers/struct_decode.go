package helpers

import (
	"reflect"

	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
)

func StructDecode(input, output interface{}, tagName string) error {
	if reflect.Indirect(reflect.ValueOf(input)).Kind() == reflect.Struct {
		structs.DefaultTagName = tagName
		input = structs.Map(input)
	}
	config := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           output,
		WeaklyTypedInput: true,
		TagName:          tagName,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	return decoder.Decode(input)
}
