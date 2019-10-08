package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testObject struct {
	OptVal StringOption `json:"optVal"`
	Val    string       `json:"val"`
}

func TestParseStringOption_Unmarshal_NoOptVal(t *testing.T) {
	testData := []byte(`
		{
			"val": "TheVal"
		}`)

	tObj := testObject{}
	err := json.Unmarshal(testData, &tObj)
	assert.NoError(t, err)
	assert.Equal(t, "TheVal", tObj.Val)
	assert.Equal(t, false, tObj.OptVal.IsValid())
}

func TestParseStringOption_Unmarshal_NullOptVal(t *testing.T) {
	testData := []byte(`
		{
			"optVal": null,
			"val": "TheVal"
		}`)

	tObj := testObject{}
	err := json.Unmarshal(testData, &tObj)
	assert.NoError(t, err)
	assert.Equal(t, "TheVal", tObj.Val)
	assert.Equal(t, false, tObj.OptVal.IsValid())
}

func TestParseStringOption_Unmarshal_ValidOptVal(t *testing.T) {
	testData := []byte(`
		{
			"optVal": "TheOptVal",
			"val": "TheVal"
		}`)

	tObj := testObject{}
	err := json.Unmarshal(testData, &tObj)
	assert.NoError(t, err)
	assert.Equal(t, tObj.Val, "TheVal")
	assert.Equal(t, true, tObj.OptVal.IsValid())
	assert.Equal(t, "TheOptVal", tObj.OptVal.String())
}

func TestParseStringOption_Marshal_NoOptVal(t *testing.T) {
	tObj := testObject{
		Val: "TheVal",
	}
	expectedResult := []byte(`{"optVal":null,"val":"TheVal"}`)
	data, err := json.Marshal(&tObj)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, data)
}

func TestParseStringOption_Marshal_ValidOptVal(t *testing.T) {
	tObj := testObject{
		OptVal: StringOption{
			valid: true,
			value: "TheOptVal",
		},
		Val: "TheVal",
	}
	expectedResult := []byte(`{"optVal":"TheOptVal","val":"TheVal"}`)
	data, err := json.Marshal(&tObj)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, data)
}
