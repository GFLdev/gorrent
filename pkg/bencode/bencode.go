package bencode

import (
	"errors"
	"fmt"
	"reflect"
)

func mapDecodeValues(val reflect.Value, decodedMap *interface{}, v any) error {
	for i := 0; i < val.NumField(); i++ {
		f := val.Type().Field(i)

		bencodeField := f.Tag.Get("bencode")
		if bencodeField == "" {
			continue
		}

		decodedVal, ok := (*decodedMap).(map[string]interface{})[bencodeField]
		if !ok {
			continue
		}

		structField := val.Field(i)
		if structField.Kind() == reflect.Struct {
			err := mapDecodeValues(structField, &decodedVal, v)
			if err != nil {
				return err
			}
			continue
		}
		if structField.CanSet() {
			if !reflect.TypeOf(decodedVal).AssignableTo(structField.Type()) {
				return fmt.Errorf(
					"cannot assign decoded value to field '%s': expected %s, got %s",
					bencodeField,
					structField.Type(),
					reflect.TypeOf(decodedVal),
				)
			}
			structField.Set(reflect.ValueOf(decodedVal))
		}
	}
	return nil
}

func Unmarshal(data []byte, v any) error {
	decoded, err := Decode(data)
	if err != nil {
		return err
	}

	obj := reflect.TypeOf(v)
	if obj.Kind() != reflect.Ptr {
		return errors.New("unmarshalling failed: v must be a pointer")
	}
	obj = obj.Elem()
	if obj.Kind() != reflect.Struct {
		if obj.Kind() != reflect.TypeOf(decoded).Kind() {
			return errors.New("unmarshalling failed: v is not the same type as the decoded value")
		}
		reflect.ValueOf(v).Elem().Set(reflect.ValueOf(decoded)) // update by reflection
		return nil
	}

	val := reflect.ValueOf(v).Elem()
	return mapDecodeValues(val, &decoded, val)
}
