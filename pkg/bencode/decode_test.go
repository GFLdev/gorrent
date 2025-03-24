package bencode

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

type decodeTestString struct {
	testData[string]
	randLen int
	delim   string
}

type decodeTestInt struct {
	testData[int]
	start string
	end   string
}

type decodeTestList struct {
	testData[[]interface{}]
	invalid bool
	start   string
	end     string
}

type decodeTestDict struct {
	testData[map[string]interface{}]
	invalid bool
	start   string
	end     string
}

func genNestedData(typeCode int, level int) (testData[interface{}], bool) {
	test := testData[interface{}]{}
	invalid := false
	switch typeCode {
	case Str:
		str := genStringDecodeTest()
		if str.randLen <= len(str.data) {
			test.data = str.data[:str.randLen]
		} else {
			test.data = str.data
		}
		test.bCode = str.bCode
		if isStringInvalid(str) {
			invalid = true
		}
	case Int:
		integer := genIntDecodeTest()
		test.data = integer.data
		test.bCode = integer.bCode
		if isIntInvalid(integer) {
			invalid = true
		}
	case List:
		list := genListDecodeTest(level)
		test.data = list.data
		test.bCode = list.bCode
		if isListInvalid(list) {
			invalid = true
		}
	case Dict:
		dict := genDictDecodeTest(level)
		test.data = dict.data
		test.bCode = dict.bCode
		if isDictInvalid(dict) {
			invalid = true
		}
	default:
		invalid = true
	}
	return test, invalid
}

func isListEmpty(bCode string) bool {
	return bCode[0] == 'l' && bCode[1] == 'e'
}

func isDictEmpty(bCode string) bool {
	return bCode[0] == 'd' && bCode[1] == 'e'
}

func isStringInvalid(s decodeTestString) bool {
	return len(s.data) < s.randLen || s.delim != ":"
}

func isIntInvalid(i decodeTestInt) bool {
	return i.start != "i" || i.end != "e" || (len(i.bCode) > 3 && i.bCode[1] == '0')
}

func isListInvalid(l decodeTestList) bool {
	isEmpty := isListEmpty(l.bCode)
	return (l.start != "l" || l.end != "e" || l.invalid) && !isEmpty
}

func isDictInvalid(d decodeTestDict) bool {
	isEmpty := isDictEmpty(d.bCode)
	return (d.start != "d" || d.end != "e" || d.invalid) && !isEmpty
}

func genStringDecodeTest() decodeTestString {
	randStr := gofakeit.Word()
	strLen := len(randStr)
	randLen := gofakeit.IntN(strLen * 2)
	delim := gofakeit.LetterN(1)
	bCode := strconv.Itoa(randLen) + delim + randStr

	return decodeTestString{
		testData: testData[string]{
			data:  randStr,
			bCode: bCode,
		},
		randLen: randLen,
		delim:   delim,
	}
}

func genIntDecodeTest() decodeTestInt {
	randInt := gofakeit.Int()
	start := gofakeit.LetterN(1)
	end := gofakeit.LetterN(1)
	bCode := start + strconv.Itoa(randInt) + end

	return decodeTestInt{
		testData: testData[int]{
			data:  randInt,
			bCode: bCode,
		},
		start: start,
		end:   end,
	}
}

func genListDecodeTest(level int) decodeTestList {
	listLen := gofakeit.IntN(10) // max of 10 items
	data := make([]interface{}, listLen)
	bCodes := make([]string, listLen)
	invalid := false

	for i := range listLen {
		var test testData[interface{}]
		randType := gofakeit.IntN(4)
		if level >= 5 { // max of 5 nested levels
			randType = gofakeit.IntN(2) // do not select list or dict
		}
		test, invalid = genNestedData(randType, level+1)
		data[i] = test.data
		bCodes[i] = test.bCode
	}

	start := gofakeit.LetterN(1)
	end := gofakeit.LetterN(1)
	bCode := start
	for i := range bCodes {
		bCode += bCodes[i]
	}
	bCode += end

	if isListEmpty(bCode) {
		data = make([]any, 0)
	}

	return decodeTestList{
		testData: testData[[]interface{}]{
			data:  data,
			bCode: bCode,
		},
		invalid: invalid,
		start:   start,
		end:     end,
	}
}

func genDictDecodeTest(level int) decodeTestDict {
	dictLen := gofakeit.IntN(10) // max of 10 items
	data := make(map[string]interface{}, dictLen)
	bCodes := make([]string, dictLen)
	invalid := false

	for i := range dictLen {
		var testKey, testVal testData[interface{}]
		var key, val any
		randKeyType := gofakeit.IntN(4)
		randValType := gofakeit.IntN(4)
		if level >= 5 { // max of 5 nested levels
			randKeyType = gofakeit.IntN(2) // do not select list or dict
			randValType = gofakeit.IntN(2)
		}

		var invalidTemp bool
		testKey, invalidTemp = genNestedData(randKeyType, level+1)
		for data[fmt.Sprintf("%v", testKey.data)] != nil { // regenerate if key already exists
			testKey, invalidTemp = genNestedData(randKeyType, level+1)
		}
		invalid = invalidTemp

		testVal, invalid = genNestedData(randValType, level+1)
		key = testKey.data
		val = testVal.data
		if reflect.TypeOf(key).Kind() != reflect.String {
			invalid = true
		}
		bCodes[i] = testKey.bCode + testVal.bCode
		data[fmt.Sprintf("%v", key)] = val
	}

	start := gofakeit.LetterN(1)
	end := gofakeit.LetterN(1)
	bCode := start
	for i := range bCodes {
		bCode += bCodes[i]
	}
	bCode += end

	if isDictEmpty(bCode) {
		data = make(map[string]interface{})
	}

	return decodeTestDict{
		testData: testData[map[string]interface{}]{
			data:  data,
			bCode: bCode,
		},
		invalid: invalid,
		start:   start,
		end:     end,
	}
}

func TestDecodeString(t *testing.T) {
	t.Parallel()
	for range *SimNumbers {
		seed := gofakeit.Int64()
		if gofakeit.Seed(seed) != nil {
			continue
		}

		test := genStringDecodeTest()
		decoded, err := Decode([]byte(test.bCode))
		if isStringInvalid(test) {
			assert.Error(t, err, FormatInfo(seed, test.bCode))
		} else {
			if assert.NoError(t, err, FormatInfo(seed, test.bCode)) {
				assert.Equal(t, test.data[:test.randLen], decoded.(string), FormatInfo(seed, test.bCode))
			}
		}
	}
}

func TestDecodeInt(t *testing.T) {
	t.Parallel()
	for range *SimNumbers {
		seed := gofakeit.Int64()
		if gofakeit.Seed(seed) != nil {
			continue
		}

		test := genIntDecodeTest()
		decoded, err := Decode([]byte(test.bCode))
		if isIntInvalid(test) {
			assert.Error(t, err, FormatInfo(seed, test.bCode))
		} else {
			if assert.NoError(t, err, FormatInfo(seed, test.bCode)) {
				assert.Equal(t, test.data, decoded.(int), FormatInfo(seed, test.bCode))
			}
		}
	}
}

func TestDecodeList(t *testing.T) {
	t.Parallel()
	for range *SimNumbers {
		seed := gofakeit.Int64()
		if gofakeit.Seed(seed) != nil {
			continue
		}

		test := genListDecodeTest(0)

		// Start rune may be valid
		if test.start == "d" || test.start == "i" {
			continue
		}

		decoded, err := Decode([]byte(test.bCode))
		if isListInvalid(test) {
			if test.bCode[1] == byte('e') {
				continue
			}
			assert.Error(t, err, FormatInfo(seed, test.bCode))
		} else {
			if assert.NoError(t, err, FormatInfo(seed, test.bCode)) {
				if len(test.data) == 0 {
					assert.Equal(t, make([]interface{}, 0), decoded, FormatInfo(seed, test.bCode))
				} else {
					assert.Equal(t, test.data, decoded, FormatInfo(seed, test.bCode))
				}
			}
		}
	}
}

func TestDecodeDict(t *testing.T) {
	t.Parallel()
	for range *SimNumbers {
		seed := gofakeit.Int64()
		if gofakeit.Seed(seed) != nil {
			continue
		}

		test := genDictDecodeTest(0)

		// Start rune may be valid
		if test.start == "l" || test.start == "i" {
			continue
		}

		decoded, err := Decode([]byte(test.bCode))
		if isDictInvalid(test) {
			if test.bCode[1] == byte('e') {
				continue
			}
			assert.Error(t, err, FormatInfo(seed, test.bCode))
		} else {
			if assert.NoError(t, err, FormatInfo(seed, test.bCode)) {
				if len(test.data) == 0 {
					assert.Equal(t, make(map[string]interface{}), decoded, FormatInfo(seed, test.bCode))
				} else {
					assert.Equal(t, test.data, decoded, FormatInfo(seed, test.bCode))
				}
			}
		}
	}
}
