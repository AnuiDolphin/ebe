package serialize

import (
	"bytes"
	"ebe-library/types"
	"ebe-library/utils"
	"fmt"
	"io"
	"reflect"
)

// serializeMap serializes a Go map to the EBE format
// Format: [Map Header] [Optional Entry Count] [Key-Value Pairs...]
// Each key and value is a self-describing EBE value with its own header
func serializeMap(value interface{}, w io.Writer) error {
	
	// Try fast paths for common map types first
	switch m := value.(type) {

	case map[string]int:
		return serializeMapStringInt(m, w)

	case map[string]string:
		return serializeMapStringString(m, w)

	case map[string]interface{}:
		return serializeMapStringInterface(m, w)

	case map[int]string:
		return serializeMapIntString(m, w)

	case map[string]int32:
		return serializeMapStringInt32(m, w)

	case map[string]bool:
		return serializeMapStringBool(m, w)

	default:
		// Fall back to generic reflection-based approach
		return serializeMapGeneric(value, w)
	}
}

// serializeMapGeneric handles arbitrary map types using reflection (fallback)
func serializeMapGeneric(value interface{}, w io.Writer) error {

	// Validate input parameter
	rv := reflect.ValueOf(value)
	
	// Ensure we have a map
	if rv.Kind() != reflect.Map {
		return fmt.Errorf("expected map, got %v", rv.Kind())
	}
	
	entryCount := rv.Len()
	
	// Write map header with entry count optimization
	if entryCount <= 7 {

		// Small maps: store count directly in header nibble
		header := types.CreateHeader(types.Map, byte(entryCount))
		if err := utils.WriteByte(w, header); err != nil {
			return fmt.Errorf("failed to write map header: %w", err)
		}
	} else {

		// Large maps: use overflow indicator (8) in header, then UInt for actual count
		header := types.CreateHeader(types.Map, 8)
		if err := utils.WriteByte(w, header); err != nil {
			return fmt.Errorf("failed to write map header: %w", err)
		}
		
		// Write actual count as standard EBE UInt
		if err := serializeUint(uint64(entryCount), w); err != nil {
			return fmt.Errorf("failed to write map entry count: %w", err)
		}
	}
	
	// Write each key-value pair using standard EBE serialization
	// Each key and value is self-describing with its own header
	for _, key := range rv.MapKeys() {

		// Serialize key
		if err := Serialize(key.Interface(), w); err != nil {
			return fmt.Errorf("failed to serialize map key: %w", err)
		}
		
		// Serialize corresponding value
		value := rv.MapIndex(key)
		if err := Serialize(value.Interface(), w); err != nil {
			return fmt.Errorf("failed to serialize map value: %w", err)
		}
	}
	
	return nil
}

// Fast path serialization functions for common map types

// serializeMapStringInt serializes map[string]int without reflection
func serializeMapStringInt(m map[string]int, w io.Writer) error {

	entryCount := len(m)
	
	// Write map header
	if err := writeMapHeader(entryCount, w); err != nil {
		return err
	}
	
	// Write key-value pairs directly without reflection
	for key, value := range m {
		if err := serializeString(key, w); err != nil {
			return fmt.Errorf("failed to serialize string key: %w", err)
		}
		if err := serializeSint(int64(value), w); err != nil {
			return fmt.Errorf("failed to serialize int value: %w", err)
		}
	}
	
	return nil
}

// serializeMapStringString serializes map[string]string without reflection
func serializeMapStringString(m map[string]string, w io.Writer) error {

	entryCount := len(m)
	
	// Write map header
	if err := writeMapHeader(entryCount, w); err != nil {
		return err
	}
	
	// Write key-value pairs directly without reflection
	for key, value := range m {
		if err := serializeString(key, w); err != nil {
			return fmt.Errorf("failed to serialize string key: %w", err)
		}
		if err := serializeString(value, w); err != nil {
			return fmt.Errorf("failed to serialize string value: %w", err)
		}
	}
	
	return nil
}

// serializeMapStringInterface serializes map[string]interface{} with minimal reflection
func serializeMapStringInterface(m map[string]interface{}, w io.Writer) error {

	entryCount := len(m)
	
	// Write map header
	if err := writeMapHeader(entryCount, w); err != nil {
		return err
	}
	
	// Write key-value pairs - keys are direct, values need Serialize()
	for key, value := range m {
		if err := serializeString(key, w); err != nil {
			return fmt.Errorf("failed to serialize string key: %w", err)
		}
		if err := Serialize(value, w); err != nil {
			return fmt.Errorf("failed to serialize interface{} value: %w", err)
		}
	}
	
	return nil
}

// serializeMapIntString serializes map[int]string without reflection
func serializeMapIntString(m map[int]string, w io.Writer) error {

	entryCount := len(m)
	
	// Write map header
	if err := writeMapHeader(entryCount, w); err != nil {
		return err
	}
	
	// Write key-value pairs directly without reflection
	for key, value := range m {
		if err := serializeSint(int64(key), w); err != nil {
			return fmt.Errorf("failed to serialize int key: %w", err)
		}
		if err := serializeString(value, w); err != nil {
			return fmt.Errorf("failed to serialize string value: %w", err)
		}
	}
	
	return nil
}

// serializeMapStringInt32 serializes map[string]int32 without reflection
func serializeMapStringInt32(m map[string]int32, w io.Writer) error {

	entryCount := len(m)
	
	// Write map header
	if err := writeMapHeader(entryCount, w); err != nil {
		return err
	}
	
	// Write key-value pairs directly without reflection
	for key, value := range m {
		if err := serializeString(key, w); err != nil {
			return fmt.Errorf("failed to serialize string key: %w", err)
		}
		if err := serializeSint(int64(value), w); err != nil {
			return fmt.Errorf("failed to serialize int32 value: %w", err)
		}
	}
	
	return nil
}

// serializeMapStringBool serializes map[string]bool without reflection
func serializeMapStringBool(m map[string]bool, w io.Writer) error {

	entryCount := len(m)
	
	// Write map header
	if err := writeMapHeader(entryCount, w); err != nil {
		return err
	}
	
	// Write key-value pairs directly without reflection
	for key, value := range m {
		if err := serializeString(key, w); err != nil {
			return fmt.Errorf("failed to serialize string key: %w", err)
		}
		if err := serializeBoolean(value, w); err != nil {
			return fmt.Errorf("failed to serialize bool value: %w", err)
		}
	}
	
	return nil
}

// writeMapHeader writes the map header with entry count optimization
func writeMapHeader(entryCount int, w io.Writer) error {

	if entryCount <= 7 {

		// Small maps: store count directly in header nibble
		header := types.CreateHeader(types.Map, byte(entryCount))
		if err := utils.WriteByte(w, header); err != nil {
			return fmt.Errorf("failed to write map header: %w", err)
		}
	} else {

		// Large maps: use overflow indicator (8) in header, then UInt for actual count
		header := types.CreateHeader(types.Map, 8)
		if err := utils.WriteByte(w, header); err != nil {
			return fmt.Errorf("failed to write map header: %w", err)
		}
		
		// Write actual count as standard EBE UInt
		if err := serializeUint(uint64(entryCount), w); err != nil {
			return fmt.Errorf("failed to write map entry count: %w", err)
		}
	}
	return nil
}

// deserializeMap deserializes a map from the EBE format with fast paths for common types
// The output parameter is already validated by getOutputValue() to be a settable pointer
func deserializeMap(r io.Reader, out interface{}) error {

	// Try fast paths for common map types first
	switch m := out.(type) {
	case *map[string]int:
		return deserializeMapStringInt(r, m)
	case *map[string]string:
		return deserializeMapStringString(r, m)
	case *map[string]interface{}:
		return deserializeMapStringInterface(r, m)
	case *map[int]string:
		return deserializeMapIntString(r, m)
	case *map[string]int32:
		return deserializeMapStringInt32(r, m)
	case *map[string]bool:
		return deserializeMapStringBool(r, m)
	default:
		// Fall back to generic reflection-based approach
		return deserializeMapGeneric(r, out)
	}
}

// deserializeMapGeneric handles arbitrary map types using reflection (fallback)
func deserializeMapGeneric(r io.Reader, out interface{}) error {

	// Get the map value (already validated as a settable pointer by main Deserialize)
	mapValue := reflect.ValueOf(out).Elem()
	if mapValue.Kind() != reflect.Map {
		return fmt.Errorf("output parameter must be a pointer to a map, got pointer to %v", mapValue.Kind())
	}
	
	// Get map type information
	mapType := mapValue.Type()
	keyType := mapType.Key()
	valueType := mapType.Elem()
	
	// Read map header
	headerType, headerValue, err := utils.ReadHeader(r)
	if err != nil {
		return fmt.Errorf("failed to read map header: %w", err)
	}
	
	if headerType != types.Map {
		return fmt.Errorf("expected Map type, got %v", types.TypeName(headerType))
	}
	
	// Determine entry count
	var entryCount uint64
	if headerValue <= 7 {

		// Small map: count stored in header
		entryCount = uint64(headerValue)
	} else if headerValue == 8 {

		// Large map: read count as UInt
		entryCount, err = deserializeUint(r)
		if err != nil {
			return fmt.Errorf("failed to read map entry count: %w", err)
		}
	} else {
		return fmt.Errorf("invalid map header value: %d", headerValue)
	}
	
	// Initialize the map if it's nil
	if mapValue.IsNil() {
		mapValue.Set(reflect.MakeMap(mapType))
	}
	
	// Read each key-value pair
	for i := uint64(0); i < entryCount; i++ {
		
		// Create new instances for key and value
		keyPtr := reflect.New(keyType)
		valuePtr := reflect.New(valueType)
		
		// Deserialize key
		if err := Deserialize(r, keyPtr.Interface()); err != nil {
			return fmt.Errorf("failed to deserialize map key %d: %w", i, err)
		}
		
		// Deserialize value
		if err := Deserialize(r, valuePtr.Interface()); err != nil {
			return fmt.Errorf("failed to deserialize map value %d: %w", i, err)
		}
		
		// Add to map
		mapValue.SetMapIndex(keyPtr.Elem(), valuePtr.Elem())
	}
	
	return nil
}

// Fast path deserialization functions for common map types

// deserializeMapStringInt deserializes map[string]int without reflection
func deserializeMapStringInt(r io.Reader, out *map[string]int) error {
	entryCount, err := readMapHeader(r)
	if err != nil {
		return err
	}
	
	// Initialize map if nil
	if *out == nil {
		*out = make(map[string]int)
	}
	
	// Read each key-value pair directly
	for i := uint64(0); i < entryCount; i++ {

		// Deserialize string key
		key, err := deserializeString(r)
		if err != nil {
			return fmt.Errorf("failed to deserialize string key %d: %w", i, err)
		}
		
		// Deserialize int value
		value, err := deserializeSint(r)
		if err != nil {
			return fmt.Errorf("failed to deserialize int value %d: %w", i, err)
		}
		
		(*out)[key] = int(value)
	}
	
	return nil
}

// deserializeMapStringString deserializes map[string]string without reflection
func deserializeMapStringString(r io.Reader, out *map[string]string) error {

	entryCount, err := readMapHeader(r)
	if err != nil {
		return err
	}
	
	// Initialize map if nil
	if *out == nil {
		*out = make(map[string]string)
	}
	
	// Read each key-value pair directly
	for i := uint64(0); i < entryCount; i++ {

		// Deserialize string key
		key, err := deserializeString(r)
		if err != nil {
			return fmt.Errorf("failed to deserialize string key %d: %w", i, err)
		}
		
		// Deserialize string value
		value, err := deserializeString(r)
		if err != nil {
			return fmt.Errorf("failed to deserialize string value %d: %w", i, err)
		}
		
		(*out)[key] = value
	}
	
	return nil
}

// deserializeMapStringInterface deserializes map[string]interface{} with minimal reflection
func deserializeMapStringInterface(r io.Reader, out *map[string]interface{}) error {

	entryCount, err := readMapHeader(r)
	if err != nil {
		return err
	}
	
	// Initialize map if nil
	if *out == nil {
		*out = make(map[string]interface{})
	}
	
	// Read each key-value pair - keys are direct, values need Deserialize()
	for i := uint64(0); i < entryCount; i++ {

		// Deserialize string key
		key, err := deserializeString(r)
		if err != nil {
			return fmt.Errorf("failed to deserialize string key %d: %w", i, err)
		}
		
		// For interface{} values, we need to peek at the type and create appropriate value
		var value interface{}
		if err := deserializeInterfaceValue(r, &value); err != nil {
			return fmt.Errorf("failed to deserialize interface{} value %d: %w", i, err)
		}
		
		(*out)[key] = value
	}
	
	return nil
}

// deserializeMapIntString deserializes map[int]string without reflection
func deserializeMapIntString(r io.Reader, out *map[int]string) error {

	entryCount, err := readMapHeader(r)
	if err != nil {
		return err
	}
	
	// Initialize map if nil
	if *out == nil {
		*out = make(map[int]string)
	}
	
	// Read each key-value pair directly
	for i := uint64(0); i < entryCount; i++ {

		// Deserialize int key
		keyVal, err := deserializeSint(r)
		if err != nil {
			return fmt.Errorf("failed to deserialize int key %d: %w", i, err)
		}
		
		// Deserialize string value
		value, err := deserializeString(r)
		if err != nil {
			return fmt.Errorf("failed to deserialize string value %d: %w", i, err)
		}
		
		(*out)[int(keyVal)] = value
	}
	
	return nil
}

// deserializeMapStringInt32 deserializes map[string]int32 without reflection
func deserializeMapStringInt32(r io.Reader, out *map[string]int32) error {

	entryCount, err := readMapHeader(r)
	if err != nil {
		return err
	}
	
	// Initialize map if nil
	if *out == nil {
		*out = make(map[string]int32)
	}
	
	// Read each key-value pair directly
	for i := uint64(0); i < entryCount; i++ {
		// Deserialize string key
		key, err := deserializeString(r)
		if err != nil {
			return fmt.Errorf("failed to deserialize string key %d: %w", i, err)
		}
		
		// Deserialize int32 value
		value, err := deserializeSint(r)
		if err != nil {
			return fmt.Errorf("failed to deserialize int32 value %d: %w", i, err)
		}
		
		(*out)[key] = int32(value)
	}
	
	return nil
}

// deserializeMapStringBool deserializes map[string]bool without reflection
func deserializeMapStringBool(r io.Reader, out *map[string]bool) error {

	entryCount, err := readMapHeader(r)
	if err != nil {
		return err
	}
	
	// Initialize map if nil
	if *out == nil {
		*out = make(map[string]bool)
	}
	
	// Read each key-value pair directly
	for i := uint64(0); i < entryCount; i++ {
		// Deserialize string key
		key, err := deserializeString(r)
		if err != nil {
			return fmt.Errorf("failed to deserialize string key %d: %w", i, err)
		}
		
		// Deserialize bool value
		value, err := deserializeBoolean(r)
		if err != nil {
			return fmt.Errorf("failed to deserialize bool value %d: %w", i, err)
		}
		
		(*out)[key] = value
	}
	
	return nil
}

// readMapHeader reads and parses the map header, returning entry count
func readMapHeader(r io.Reader) (uint64, error) {

	headerType, headerValue, err := utils.ReadHeader(r)
	if err != nil {
		return 0, fmt.Errorf("failed to read map header: %w", err)
	}
	
	if headerType != types.Map {
		return 0, fmt.Errorf("expected Map type, got %v", types.TypeName(headerType))
	}
	
	// Determine entry count
	var entryCount uint64
	if headerValue <= 7 {
		// Small map: count stored in header
		entryCount = uint64(headerValue)
	} else if headerValue == 8 {
		// Large map: read count as UInt
		entryCount, err = deserializeUint(r)
		if err != nil {
			return 0, fmt.Errorf("failed to read map entry count: %w", err)
		}
	} else {
		return 0, fmt.Errorf("invalid map header value: %d", headerValue)
	}
	
	return entryCount, nil
}

// deserializeInterfaceValue deserializes a value into interface{} by peeking at the type
func deserializeInterfaceValue(r io.Reader, out *interface{}) error {

	// Peek at the header to determine type
	header, err := utils.ReadByte(r)
	if err != nil {
		return fmt.Errorf("failed to read value header: %w", err)
	}
	
	headerType := types.TypeFromHeader(header)
	headerReader := io.MultiReader(bytes.NewReader([]byte{header}), r)
	
	// Create appropriate concrete type based on header
	switch headerType {

	case types.UNibble, types.UInt:
		var value uint64
		if err := Deserialize(headerReader, &value); err != nil {
			return err
		}
		*out = value

	case types.SNibble, types.SInt:
		var value int64
		if err := Deserialize(headerReader, &value); err != nil {
			return err
		}
		*out = value

	case types.Float:
		var value float64
		if err := Deserialize(headerReader, &value); err != nil {
			return err
		}
		*out = value

	case types.Boolean:
		var value bool
		if err := Deserialize(headerReader, &value); err != nil {
			return err
		}
		*out = value

	case types.String:
		var value string
		if err := Deserialize(headerReader, &value); err != nil {
			return err
		}
		*out = value

	case types.Buffer:
		var value []byte
		if err := Deserialize(headerReader, &value); err != nil {
			return err
		}
		*out = value

	case types.Array:
		// For arrays in interface{}, we need to determine element type dynamically
		// Fall back to generic deserialization
		return fmt.Errorf("arrays in interface{} values require generic deserialization")

	case types.Map:
		// For nested maps in interface{}, fall back to generic deserialization
		return fmt.Errorf("nested maps in interface{} values require generic deserialization")
		
	default:
		return fmt.Errorf("unsupported type in interface{}: %s", types.TypeName(headerType))
	}
	
	return nil
}