package apirouter

import "testing"

//TestSnakeCase test our snake case method
func TestSnakeCase(t *testing.T) {

	//Test a valid case
	word := "testCamelCase"
	result := SnakeCase(word)

	if result != "test_camel_case" {
		t.Fatal("SnakeCase conversion failed!", result)
	}

	//Test a valid case
	word = "TestCamelCase"
	result = SnakeCase(word)

	if result != "test_camel_case" {
		t.Fatal("SnakeCase conversion failed!", result)
	}

	//Test a valid case
	word = "TEstCamelCase"
	result = SnakeCase(word)

	if result != "test_camel_case" {
		t.Fatal("SnakeCase conversion failed!", result)
	}

	//Test a valid case
	word = "testCamelCASE"
	result = SnakeCase(word)

	if result != "test_camel_case" {
		t.Fatal("SnakeCase conversion failed!", result)
	}

	//Test a valid case
	word = "testCamel!CASE"
	result = SnakeCase(word)

	if result != "test_camel_case" {
		t.Fatal("SnakeCase conversion failed!", result)
	}
}

//TestFindString test our find string method
func TestFindString(t *testing.T) {
	//needle string, haystack []string
	haystack := []string{"test", "stack"}
	needle := "stack"

	if index := FindString(needle, haystack); index == -1 {
		t.Fatal("FindString does not work correctly!")
	}

	if index := FindString("wrong", haystack); index >= 0 {
		t.Fatal("FindString does not work correctly!")
	}

}
