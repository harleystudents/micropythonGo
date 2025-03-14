package micropythongo

import (
	"reflect"
	"testing"
)

func TestParseList(t *testing.T) {
	pyList := "['hello', 'world']"

	parsed := parsePythonList(pyList)

	expected := []string{"hello", "world"}
	if len(parsed) != len(expected) {
		t.Errorf("parsed list has incorrect length: got %v, want %v", len(parsed), len(expected))
	}
	for i, v := range parsed {
		if v != expected[i] {
			t.Errorf("parsed list has incorrect value at index %d: got %v, want %v", i, v, expected[i])
		}
	}

	parsed = parsePythonList("")
	check := []string{}

	if !reflect.DeepEqual(parsed, check) {
		t.Error("failed to recognixe empty string")
	}
}
