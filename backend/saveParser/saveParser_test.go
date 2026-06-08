package saveparser

import (
	"encoding/json"
	"testing"
)

func TestSaveParser(t *testing.T) {
	filePath := "example.sav"
	_, err := ConvertUesaveToJSON(filePath)
	if err != nil {
		t.Errorf("Failed to parse file: %s\n", err)
	}
}

func TestSaveParserStructure(t *testing.T) {
	filePath := "example.sav"
	jsonStr, err := ConvertUesaveToJSON(filePath)
	if err != nil {
		t.Fatalf("Failed to parse file: %s\n", err)
	}

	var data map[string]interface{}
	err = json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %s\n", err)
	}

	expectedKeys := []string{"header", "schemas", "root", "extra"}
	for _, key := range expectedKeys {
		if _, ok := data[key]; !ok {
			keys := make([]string, 0, len(data))
			for k := range data {
				keys = append(keys, k)
			}
			t.Errorf("Expected key %q not found in JSON top level. Actual keys: %v", key, keys)
		}
	}
}
