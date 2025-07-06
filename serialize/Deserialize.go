package serialize

import (
	"ebe/types"
	"ebe/utils"
	"fmt"
	"reflect"
)

// Deserialize reads the serialized type from the header and deserializes into the provided output parameter
func Deserialize(data []byte, out interface{}) ([]byte, error) {

	// Validate and get the output value from within the interface{}
	outValue, err := getOutputValue(out)
	if err != nil {
		return data, err
	}

	// Check if we're deserializing into a struct
	if outValue.Kind() == reflect.Struct {
		return deserializeStruct(data, outValue)
	}

	// Deserialize simple types
	return deserializeSimpleType(data, outValue)
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

// deserializeStruct deserializes data into a struct by deserializing each field in order
func deserializeStruct(data []byte, structValue reflect.Value) ([]byte, error) {

	remaining := data
	structType := structValue.Type()

	// Iterate through each field in the struct
	for i := 0; i < structValue.NumField(); i++ {

		field := structValue.Field(i)
		fieldType := structType.Field(i)

		// Skip unexported fields
		if !field.CanSet() {
			continue
		}

		// Create a pointer to the field for deserialization
		fieldPtr := field.Addr().Interface()

		// Recursively call Deserialize to deserialize into this field
		newRemaining, err := Deserialize(remaining, fieldPtr)
		if err != nil {
			return remaining, fmt.Errorf("failed to deserialize field '%s': %w", fieldType.Name, err)
		}

		remaining = newRemaining
	}

	return remaining, nil
}

// deserializeSimpleType deserializes data into a simple (non-struct) type
func deserializeSimpleType(data []byte, outValue reflect.Value) ([]byte, error) {

	// Make sure there is data to deserialize for simple types
	if len(data) == 0 {
		return data, fmt.Errorf("empty data")
	}

	header := data[0]
	headerType := types.TypeFromHeader(header)

	switch headerType {
	case types.UNibble:
		value, remaining, err := DeserializeUNibble(data)
		if err != nil {
			return remaining, err
		}
		if err := utils.SetValueWithConversion(outValue, value); err != nil {
			return remaining, fmt.Errorf("failed to set UNibble value: %w", err)
		}
		return remaining, nil

	case types.SNibble:
		value, remaining, err := DeserializeSNibble(data)
		if err != nil {
			return remaining, err
		}
		if err := utils.SetValueWithConversion(outValue, value); err != nil {
			return remaining, fmt.Errorf("failed to set SNibble value: %w", err)
		}
		return remaining, nil

	case types.UInt:
		value, remaining, err := DeserializeUint64(data)
		if err != nil {
			return remaining, err
		}
		if err := utils.SetValueWithConversion(outValue, value); err != nil {
			return remaining, fmt.Errorf("failed to set UInt value: %w", err)
		}
		return remaining, nil

	case types.SInt:
		value, remaining, err := DeserializeSint64(data)
		if err != nil {
			return remaining, err
		}
		if err := utils.SetValueWithConversion(outValue, value); err != nil {
			return remaining, fmt.Errorf("failed to set SInt value: %w", err)
		}
		return remaining, nil

	case types.Float:
		value, remaining, err := DeserializeFloat64(data)
		if err != nil {
			return remaining, err
		}
		if err := utils.SetValueWithConversion(outValue, value); err != nil {
			return remaining, fmt.Errorf("failed to set Float value: %w", err)
		}
		return remaining, nil

	case types.Boolean:
		value, remaining, err := DeserializeBoolean(data)
		if err != nil {
			return remaining, err
		}
		if err := utils.SetValueWithConversion(outValue, value); err != nil {
			return remaining, fmt.Errorf("failed to set Boolean value: %w", err)
		}
		return remaining, nil

	case types.String:
		value, remaining, err := DeserializeString(data)
		if err != nil {
			return remaining, err
		}
		if err := utils.SetValueWithConversion(outValue, value); err != nil {
			return remaining, fmt.Errorf("failed to set String value: %w", err)
		}
		return remaining, nil

	case types.Buffer:
		value, remaining, err := DeserializeBuffer(data)
		if err != nil {
			return remaining, err
		}
		if err := utils.SetValueWithConversion(outValue, value); err != nil {
			return remaining, fmt.Errorf("failed to set Buffer value: %w", err)
		}
		return remaining, nil

	default:
		return data, fmt.Errorf("unsupported type: %s", types.TypeName(headerType))
	}
}
