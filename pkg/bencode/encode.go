package bencode

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"sync"
)

func formatErrorSlice(errs []error) string {
	if len(errs) == 1 {
		return errs[0].Error()
	}
	str := "\n"
	for i, err := range errs {
		str += fmt.Sprintf("%d: %s\n", i+1, err.Error())
	}
	str = str[:len(str)-1]
	return str
}

func insertSortKVS(kv []byte, kvs *[][]byte) {
	// Binary search
	var idx int
	for l := 0; l < len(*kvs); l++ {
		if len((*kvs)[l]) == 0 {
			idx = l // default: insert at the end
			break
		}
	}
	l := 0
	r := len(*kvs) - 1
	for l <= r {
		m := (l + r) / 2
		if bytes.Compare(kv, (*kvs)[m]) == -1 { // less than
			r = m - 1
		} else if bytes.Compare(kv, (*kvs)[m]) == 1 { // greater than
			l = m + 1
		} else {
			idx = m
			break
		}
	}

	// Insertion sort with previously sorted slice
	// Copy to buffer
	buf := make([][]byte, len((*kvs)[idx:]))
	copy(buf, (*kvs)[idx:])

	// Copy to slice
	(*kvs)[idx] = kv
	copy((*kvs)[idx+1:], buf)
	return
}

func encodeString(v interface{}) ([]byte, error) {
	s, ok := v.(string)
	if !ok {
		return nil, fmt.Errorf("cannot encode string: expected string, got %s", reflect.TypeOf(v))
	}
	strLen := len(s)
	strLenStr := strconv.Itoa(strLen)

	// Encoding
	encoded := make([]byte, len(strLenStr)+1+strLen)
	n := copy(encoded[:], strLenStr)
	n += copy(encoded[n:], ":")
	n += copy(encoded[n:], s)
	return encoded, nil
}

func encodeInt(v interface{}) ([]byte, error) {
	i, ok := v.(int)
	if !ok {
		return nil, fmt.Errorf("cannot encode int: expected int, got %s", reflect.TypeOf(v))
	}
	intStr := strconv.Itoa(i)

	// Encoding
	encoded := make([]byte, len(intStr)+2)
	n := copy(encoded[:], "i")
	n += copy(encoded[n:], intStr)
	n += copy(encoded[n:], "e")
	return encoded, nil
}

func encodeList(v interface{}) ([]byte, error) {
	var wg sync.WaitGroup
	var mux sync.Mutex
	var errs []error

	// Type check and create []interface
	if reflect.ValueOf(v).Kind() != reflect.Slice || reflect.ValueOf(v).Kind() == reflect.Array {
		return nil, fmt.Errorf("cannot encode list: expected slice or array, got %s", reflect.TypeOf(v))
	}
	s := make([]interface{}, reflect.ValueOf(v).Len())
	for i := 0; i < reflect.ValueOf(v).Len(); i++ {
		s[i] = reflect.ValueOf(v).Index(i).Interface()
	}

	// Buffer and encoding each values concurrently
	elem := make([][]byte, len(s))
	totalBytes := 2
	wg.Add(len(s))
	for i, val := range s {
		go func(idx int, val interface{}) {
			defer wg.Done()
			var err error

			elem[idx], err = Encode(val)
			if err != nil {
				mux.Lock()
				errs = append(errs, err)
				mux.Unlock()
				return
			}

			mux.Lock()
			totalBytes += len(elem[idx])
			mux.Unlock()
		}(i, val)
	}
	wg.Wait()
	if len(errs) > 0 {
		return nil, fmt.Errorf("cannot encode list: %s", formatErrorSlice(errs))
	}

	// Encoding the list
	encoded := make([]byte, totalBytes)
	i := copy(encoded[:], "l")
	for _, val := range elem {
		i += copy(encoded[i:], val)
	}
	i += copy(encoded[i:], "e")
	return encoded, nil
}

func encodeMap(v interface{}) ([]byte, error) {
	var wg sync.WaitGroup
	var mux sync.Mutex
	var errs []error

	// Type check and create map[string]interface
	if reflect.ValueOf(v).Kind() != reflect.Map {
		return nil, fmt.Errorf("cannot encode map: expected map, got %s", reflect.TypeOf(v))
	}
	m := make(map[string]interface{}, reflect.ValueOf(v).Len())
	for _, k := range reflect.ValueOf(v).MapKeys() {
		m[k.String()] = reflect.ValueOf(v).MapIndex(k).Interface()
	}

	// Buffer and encoding each key-pair values concurrently
	elem := make([][]byte, len(m))
	totalBytes := 2
	i := -1
	wg.Add(len(m))
	for k, v := range m {
		i++
		go func(key string, val interface{}, idx int) {
			defer wg.Done()

			keyStr, err := encodeString(key)
			if err != nil {
				mux.Lock()
				errs = append(errs, err)
				mux.Unlock()
				return
			}

			valStr, err := Encode(val)
			if err != nil {
				mux.Lock()
				errs = append(errs, err)
				mux.Unlock()
				return
			}

			mux.Lock()
			insertSortKVS(append(keyStr, valStr...), &elem)
			totalBytes += len(keyStr) + len(valStr)
			mux.Unlock()
		}(k, v, i)
	}
	wg.Wait()

	if len(errs) > 0 {
		return nil, fmt.Errorf("cannot encode map: %s", formatErrorSlice(errs))
	}

	// Encoding the map
	encoded := make([]byte, totalBytes)
	i = copy(encoded[:], "d")
	for _, val := range elem {
		i += copy(encoded[i:], val)
	}
	i += copy(encoded[i:], "e")
	return encoded, nil
}

func Encode(v interface{}) ([]byte, error) {
	if v == nil {
		return nil, nil
	}

	switch reflect.TypeOf(v).Kind() {
	case reflect.String:
		return encodeString(v)
	case reflect.Int:
		return encodeInt(v)
	case reflect.Slice:
		return encodeList(v)
	case reflect.Map:
		return encodeMap(v)
	default:
		return nil, fmt.Errorf("cannot encode type %s", reflect.TypeOf(v))
	}
}
