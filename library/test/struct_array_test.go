package test

import (
	"bytes"
	"ebe/serialize"
	"reflect"
	"testing"
)

// Test structs for struct arrays
type PersonStruct struct {
	ID   uint32
	Name string
	Age  int32
}

type CompanyStruct struct {
	ID       uint64
	Name     string
	Founded  int32
	Revenue  float64
}

type EmptyStruct struct {
}

func TestBasicStructArray(t *testing.T) {
	people := []PersonStruct{
		{ID: 1, Name: "Alice", Age: 30},
		{ID: 2, Name: "Bob", Age: 25},
	}

	var buf bytes.Buffer
	err := serialize.Serialize(people, &buf)
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	var result []PersonStruct
	err = serialize.Deserialize(bytes.NewReader(buf.Bytes()), &result)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	if !reflect.DeepEqual(people, result) {
		t.Errorf("Round trip mismatch: expected %+v, got %+v", people, result)
	}
}

func TestEmptyStructArray(t *testing.T) {
	var people []PersonStruct

	var buf bytes.Buffer
	err := serialize.Serialize(people, &buf)
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	var result []PersonStruct
	err = serialize.Deserialize(bytes.NewReader(buf.Bytes()), &result)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected empty array, got %+v", result)
	}
}

func TestSingleElementStructArray(t *testing.T) {
	people := []PersonStruct{{ID: 42, Name: "Solo", Age: 99}}

	var buf bytes.Buffer
	err := serialize.Serialize(people, &buf)
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	var result []PersonStruct
	err = serialize.Deserialize(bytes.NewReader(buf.Bytes()), &result)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	if !reflect.DeepEqual(people, result) {
		t.Errorf("Round trip mismatch: expected %+v, got %+v", people, result)
	}
}

func TestLargeStructArray(t *testing.T) {
	// Test with more than 7 elements to trigger overflow handling
	people := make([]PersonStruct, 10)
	for i := range people {
		people[i] = PersonStruct{
			ID:   uint32(i + 1),
			Name: "Person" + string(rune('A'+i)),
			Age:  int32(20 + i),
		}
	}

	var buf bytes.Buffer
	err := serialize.Serialize(people, &buf)
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	var result []PersonStruct
	err = serialize.Deserialize(bytes.NewReader(buf.Bytes()), &result)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	if !reflect.DeepEqual(people, result) {
		t.Errorf("Round trip mismatch: expected %+v, got %+v", people, result)
	}
}

func TestDifferentStructTypes(t *testing.T) {
	companies := []CompanyStruct{
		{ID: 1, Name: "Tech Corp", Founded: 2020, Revenue: 1000000.0},
		{ID: 2, Name: "Data Inc", Founded: 2019, Revenue: 500000.0},
	}

	var buf bytes.Buffer
	err := serialize.Serialize(companies, &buf)
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	var result []CompanyStruct
	err = serialize.Deserialize(bytes.NewReader(buf.Bytes()), &result)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	if !reflect.DeepEqual(companies, result) {
		t.Errorf("Round trip mismatch: expected %+v, got %+v", companies, result)
	}
}

func TestEmptyStructsArray(t *testing.T) {
	emptyStructs := []EmptyStruct{{}, {}, {}}

	var buf bytes.Buffer
	err := serialize.Serialize(emptyStructs, &buf)
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	var result []EmptyStruct
	err = serialize.Deserialize(bytes.NewReader(buf.Bytes()), &result)
	if err != nil {
		t.Fatalf("Deserialization failed: %v", err)
	}

	if len(result) != 3 {
		t.Errorf("Expected 3 empty structs, got %d", len(result))
	}
}

func TestStructArraySizeComparison(t *testing.T) {
	// Test that struct arrays produce reasonable sizes
	people := []PersonStruct{
		{ID: 1, Name: "Alice", Age: 30},
		{ID: 2, Name: "Bob", Age: 25},
	}

	var buf bytes.Buffer
	err := serialize.Serialize(people, &buf)
	if err != nil {
		t.Fatalf("Serialization failed: %v", err)
	}

	// Log the size for analysis
	t.Logf("Serialized 2 Person structs: %d bytes", buf.Len())
	t.Logf("Data: %x", buf.Bytes())

	// Verify reasonable size (not too bloated)
	if buf.Len() > 100 {
		t.Errorf("Serialized size seems too large: %d bytes", buf.Len())
	}
}