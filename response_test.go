package apirouter

import (
	"bytes"
	"encoding/json"
	"testing"
)

//TestJsonEncode tests the the json encoder removes fields that are not approved
func TestJsonEncode(t *testing.T) {

	//Setup a mock struct for testing
	type TestStruct struct {
		TestKey    string `json:"test_key"`
		TestKeyTwo string `json:"test_key_two"`
		notAllowed string
	}

	//Base model and test model
	var model = new(TestStruct)
	var modelTest = new(TestStruct)
	var allowedFields = []string{"test_key", "test_key_two"} //notice omitted: notAllowed

	//Set the testing data
	model.TestKey = "TestValue1"
	model.TestKeyTwo = "TestValue2"
	model.notAllowed = "PrivateValue"

	//Set the buffer and encoder
	var b bytes.Buffer
	enc := json.NewEncoder(&b)

	//Run the encoder
	err := JSONEncode(enc, model, allowedFields)
	if err != nil {
		t.Fatal(err)
	}

	//Now unmarshal and test
	err = json.Unmarshal([]byte(b.String()), &modelTest)
	if err != nil {
		t.Fatal(err)
	}

	//Test for our fields and values now
	if modelTest.TestKey != "TestValue1" {
		t.Fatal("TestKey does not have the right value! Encoding failed.", modelTest.TestKey)
	} else if modelTest.TestKeyTwo != "TestValue2" {
		t.Fatal("TestKeyTwo does not have the right value! Encoding failed.", modelTest.TestKeyTwo)
	} else if modelTest.notAllowed == "PrivateValue" {
		t.Fatal("Field not removed! notAllowed does not have the right value! Encoding failed.", modelTest.notAllowed)
	}
}

//TestJsonEncodeSubstruct tests the the json encoder removes fields that are not approved
func TestJsonEncodeSubstruct(t *testing.T) {

	//Setup a mock substruct
	type TestSubStruct struct {
		TestSubKey string `json:"test_sub_key"`
	}
	//Setup a mock struct for testing
	type TestStruct struct {
		TestKey    string        `json:"test_key"`
		TestKeyTwo TestSubStruct `json:"test_key_two"`
		NotAllowed string        `json:"not_allowed"`
	}

	//Base model and test model
	var model = new(TestStruct)
	var modelTest = new(TestStruct)
	var allowedFields = []string{"test_key", "test_key_two"} //notice omitted: notAllowed

	//Set the testing data
	model.TestKey = "TestValue1"
	model.TestKeyTwo.TestSubKey = "TestSubValue"
	model.NotAllowed = "PrivateValue"

	//Set the buffer and encoder
	var b bytes.Buffer
	enc := json.NewEncoder(&b)

	//Run the encoder
	err := JSONEncode(enc, model, allowedFields)
	if err != nil {
		t.Fatal(err)
	}

	//Now unmarshal and test
	err = json.Unmarshal([]byte(b.String()), &modelTest)
	if err != nil {
		t.Fatal(err)
	}

	//Test for our fields and values now
	if modelTest.TestKey != "TestValue1" {
		t.Fatal("TestKey does not have the right value! Encoding failed.", modelTest.TestKey)
	} else if modelTest.TestKeyTwo.TestSubKey != "TestSubValue" {
		t.Fatal("TestKeyTwo does not have the right value! Encoding failed.", modelTest.TestKeyTwo)
	} else if modelTest.NotAllowed == "PrivateValue" {
		t.Fatal("Field not removed! notAllowed does not have the right value! Encoding failed.", modelTest.NotAllowed)
	}

}
