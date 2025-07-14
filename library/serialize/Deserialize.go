package serialize

import (
	"ebe/types"
	"ebe/utils"
	"fmt"
	"io"
	"reflect"
)

// Deserialize reads the serialized type from the header and deserializes into the provided output parameter
func Deserialize(r io.Reader, out interface{}) error {
	
	// Validate and get the output value from within the interface{}
	outValue, err := getOutputValue(out)
	if err != nil {
		return err
	}

	// Check if this is an empty struct first - they serialize to 0 bytes (no header)
	if outValue.Kind() == reflect.Struct && isStructEmpty(outValue) {
		return nil
	}
	
	// Read the header byte and delegate to the internal deserializer
	header, err := utils.ReadByte(r)
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}
	
	return deserializeWithHeaderInternal(r, header, out, outValue)
}

// deserializeWithHeader deserializes data with a pre-read header byte (internal use only)
func deserializeWithHeader(r io.Reader, header byte, out interface{}) error {

	// Validate and get the output value from within the interface{}
	outValue, err := getOutputValue(out)
	if err != nil {
		return err
	}

	return deserializeWithHeaderInternal(r, header, out, outValue)
}

// deserializeWithHeaderInternal performs the actual deserialization with pre-validated outValue
func deserializeWithHeaderInternal(r io.Reader, header byte, out interface{}, outValue reflect.Value) error {
	
	headerType := types.TypeFromHeader(header)

	// For JSON, parse with header parameter
	// This has to happen before struct deserialization so we can handle the JSON type correctly
	if headerType == types.Json {
		return deserializeJson(r, header, out)
	}

	// For structs, validate that we have a struct type header and use header-aware struct deserialization
	if outValue.Kind() == reflect.Struct {
		if headerType != types.Struct {
			return fmt.Errorf("expected Struct header type for struct value, got %v", types.TypeName(headerType))
		}
		return deserializeStruct(r, header, outValue)
	}

	// For simple types, pass the header directly
	return deserializeSimpleType(r, header, outValue)
}

// getOutputValue validates the output parameter and returns the reflect.Value to set
func getOutputValue(out interface{}) (reflect.Value, error) {

	// Make sure the output parameter is not nil
	if out == nil {
		return reflect.Value{}, fmt.Errorf("output parameter cannot be nil")
	}

	// Get the reflect value of the output parameter and ensure it's a pointer so we can assign values to it
	refectValue := reflect.ValueOf(out)
	if refectValue.Kind() != reflect.Ptr {
		return reflect.Value{}, fmt.Errorf("output parameter must be a pointer")
	}

	// Get the element that the pointer points to
	outValue := refectValue.Elem()
	if !outValue.CanSet() {
		return reflect.Value{}, fmt.Errorf("output parameter must be settable")
	}

	return outValue, nil
}

func isStructEmpty(outValue reflect.Value) bool {

	// Check if the value is a struct and has no exported fields
	if outValue.Kind() != reflect.Struct {
		return false
	}

	// Use cached struct information for performance
	structInfo, err := typeCache.GetStructInfo(outValue.Type())
	if err != nil {
		// Fallback to direct reflection on error
		for i := 0; i < outValue.NumField(); i++ {
			field := outValue.Type().Field(i)
			if field.PkgPath == "" { // Exported field
				return false
			}
		}
		return true
	}
	
	return structInfo.Empty
}

// deserializeSimpleType deserializes data from a stream into a simple (non-struct) type
func deserializeSimpleType(r io.Reader, header byte, outValue reflect.Value) error {

	headerType := types.TypeFromHeader(header)

	switch headerType {
	case types.UNibble:
		value := types.ValueFromHeader(header)
		if err := utils.SetValueWithConversion(outValue, value); err != nil {
			return fmt.Errorf("failed to set UNibble value: %w", err)
		}
		return nil

	case types.SNibble:
		var negative = (header & 0x8) != 0
		var magnitude = header & 0x7
		value := int8(magnitude)
		if negative {
			value = -value
		}
		if err := utils.SetValueWithConversion(outValue, value); err != nil {
			return fmt.Errorf("failed to set SNibble value: %w", err)
		}
		return nil

	case types.Boolean:
		value, err := deserializeBoolean(r, header)
		if err != nil {
			return err
		}
		if err := utils.SetValueWithConversion(outValue, value); err != nil {
			return fmt.Errorf("failed to set Boolean value: %w", err)
		}
		return nil

	case types.UInt:
		value, err := deserializeUint(r, header)
		if err != nil {
			return err
		}
		if err := utils.SetValueWithConversion(outValue, value); err != nil {
			return fmt.Errorf("failed to set UInt value: %w", err)
		}
		return nil

	case types.SInt:
		value, err := deserializeSint(r, header)
		if err != nil {
			return err
		}
		if err := utils.SetValueWithConversion(outValue, value); err != nil {
			return fmt.Errorf("failed to set SInt value: %w", err)
		}
		return nil

	case types.Float:
		value, err := deserializeFloat(r, header)
		if err != nil {
			return err
		}
		if err := utils.SetValueWithConversion(outValue, value); err != nil {
			return fmt.Errorf("failed to set Float value: %w", err)
		}
		return nil

	case types.String:
		value, err := deserializeString(r, header)
		if err != nil {
			return err
		}
		if err := utils.SetValueWithConversion(outValue, value); err != nil {
			return fmt.Errorf("failed to set String value: %w", err)
		}
		return nil

	case types.Buffer:
		value, err := deserializeBuffer(r, header)
		if err != nil {
			return err
		}
		if err := utils.SetValueWithConversion(outValue, value); err != nil {
			return fmt.Errorf("failed to set Buffer value: %w", err)
		}
		return nil

	case types.Array:
		if err := deserializeArray(r, header, outValue.Addr().Interface()); err != nil {
			return err
		}
		return nil

	case types.Map:
		if err := deserializeMap(r, header, outValue.Addr().Interface()); err != nil {
			return err
		}
		return nil

	case types.Struct:
		if err := deserializeStruct(r, header, outValue); err != nil {
			return err
		}
		return nil

	default:
		return fmt.Errorf("unsupported type: %s", types.TypeName(headerType))
	}
}

// Type-specific deserializers that avoid reflection overhead
// These provide fast paths for the most commonly used types

// deserializeInt64 deserializes an int64 value directly without reflection
func deserializeInt64(r io.Reader) (int64, error) {
	header, err := utils.ReadByte(r)
	if err != nil {
		return 0, fmt.Errorf("failed to read header: %w", err)
	}

	headerType := types.TypeFromHeader(header)

	switch headerType {
	case types.SNibble:
		var negative = (header & 0x8) != 0
		var magnitude = header & 0x7
		value := int64(magnitude)
		if negative {
			value = -value
		}
		return value, nil

	case types.UNibble:
		value := types.ValueFromHeader(header)
		return int64(value), nil

	case types.SInt:
		return deserializeSint(r, header)

	case types.UInt:
		uvalue, err := deserializeUint(r, header)
		if err != nil {
			return 0, err
		}
		// Check for overflow when converting uint64 to int64
		if uvalue > 9223372036854775807 { // math.MaxInt64
			return 0, fmt.Errorf("uint64 value %d overflows int64", uvalue)
		}
		return int64(uvalue), nil

	default:
		return 0, fmt.Errorf("cannot deserialize %s as int64", types.TypeName(headerType))
	}
}

// deserializeUint64 deserializes a uint64 value directly without reflection
func deserializeUint64(r io.Reader) (uint64, error) {
	header, err := utils.ReadByte(r)
	if err != nil {
		return 0, fmt.Errorf("failed to read header: %w", err)
	}

	headerType := types.TypeFromHeader(header)

	switch headerType {
	case types.UNibble:
		value := types.ValueFromHeader(header)
		return uint64(value), nil

	case types.UInt:
		return deserializeUint(r, header)

	case types.SNibble:
		// Convert SNibble to uint64 if non-negative
		var negative = (header & 0x8) != 0
		if negative {
			return 0, fmt.Errorf("cannot convert negative SNibble to uint64")
		}
		var magnitude = header & 0x7
		return uint64(magnitude), nil

	case types.SInt:
		svalue, err := deserializeSint(r, header)
		if err != nil {
			return 0, err
		}
		if svalue < 0 {
			return 0, fmt.Errorf("cannot convert negative int64 %d to uint64", svalue)
		}
		return uint64(svalue), nil

	default:
		return 0, fmt.Errorf("cannot deserialize %s as uint64", types.TypeName(headerType))
	}
}

// deserializeFloat64 deserializes a float64 value directly without reflection
func deserializeFloat64(r io.Reader) (float64, error) {
	header, err := utils.ReadByte(r)
	if err != nil {
		return 0, fmt.Errorf("failed to read header: %w", err)
	}

	headerType := types.TypeFromHeader(header)

	switch headerType {
	case types.Float:
		return deserializeFloat(r, header)

	case types.UNibble:
		value := types.ValueFromHeader(header)
		return float64(value), nil

	case types.SNibble:
		var negative = (header & 0x8) != 0
		var magnitude = header & 0x7
		value := float64(magnitude)
		if negative {
			value = -value
		}
		return value, nil

	case types.UInt:
		uvalue, err := deserializeUint(r, header)
		if err != nil {
			return 0, err
		}
		return float64(uvalue), nil

	case types.SInt:
		svalue, err := deserializeSint(r, header)
		if err != nil {
			return 0, err
		}
		return float64(svalue), nil

	default:
		return 0, fmt.Errorf("cannot deserialize %s as float64", types.TypeName(headerType))
	}
}


