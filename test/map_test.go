package main_test

import (
	"bytes"
	"ebe/serialize"
	"ebe/types"
	"ebe/utils"
	"fmt"
	"reflect"
	"testing"
)

// TestBasicMaps tests basic map functionality with different key/value types
func TestBasicMaps(t *testing.T) {
	// String to int map
	testMap(t, map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}, "map[string]int")

	// String to string map
	testMap(t, map[string]string{
		"hello": "world",
		"foo":   "bar",
	}, "map[string]string")

	// Int to string map
	testMap(t, map[int]string{
		1: "one",
		2: "two",
		3: "three",
	}, "map[int]string")

	// Bool to int map
	testMap(t, map[bool]int{
		true:  1,
		false: 0,
	}, "map[bool]int")
}

// TestNumericKeyMaps tests maps with various numeric key types
func TestNumericKeyMaps(t *testing.T) {
	// Different integer key types
	testMap(t, map[int8]string{
		-1: "negative one",
		0:  "zero",
		1:  "positive one",
	}, "map[int8]string")

	testMap(t, map[uint64]string{
		0:    "zero",
		255:  "max uint8",
		1000: "thousand",
	}, "map[uint64]string")

	testMap(t, map[int32]float64{
		-100: -1.5,
		0:    0.0,
		100:  1.5,
	}, "map[int32]float64")
}

// TestVariousValueTypes tests maps with different value types
func TestVariousValueTypes(t *testing.T) {
	// String to bytes
	testMap(t, map[string][]byte{
		"binary": {0x01, 0x02, 0x03},
		"empty":  {},
		"text":   []byte("hello"),
	}, "map[string][]byte")

	// String to bool
	testMap(t, map[string]bool{
		"enabled":  true,
		"disabled": false,
		"active":   true,
	}, "map[string]bool")

	// String to float (test with simpler values to avoid precision issues)
	testMap(t, map[string]float64{
		"zero": 0.0,
		"one":  1.0,
		"two":  2.0,
	}, "map[string]float64")
}

// TestEmptyMap tests empty map serialization/deserialization
func TestEmptyMap(t *testing.T) {
	testMap(t, map[string]int{}, "empty map[string]int")
	testMap(t, map[int]string{}, "empty map[int]string")
	testMap(t, map[bool][]byte{}, "empty map[bool][]byte")
}

// TestLargeMaps tests maps that exceed the 7-entry header optimization
func TestLargeMaps(t *testing.T) {
	// Create a map with exactly 8 entries (triggers overflow encoding)
	largeMap := make(map[string]int)
	for i := 0; i < 8; i++ {
		largeMap[fmt.Sprintf("key%d", i)] = i * i
	}
	testMap(t, largeMap, "map with 8 entries")

	// Create a larger map
	veryLargeMap := make(map[int]string)
	for i := 0; i < 100; i++ {
		veryLargeMap[i] = fmt.Sprintf("value_%d", i)
	}
	testMap(t, veryLargeMap, "map with 100 entries")
}

// TestMixedValueMaps tests maps where values can have different types
// Note: This tests the serialization framework, not mixed types in a single map
func TestMixedValueMaps(t *testing.T) {
	// Test maps with interface{} values containing different types
	testMapInterface(t, map[string]interface{}{
		"int":    42,
		"string": "hello",
		"bool":   true,
		"float":  3.14,
	}, "map[string]interface{} with mixed types")
}

// TestNestedMaps tests maps containing other maps
func TestNestedMaps(t *testing.T) {
	nestedMap := map[string]map[string]int{
		"group1": {
			"a": 1,
			"b": 2,
		},
		"group2": {
			"x": 10,
			"y": 20,
			"z": 30,
		},
	}
	testMap(t, nestedMap, "nested map[string]map[string]int")
}

// TestMapHeaderEncoding tests the header encoding for different entry counts
func TestMapHeaderEncoding(t *testing.T) {
	testCases := []struct {
		entryCount   int
		expectedByte byte
		description  string
	}{
		{0, 0xB0, "empty map (count=0 in header)"}, // Type=11 (0xB), Count=0
		{1, 0xB1, "1 entry (count=1 in header)"},  // Type=11 (0xB), Count=1
		{7, 0xB7, "7 entries (count=7 in header)"}, // Type=11 (0xB), Count=7
		{8, 0xB8, "8 entries (overflow indicator)"}, // Type=11 (0xB), Count=8 (overflow)
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			// Create map with specified entry count
			testMap := make(map[string]int)
			for i := 0; i < tc.entryCount; i++ {
				testMap[fmt.Sprintf("key%d", i)] = i
			}

			// Serialize the map
			var buf bytes.Buffer
			err := serialize.Serialize(testMap, &buf)
			if err != nil {
				t.Fatalf("Failed to serialize map: %v", err)
			}

			// Check that the first byte matches expected header
			data := buf.Bytes()
			if len(data) == 0 {
				t.Fatal("No data serialized")
			}

			if data[0] != tc.expectedByte {
				t.Errorf("Expected header byte 0x%02X, got 0x%02X", tc.expectedByte, data[0])
			}

			// For overflow case (8+ entries), verify UInt encoding follows
			if tc.entryCount >= 8 {
				// Skip the map header byte
				reader := bytes.NewReader(data[1:])
				
				// The next bytes should be a UInt encoding of the actual count
				actualCount, err := deserializeUintFromReader(reader)
				if err != nil {
					t.Fatalf("Failed to read entry count: %v", err)
				}
				
				if actualCount != uint64(tc.entryCount) {
					t.Errorf("Expected entry count %d, got %d", tc.entryCount, actualCount)
				}
			}
		})
	}
}

// Helper function to deserialize UInt from reader (for testing)
func deserializeUintFromReader(r *bytes.Reader) (uint64, error) {
	headerType, headerValue, err := utils.ReadHeader(r)
	if err != nil {
		return 0, err
	}

	if headerType == types.UNibble {
		return uint64(headerValue), nil
	}

	if headerType != types.UInt {
		return 0, fmt.Errorf("expected UInt, got %v", headerType)
	}

	length := headerValue
	data := make([]byte, length)
	_, err = r.Read(data)
	if err != nil {
		return 0, err
	}

	var value uint64 = 0
	for i := 0; i < int(length); i++ {
		value = value << 8
		value = value | uint64(data[i])
	}

	return value, nil
}

// testMap is a helper function that tests map serialization and deserialization
func testMap(t *testing.T, originalMap interface{}, description string) {
	var buf bytes.Buffer

	// Serialize the map
	err := serialize.Serialize(originalMap, &buf)
	if err != nil {
		t.Errorf("%s: Serialization failed: %v", description, err)
		return
	}

	// Create a new map of the same type for deserialization
	mapType := reflect.TypeOf(originalMap)
	newMapPtr := reflect.New(mapType)
	newMap := newMapPtr.Interface()

	// Deserialize the map
	err = serialize.Deserialize(bytes.NewReader(buf.Bytes()), newMap)
	if err != nil {
		t.Errorf("%s: Deserialization failed: %v", description, err)
		return
	}

	// Compare the maps
	if !reflect.DeepEqual(originalMap, newMapPtr.Elem().Interface()) {
		t.Errorf("%s: Maps not equal after round-trip", description)
		t.Logf("Original: %+v", originalMap)
		t.Logf("Deserialized: %+v", newMapPtr.Elem().Interface())
	} else {
		t.Logf("%s: PASS", description)
	}
}

// testMapInterface is a special helper for maps with interface{} values
func testMapInterface(t *testing.T, originalMap map[string]interface{}, description string) {
	var buf bytes.Buffer

	// Serialize the map
	err := serialize.Serialize(originalMap, &buf)
	if err != nil {
		t.Errorf("%s: Serialization failed: %v", description, err)
		return
	}

	// For interface{} maps, we need to handle comparison carefully
	// since the types might not match exactly after deserialization
	newMap := make(map[string]interface{})

	// Deserialize the map
	err = serialize.Deserialize(bytes.NewReader(buf.Bytes()), &newMap)
	if err != nil {
		t.Errorf("%s: Deserialization failed: %v", description, err)
		return
	}

	// Compare keys and values using custom logic
	if len(originalMap) != len(newMap) {
		t.Errorf("%s: Map lengths differ: original=%d, deserialized=%d", 
			description, len(originalMap), len(newMap))
		return
	}

	for key, originalValue := range originalMap {
		deserializedValue, exists := newMap[key]
		if !exists {
			t.Errorf("%s: Key %q missing in deserialized map", description, key)
			continue
		}

		// Use utils.CompareValue for type-aware comparison
		if !utils.CompareValue(originalValue, deserializedValue) {
			t.Errorf("%s: Value mismatch for key %q: original=%v (%T), deserialized=%v (%T)", 
				description, key, originalValue, originalValue, deserializedValue, deserializedValue)
		}
	}

	t.Logf("%s: PASS", description)
}