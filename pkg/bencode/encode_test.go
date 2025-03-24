package bencode

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"sort"
	"strconv"
	"testing"
)

func switchGenType(typeCode int, level int) testData[interface{}] {
	var test testData[interface{}]
	switch typeCode {
	case Str:
		val := genStringEncodeTest()
		test.data, test.bCode = val.data, val.bCode
	case Int:
		val := genIntEncodeTest()
		test.data, test.bCode = val.data, val.bCode
	case List:
		val := genListEncodeTest(level)
		test.data, test.bCode = val.data, val.bCode
	case Dict:
		val := genDictEncodeTest(level)
		test.data, test.bCode = val.data, val.bCode
	}
	return test
}

func genStringEncodeTest() testData[string] {
	test := gofakeit.Word()
	return testData[string]{
		data:  test,
		bCode: strconv.Itoa(len(test)) + ":" + test,
	}
}

func genIntEncodeTest() testData[int] {
	test := gofakeit.Int()
	return testData[int]{
		data:  test,
		bCode: "i" + strconv.Itoa(test) + "e",
	}
}

func genListEncodeTest(level int) testData[[]interface{}] {
	listLen := gofakeit.IntN(10) // max of 10 items
	list := make([]interface{}, listLen)
	bCodes := make([]string, listLen)

	// Generate data
	for i := range listLen {
		randType := gofakeit.IntN(4)
		if level >= 5 { // max of 5 nested levels
			randType = gofakeit.IntN(2) // do not select list or dict
		}
		item := switchGenType(randType, level+1)
		list[i], bCodes[i] = item.data, item.bCode
	}

	// Generate expected data
	expected := "l"
	for _, bCode := range bCodes {
		expected += bCode
	}
	expected += "e"
	return testData[[]interface{}]{
		data:  list,
		bCode: expected,
	}
}

func genDictEncodeTest(level int) testData[map[string]interface{}] {
	dictLen := gofakeit.IntN(10) // max of 10 items
	dict := make(map[string]interface{}, dictLen)
	bCodes := make(map[string]string, dictLen)
	keys := make([]string, dictLen)

	// Generate data
	for i := range dictLen {
		randType := gofakeit.IntN(4)
		if level >= 5 { // max of 5 nested levels
			randType = gofakeit.IntN(2) // do not select list or dict
		}
		key := gofakeit.Word()
		for dict[key] != nil { // regenerate if key already exists
			key = gofakeit.Word()
		}
		keys[i] = key
		item := switchGenType(randType, level+1)
		dict[key], bCodes[key] = item.data, item.bCode
	}

	// Generate expected data
	sort.Strings(keys)
	expected := "d"
	for _, key := range keys {
		keyBCode, err := Encode(key)
		if err != nil {
			delete(dict, key)
			delete(bCodes, key)
			continue
		}
		expected += string(keyBCode)
		expected += bCodes[key]
	}
	expected += "e"
	return testData[map[string]interface{}]{
		data:  dict,
		bCode: expected,
	}
}

func TestEncodeString(t *testing.T) {
	t.Parallel()
	for range *SimNumbers {
		seed := gofakeit.Int64()
		if gofakeit.Seed(seed) != nil {
			continue
		}

		mock := switchGenType(Str, 0)
		test, expected := mock.data.(string), mock.bCode
		bCode, err := Encode(test)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, string(bCode), FormatSeed(seed))
		}
	}
}

func TestEncodeInt(t *testing.T) {
	t.Parallel()
	for range *SimNumbers {
		seed := gofakeit.Int64()
		if gofakeit.Seed(seed) != nil {
			continue
		}

		mock := switchGenType(Int, 0)
		test, expected := mock.data.(int), mock.bCode
		bCode, err := Encode(test)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, string(bCode), FormatSeed(seed))
		}
	}
}

func TestEncodeList(t *testing.T) {
	t.Parallel()
	for range *SimNumbers {
		seed := gofakeit.Int64()
		if gofakeit.Seed(seed) != nil {
			continue
		}

		mock := switchGenType(List, 0)
		test, expected := mock.data.([]interface{}), mock.bCode
		bCode, err := Encode(test)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, string(bCode), FormatSeed(seed))
		}
	}
}

func TestEncodeDict(t *testing.T) {
	t.Parallel()
	for range *SimNumbers {
		seed := gofakeit.Int64()
		if gofakeit.Seed(seed) != nil {
			continue
		}

		mock := switchGenType(Dict, 0)
		test, expected := mock.data.(map[string]interface{}), mock.bCode
		bCode, err := Encode(test)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, string(bCode), FormatSeed(seed))
		}
	}
}
