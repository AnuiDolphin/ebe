package serialize

import (
	"bytes"
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

	// Check if this is an empty struct first (before reading header)
	if outValue.Kind() == reflect.Struct && isStructEmpty(outValue) {
		return nil
	}

	// Peek at the header type to determine how to deserialize
	header, err := utils.ReadByte(r)
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}
	headerType := types.TypeFromHeader(header)
	headerReader := io.MultiReader(bytes.NewReader([]byte{header}), r)

	// For JSON, we need to pass the header back into the stream
	if headerType == types.Json {
		return deserializeJson(headerReader, out)
	}

	// For structs, pass the header back into the stream
	if outValue.Kind() == reflect.Struct {
		return deserializeStruct(headerReader, outValue)
	}

	// For simple types, pass the header back into the stream
	return deserializeSimpleType(headerReader, header, outValue)
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

	// Check if all fields are unexported (private)
	for i := 0; i < outValue.NumField(); i++ {
		field := outValue.Type().Field(i)
		if field.PkgPath == "" { // Exported field
			return false
		}
	}
	return true
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
		value, err := deserializeBoolean(r)
		if err != nil {
			return err
		}
		if err := utils.SetValueWithConversion(outValue, value); err != nil {
			return fmt.Errorf("failed to set Boolean value: %w", err)
		}
		return nil

	case types.UInt:
		value, err := deserializeUint(r)
		if err != nil {
			return err
		}
		if err := utils.SetValueWithConversion(outValue, value); err != nil {
			return fmt.Errorf("failed to set UInt value: %w", err)
		}
		return nil

	case types.SInt:
		value, err := deserializeSint(r)
		if err != nil {
			return err
		}
		if err := utils.SetValueWithConversion(outValue, value); err != nil {
			return fmt.Errorf("failed to set SInt value: %w", err)
		}
		return nil

	case types.Float:
		value, err := deserializeFloat(r)
		if err != nil {
			return err
		}
		if err := utils.SetValueWithConversion(outValue, value); err != nil {
			return fmt.Errorf("failed to set Float value: %w", err)
		}
		return nil

	case types.String:
		value, err := deserializeString(r)
		if err != nil {
			return err
		}
		if err := utils.SetValueWithConversion(outValue, value); err != nil {
			return fmt.Errorf("failed to set String value: %w", err)
		}
		return nil

	case types.Buffer:
		value, err := deserializeBuffer(r)
		if err != nil {
			return err
		}
		if err := utils.SetValueWithConversion(outValue, value); err != nil {
			return fmt.Errorf("failed to set Buffer value: %w", err)
		}
		return nil

	case types.Array:
		if err := deserializeArray(r, outValue.Addr().Interface()); err != nil {
			return err
		}
		return nil

	default:
		return fmt.Errorf("unsupported type: %s", types.TypeName(headerType))
	}
}
