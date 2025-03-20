package bencode

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

const SimNumbers = 10000

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

func printSeed(seed int64) string {
	return "seed: " + strconv.Itoa(int(seed))
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
		switch randType {
		case str:
			test := genTestString()
			if test.randLen <= len(test.randStr) {
				data[i] = test.randStr[:test.randLen]
			} else {
				data[i] = test.randStr
			}
			bCodeData += test.bCode
			if isStringInvalid(test) {
				invalid = true
			}
		case integer:
			test := genTestInt()
			data[i] = test.randInt
			bCodeData += test.bCode
			if isIntInvalid(test) {
				invalid = true
			}
		case list:
			test := genTestList(level + 1)
			data[i] = test.data
			bCodeData += test.bCode
			if isListInvalid(test) {
				invalid = true
			}
		default:
			continue
		}
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
			assert.Error(t, err, printSeed(seed))
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
			assert.Error(t, err, printSeed(seed))
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
		decoded, err := Decode([]byte(test.bCode))
		if isListInvalid(test) {
			assert.Error(t, err, printSeed(seed))
		} else {
			if assert.NoError(t, err, printSeed(seed)) {
				assert.Equal(t, test.bCodeData, decoded.([]any), printSeed(seed))
			}
		}
	}
}
