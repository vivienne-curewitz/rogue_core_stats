package saveparser

import (
	"encoding/json"
	"os"
	"testing"
)

func readFile(filePath string) ([]byte, error) {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return fileBytes, nil
}

func TestSaveParser(t *testing.T) {
	filePath := "example.sav"
	data, err := readFile(filePath)
	if err != nil {
		t.Errorf("Could not load testing file: %s\n", err)
	}
	_, err = ConvertUesaveToJSON(data)
	if err != nil {
		t.Errorf("Failed to parse file: %s\n", err)
	}
}

func TestSaveParserStructure(t *testing.T) {
	filePath := "example.sav"
	fdata, err := readFile(filePath)
	if err != nil {
		t.Errorf("Could not load testing file: %s\n", err)
	}
	jsonStr, err := ConvertUesaveToJSON(fdata)
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
