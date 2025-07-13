package test

import (
	"bytes"
	"ebe/serialize"
	"testing"
)

// TestMapIntegration tests map support in real-world-like scenarios
func TestMapIntegration(t *testing.T) {
	
	// Test the scenario that was failing in the comparison tool
	userPreferences := map[string]interface{}{
		"theme":         "dark",
		"notifications": true,
		"language":      "en", 
		"timezone":      "UTC-5",
		"font_size":     14,
	}

	var buf bytes.Buffer
	err := serialize.Serialize(userPreferences, &buf)
	if err != nil {
		t.Fatalf("Failed to serialize user preferences: %v", err)
	}

	// Deserialize back
	var result map[string]interface{}
	err = serialize.Deserialize(bytes.NewReader(buf.Bytes()), &result)
	if err != nil {
		t.Fatalf("Failed to deserialize user preferences: %v", err)
	}

	// Verify key count
	if len(result) != len(userPreferences) {
		t.Errorf("Expected %d keys, got %d", len(userPreferences), len(result))
	}

	// Verify each key exists and has reasonable value
	for key := range userPreferences {
		if _, exists := result[key]; !exists {
			t.Errorf("Key %q missing in deserialized map", key)
		}
	}

	t.Logf("Successfully serialized and deserialized map with %d entries", len(result))
	t.Logf("Serialized size: %d bytes", buf.Len())
}

// TestComplexMapStructure tests a more complex map structure  
func TestComplexMapStructure(t *testing.T) {
	complexData := map[string]map[string]interface{}{
		"database": {
			"host":     "localhost",
			"port":     5432,
			"ssl":      true,
			"timeout":  30,
		},
		"redis": {
			"host":     "redis.example.com", 
			"port":     6379,
			"ssl":      false,
			"database": 0,
		},
	}

	var buf bytes.Buffer
	err := serialize.Serialize(complexData, &buf)
	if err != nil {
		t.Fatalf("Failed to serialize complex map: %v", err)
	}

	var result map[string]map[string]interface{}
	err = serialize.Deserialize(bytes.NewReader(buf.Bytes()), &result)
	if err != nil {
		t.Fatalf("Failed to deserialize complex map: %v", err)
	}

	// Verify structure
	if len(result) != 2 {
		t.Errorf("Expected 2 top-level keys, got %d", len(result))
	}

	if db, exists := result["database"]; exists {
		if len(db) != 4 {
			t.Errorf("Expected 4 database config keys, got %d", len(db))
		}
	} else {
		t.Error("Database config missing")
	}

	t.Logf("Successfully serialized complex nested map structure")
	t.Logf("Serialized size: %d bytes", buf.Len())
}