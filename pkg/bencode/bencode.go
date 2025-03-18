/*
Package bencode provides encoding and decoding of data in the Bencode format.
*/
package bencode

import (
	"errors"
	"fmt"
	"reflect"
)

// mapToStruct set decoded bencode data (map) to a struct by matching struct fields with "bencode" tags.
func mapToStruct(val reflect.Value, decodedMap *interface{}, v any) error {
	for i := 0; i < val.NumField(); i++ {
		f := val.Type().Field(i)

		// Get bencode field and use it as key in decodedMap to find its value
		bencodeField := f.Tag.Get("bencode")
		if bencodeField == "" {
			continue
		}
		decodedVal, ok := (*decodedMap).(map[string]interface{})[bencodeField]
		if !ok {
			continue
		}

		structField := val.Field(i)
		if structField.Kind() == reflect.Struct { // recursively convert map fields if value is struct
			err := mapToStruct(structField, &decodedVal, v)
			if err != nil {
				return err
			}
			continue
		}
		if structField.CanSet() { // set value to given field
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

// structToMap converts a struct to a map using the "bencode" tags as keys.
func structToMap(v any) (map[string]interface{}, error) {
	// Check if v is a struct
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

		// Get bencode tag to use as key
		key := f.Tag.Get("bencode")
		if key == "" {
			continue
		}

		// Set values
		structField := elem.Field(i)
		if structField.Kind() == reflect.Struct { // recursively convert nested structs to map
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

// Unmarshal decodes bencoded data into the structure or variable provided by v, which must be a pointer.
func Unmarshal(data []byte, v any) error {
	// Decode bencode
	decoded, err := Decode(data)
	if err != nil {
		return err
	}

	// Check if v is a pointer
	obj := reflect.TypeOf(v)
	if obj.Kind() != reflect.Ptr {
		return errors.New("unmarshalling failed: v must be a pointer")
	}
	obj = obj.Elem()

	// If struct, set its fields value with decoded bencode, by "bencode" tag
	if obj.Kind() == reflect.Struct {
		val := reflect.ValueOf(v).Elem()
		return mapToStruct(val, &decoded, val)
	}

	// Otherwise, set decoded value
	if obj.Kind() != reflect.TypeOf(decoded).Kind() {
		return errors.New("unmarshalling failed: v is not the same type as the decoded value")
	}
	reflect.ValueOf(v).Elem().Set(reflect.ValueOf(decoded)) // update by reflection
	return nil
}

// Marshal encodes the given value into a byte slice using with bencode format.
func Marshal(v any) ([]byte, error) {
	// Check if v is a pointer
	obj := reflect.TypeOf(v)
	if obj.Kind() == reflect.Ptr {
		obj = obj.Elem()
	}

	// If v is a struct, map its values with "bencode" tag as keys
	var val any
	if obj.Kind() == reflect.Struct {
		var err error
		val, err = structToMap(v)
		if err != nil {
			return nil, fmt.Errorf("marshalling failed: %w", err)
		}
	} else {
		val = v
	}

	// Encode
	encoded, err := Encode(val)
	if err != nil {
		return nil, fmt.Errorf("marshalling failed: %w", err)
	}

	return encoded, nil
}
