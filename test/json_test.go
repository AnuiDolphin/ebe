package main_test

import (
	"bytes"
	"ebe/serialize"
	"encoding/json"
	"testing"
)

// Test data structures for JSON serialization
type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Company struct {
	Name      string   `json:"name"`
	Employees []Person `json:"employees"`
	Active    bool     `json:"active"`
}

func TestJsonSerializeDeserialize(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
	}{
		{
			name:  "simple object",
			value: Person{Name: "John", Age: 30},
		},
		{
			name: "complex object",
			value: Company{
				Name: "Tech Corp",
				Employees: []Person{
					{Name: "Alice", Age: 25},
					{Name: "Bob", Age: 35},
				},
				Active: true,
			},
		},
		{
			name:  "array",
			value: []string{"apple", "banana", "cherry"},
		},
		{
			name:  "map",
			value: map[string]int{"a": 1, "b": 2, "c": 3},
		},
		{
			name:  "empty object",
			value: Person{},
		},
		{
			name:  "small string", // Should use header nibble
			value: "hi",
		},
		{
			name:  "large object", // Should use UInt for length
			value: Company{Name: "Very Long Company Name That Exceeds Seven Characters", Active: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to JSON first using standard library
			jsonBytes, err := json.Marshal(tt.value)
			if err != nil {
				t.Fatalf("json.Marshal failed: %v", err)
			}

			// Serialize the JSON as json.RawMessage using our library
			var buffer bytes.Buffer
			err = serialize.SerializeJson(json.RawMessage(jsonBytes), &buffer)
			if err != nil {
				t.Fatalf("SerializeJson failed: %v", err)
			}

			serializedData := buffer.Bytes()
			if len(serializedData) == 0 {
				t.Fatal("SerializeJson produced empty data")
			}

			// Create a variable of the same type for deserialization
			// We need to handle different types appropriately
			var result interface{}
			switch tt.value.(type) {
			case Person:
				result = &Person{}
			case Company:
				result = &Company{}
			case []string:
				result = &[]string{}
			case map[string]int:
				result = &map[string]int{}
			case string:
				var s string
				result = &s
			}

			// Deserialize using our library
			remaining, err := serialize.DeserializeJson(serializedData, result)
			if err != nil {
				t.Fatalf("DeserializeJson failed: %v", err)
			}

			if len(remaining) != 0 {
				t.Errorf("DeserializeJson left %d bytes remaining", len(remaining))
			}

			// Compare the results by marshaling both to JSON and comparing
			originalJSON, err := json.Marshal(tt.value)
			if err != nil {
				t.Fatalf("Failed to marshal original value: %v", err)
			}

			// Dereference the result pointer for comparison
			var resultValue interface{}
			switch v := result.(type) {
			case *Person:
				resultValue = *v
			case *Company:
				resultValue = *v
			case *[]string:
				resultValue = *v
			case *map[string]int:
				resultValue = *v
			case *string:
				resultValue = *v
			}

			resultJSON, err := json.Marshal(resultValue)
			if err != nil {
				t.Fatalf("Failed to marshal result value: %v", err)
			}

			if !bytes.Equal(originalJSON, resultJSON) {
				t.Errorf("JSON content mismatch:\nOriginal: %s\nResult:   %s", originalJSON, resultJSON)
			}
		})
	}
}

func TestJsonIntegrationWithSerialize(t *testing.T) {
	// Test that we can use json.RawMessage with the main Serialize function
	person := Person{Name: "Alice", Age: 30}

	// Create JSON using standard library
	jsonBytes, err := json.Marshal(person)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}
	jsonValue := json.RawMessage(jsonBytes)

	var buffer bytes.Buffer
	err = serialize.Serialize(jsonValue, &buffer)
	if err != nil {
		t.Fatalf("Serialize failed: %v", err)
	}

	// Deserialize using main Deserialize function (should auto-detect JSON type)
	var result Person
	remaining, err := serialize.Deserialize(bytes.NewReader(buffer.Bytes()), &result)
	if err != nil {
		t.Fatalf("Deserialize failed: %v", err)
	}

	if len(remaining) != 0 {
		t.Errorf("Deserialize left %d bytes remaining", len(remaining))
	}

	if result.Name != person.Name || result.Age != person.Age {
		t.Errorf("Round-trip failed: expected %+v, got %+v", person, result)
	}
}

func TestJsonVsRegularSerialization(t *testing.T) {
	// Test the difference between JSON and regular struct serialization
	person := Person{Name: "Bob", Age: 25}

	// Serialize as regular struct
	var structBuffer bytes.Buffer
	err := serialize.Serialize(person, &structBuffer)
	if err != nil {
		t.Fatalf("Struct serialize failed: %v", err)
	}

	// Serialize as JSON
	var jsonBuffer bytes.Buffer
	jsonBytes, err := json.Marshal(person)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}
	jsonValue := json.RawMessage(jsonBytes)
	err = serialize.Serialize(jsonValue, &jsonBuffer)
	if err != nil {
		t.Fatalf("JSON serialize failed: %v", err)
	}

	// The serialized data should be different
	if bytes.Equal(structBuffer.Bytes(), jsonBuffer.Bytes()) {
		t.Error("JSON and struct serialization should produce different output")
	}

	// Struct should deserialize with regular Deserialize
	var structResult Person
	_, err = serialize.Deserialize(bytes.NewReader(structBuffer.Bytes()), &structResult)
	if err != nil {
		t.Fatalf("Struct deserialize failed: %v", err)
	}

	// JSON should deserialize with main Deserialize (auto-detects JSON type)
	var jsonResult Person
	_, err = serialize.Deserialize(bytes.NewReader(jsonBuffer.Bytes()), &jsonResult)
	if err != nil {
		t.Fatalf("JSON deserialize failed: %v", err)
	}

	// Both results should equal the original
	if structResult != person {
		t.Errorf("Struct round-trip failed: expected %+v, got %+v", person, structResult)
	}
	if jsonResult != person {
		t.Errorf("JSON round-trip failed: expected %+v, got %+v", person, jsonResult)
	}
}
