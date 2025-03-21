package bencode

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

const SimNumbers = 100000

const (
	str int = iota
	integer
	list
	dict
)

type testString struct {
	randStr string
	randLen int
	delim   string
	bCode   string
}

type testInt struct {
	randInt int
	start   string
	end     string
	bCode   string
}

type testList struct {
	data      []any
	bCodeData string
	invalid   bool
	start     string
	end       string
	bCode     string
}

type testDict struct {
	data      map[string]any
	bCodeData string
	invalid   bool
	start     string
	end       string
	bCode     string
}

func printSeed(seed int64) string {
	return "seed: " + strconv.Itoa(int(seed))
}

func switchNestedData(randType int, data *any, bCode *string, level int) bool {
	invalid := false
	switch randType {
	case str:
		test := genTestString()
		if test.randLen <= len(test.randStr) {
			*data = test.randStr[:test.randLen]
		} else {
			*data = test.randStr
		}
		*bCode += test.bCode
		if isStringInvalid(test) {
			invalid = true
		}
	case integer:
		test := genTestInt()
		*data = test.randInt
		*bCode += test.bCode
		if isIntInvalid(test) {
			invalid = true
		}
	case list:
		test := genTestList(level + 1)
		*data = test.data
		*bCode += test.bCode
		if isListInvalid(test) {
			invalid = true
		}
	case dict:
		test := genTestDict(level + 1)
		*data = test.data
		*bCode += test.bCode
		if isDictInvalid(test) {
			invalid = true
		}
	default:
		invalid = true
	}
	return invalid
}

func isStringInvalid(s testString) bool {
	return len(s.randStr) < s.randLen || s.delim != ":"
}

func isIntInvalid(i testInt) bool {
	return i.start != "i" || i.end != "e" || i.randInt < 0
}

func isListInvalid(l testList) bool {
	return l.start != "l" || l.end != "e" || l.invalid
}

func isDictInvalid(d testDict) bool {
	return d.start != "d" || d.end != "e" || d.invalid
}

func genTestString() testString {
	randStr := gofakeit.Word()
	strLen := len(randStr)
	randLen := gofakeit.IntN(strLen * 2)
	delim := gofakeit.LetterN(1)
	bCode := strconv.Itoa(randLen) + delim + randStr

	return testString{
		randStr: randStr,
		randLen: randLen,
		delim:   delim,
		bCode:   bCode,
	}
}

func genTestInt() testInt {
	randInt := gofakeit.Int()
	start := gofakeit.LetterN(1)
	end := gofakeit.LetterN(1)
	bCode := start + strconv.Itoa(randInt) + end

	return testInt{
		randInt: randInt,
		start:   start,
		end:     end,
		bCode:   bCode,
	}
}

func genTestList(level int) testList {
	listLen := gofakeit.IntN(10) // max of 10 items
	data := make([]any, listLen)
	bCodeData := ""
	invalid := false
	for i := range listLen {
		randType := gofakeit.IntN(3)
		if level >= 5 { // max of 5 nested levels
			randType = gofakeit.IntN(1) // do not select list or dict
		}
		invalid = switchNestedData(randType, &(data[i]), &bCodeData, level)
	}

	start := gofakeit.LetterN(1)
	end := gofakeit.LetterN(1)
	bCode := start + bCodeData + end

	return testList{
		data:      data,
		bCodeData: bCodeData,
		invalid:   invalid,
		start:     start,
		end:       end,
		bCode:     bCode,
	}
}

func genTestDict(level int) testDict {
	dictLen := gofakeit.IntN(10) // max of 10 items
	data := make(map[string]any, dictLen)
	bCodeData := ""
	invalid := false
	for range dictLen {
		randKeyType := gofakeit.IntN(3)
		randValType := gofakeit.IntN(3)
		if level >= 5 { // max of 5 nested levels
			randKeyType = gofakeit.IntN(1) // do not select list or dict
			randValType = gofakeit.IntN(1)
		}
		var key any
		var val any
		invalid = switchNestedData(randKeyType, &key, &bCodeData, level)
		invalid = switchNestedData(randValType, &key, &bCodeData, level)
		if reflect.TypeOf(key).Kind() != reflect.String {
			invalid = true
		}
		data[fmt.Sprintf("%v", key)] = val
	}

	start := gofakeit.LetterN(1)
	end := gofakeit.LetterN(1)
	bCode := start + bCodeData + end

	return testDict{
		data:      data,
		bCodeData: bCodeData,
		invalid:   invalid,
		start:     start,
		end:       end,
		bCode:     bCode,
	}
}

func TestDecodeString(t *testing.T) {
	t.Parallel()
	for range SimNumbers {
		seed := gofakeit.Int64()
		if gofakeit.Seed(seed) != nil {
			continue
		}

		test := genTestString()
		decoded, err := Decode([]byte(test.bCode))
		if isStringInvalid(test) {
			assert.Error(t, err, printSeed(seed)+"\nBencode: "+test.bCode)
		} else {
			if assert.NoError(t, err, printSeed(seed)) {
				assert.Equal(t, test.randStr[:test.randLen], decoded.(string), printSeed(seed))
			}
		}
	}
}

func TestDecodeInt(t *testing.T) {
	t.Parallel()
	for range SimNumbers {
		seed := gofakeit.Int64()
		if gofakeit.Seed(seed) != nil {
			continue
		}

		test := genTestInt()
		decoded, err := Decode([]byte(test.bCode))
		if isIntInvalid(test) {
			assert.Error(t, err, printSeed(seed)+"\nBencode: "+test.bCode)
		} else {
			if assert.NoError(t, err, printSeed(seed)) {
				assert.Equal(t, test.randInt, decoded.(int), printSeed(seed))
			}
		}
	}
}

func TestDecodeList(t *testing.T) {
	t.Parallel()
	for range SimNumbers {
		seed := gofakeit.Int64()
		if gofakeit.Seed(seed) != nil {
			continue
		}

		test := genTestList(0)

		// Start rune may be valid
		if test.start == "d" || test.start == "i" {
			continue
		}

		decoded, err := Decode([]byte(test.bCode))
		if isListInvalid(test) {
			if test.bCode[1] == byte('e') {
				continue
			}
			assert.Error(t, err, printSeed(seed)+"\nBencode: "+test.bCode)
		} else {
			if assert.NoError(t, err, printSeed(seed)) {
				if len(test.data) == 0 {
					assert.Equal(t, make([]interface{}, 0), decoded, printSeed(seed))
				} else {
					assert.Equal(t, test.bCodeData, decoded.([]any), printSeed(seed))
				}
			}
		}
	}
}

func TestDecodeDict(t *testing.T) {
	t.Parallel()
	for range SimNumbers {
		seed := gofakeit.Int64()
		if gofakeit.Seed(seed) != nil {
			continue
		}

		test := genTestDict(0)

		// Start rune may be valid
		if test.start == "l" || test.start == "i" {
			continue
		}

		decoded, err := Decode([]byte(test.bCode))
		if isDictInvalid(test) {
			if test.bCode[1] == byte('e') {
				continue
			}
			assert.Error(t, err, printSeed(seed)+"\nBencode: "+test.bCode)
		} else {
			if assert.NoError(t, err, printSeed(seed)) {
				if len(test.data) == 0 {
					assert.Equal(t, make(map[string]interface{}, 0), decoded, printSeed(seed))
				} else {
					assert.Equal(t, test.bCodeData, decoded.([]any), printSeed(seed))
				}
			}
		}
	}
}
