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

func structToMap(v any) (map[string]interface{}, error) {
	elem := reflect.ValueOf(v)
	if elem.Kind() == reflect.Ptr {
		elem = elem.Elem()
	}
	if elem.Kind() != reflect.Struct {
		return nil, errors.New("cannot convert to map: v is not a struct")
	}

	mappedStruct := make(map[string]interface{})
	for i := 0; i < elem.NumField(); i++ {
		var val interface{}

		f := elem.Type().Field(i)
		key := f.Tag.Get("bencode")
		if key == "" {
			continue
		}

		structField := elem.Field(i)
		if structField.Kind() == reflect.Struct {
			var err error
			val, err = structToMap(structField.Interface())
			if err != nil {
				return nil, fmt.Errorf("cannot convert struct '%s' to map: %w", key, err)
			}
		} else {
			val = structField.Interface()
		}

		mappedStruct[key] = val
	}

	return mappedStruct, nil
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
	if obj.Kind() == reflect.Struct {
		val := reflect.ValueOf(v).Elem()
		return mapDecodeValues(val, &decoded, val)
	}

	if obj.Kind() != reflect.TypeOf(decoded).Kind() {
		return errors.New("unmarshalling failed: v is not the same type as the decoded value")
	}
	reflect.ValueOf(v).Elem().Set(reflect.ValueOf(decoded)) // update by reflection
	return nil
}

func Marshal(v any) ([]byte, error) {
	obj := reflect.TypeOf(v)
	if obj.Kind() != reflect.Ptr {
		return nil, errors.New("marshalling failed: v must be a pointer")
	}

	obj = obj.Elem()
	if obj.Kind() == reflect.Struct {
		var err error
		v, err = structToMap(v)
		if err != nil {
			return nil, fmt.Errorf("marshalling failed: %w", err)
		}
	}

	encoded, err := Encode(v)
	if err != nil {
		return nil, fmt.Errorf("marshalling failed: %w", err)
	}

	return encoded, nil
}
