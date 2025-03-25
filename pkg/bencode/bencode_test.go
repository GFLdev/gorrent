package bencode

import (
	"flag"
	"strconv"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

var SimNumbers = flag.Int("sim", 1000, "number of simulations to run")

const (
	Str int = iota
	Int
	List
	Dict
)

type testData[T interface{}] struct {
	data  T
	bCode string
}

func FormatSeed(seed int64) string {
	return "seed: " + strconv.Itoa(int(seed))
}

func FormatInfo(seed int64, bCode string) string {
	return FormatSeed(seed) + "\nbencode: " + bCode
}

type TestMarshalAndUnmarshalData struct {
	String  string                 `bencode:"string"`
	Integer int                    `bencode:"integer"`
	List    []interface{}          `bencode:"list"`
	Dict    map[string]interface{} `bencode:"dict"`
	Nothing interface{}
}

func TestMarshalAndUnmarshal(t *testing.T) {
	t.Parallel()
	for range *SimNumbers {
		seed := gofakeit.Int64()
		if gofakeit.Seed(seed) != nil {
			continue
		}

		strTest := genStringEncodeTest()
		intTest := genIntEncodeTest()
		listTest := genListEncodeTest(0)
		dictTest := genDictEncodeTest(0)
		nothingTest := switchGenType(gofakeit.IntN(4), 0)

		// Data
		testUnmarshal := TestMarshalAndUnmarshalData{}
		testMarshal := TestMarshalAndUnmarshalData{
			String:  strTest.data,
			Integer: intTest.data,
			List:    listTest.data,
			Dict:    dictTest.data,
			Nothing: nothingTest.data,
		}
		expectedBCode := "d4:dict" + dictTest.bCode +
			"7:integer" + intTest.bCode +
			"4:list" + listTest.bCode +
			"6:string" + strTest.bCode + "e"

		// Tests
		bCode, errMarshal := Marshal(&testMarshal)
		errUnmarshal := Unmarshal(bCode, &testUnmarshal)
		testMarshal.Nothing = nil
		if assert.NoError(t, errMarshal, FormatInfo(seed, expectedBCode)) { // marshal no error
			if assert.Equal(t, expectedBCode, string(bCode), FormatSeed(seed)) { // marshal data
				if assert.NoError(t, errUnmarshal, FormatInfo(seed, string(bCode))) { // unmarshal error
					assert.Equal(t, testMarshal, testUnmarshal, FormatSeed(seed)) // unmarshal data
				}
			}
		}
	}
}
