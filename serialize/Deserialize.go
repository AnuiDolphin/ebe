package serialize

import (
	"ebe/types"
	"fmt"
)

// Deserialize reads the serialized type from the header and calls the appropriate deserialize method
func Deserialize(data []byte) (interface{}, []byte, error) {
	if len(data) == 0 {
		return nil, data, fmt.Errorf("empty data")
	}

	header := data[0]
	headerType := types.TypeFromHeader(header)

	switch headerType {
	case types.UNibble:
		value, remaining, err := DeserializeUNibble(data)
		if err != nil {
			return nil, remaining, err
		}
		return value, remaining, nil

	case types.SNibble:
		value, remaining, err := DeserializeSNibble(data)
		if err != nil {
			return nil, remaining, err
		}
		return value, remaining, nil

	case types.UInt:
		value, remaining, err := DeserializeUint64(data)
		if err != nil {
			return nil, remaining, err
		}
		return value, remaining, nil

	case types.SInt:
		value, remaining, err := DeserializeSint64(data)
		if err != nil {
			return nil, remaining, err
		}
		return value, remaining, nil

	case types.Float:
		value, remaining, err := DeserializeFloat64(data)
		if err != nil {
			return nil, remaining, err
		}
		return value, remaining, nil

	case types.Boolean:
		value, remaining, err := DeserializeBoolean(data)
		if err != nil {
			return nil, remaining, err
		}
		return value, remaining, nil

	case types.String:
		value, remaining, err := DeserializeString(data)
		if err != nil {
			return nil, remaining, err
		}
		return value, remaining, nil

	case types.Buffer:
		value, remaining, err := DeserializeBuffer(data)
		if err != nil {
			return nil, remaining, err
		}
		return value, remaining, nil

	default:
		return nil, data, fmt.Errorf("unsupported type: %s", types.TypeName(headerType))
	}
}

// DeserializeAll deserializes multiple values from a byte slice
func DeserializeAll(data []byte) ([]interface{}, error) {
	var results []interface{}
	remaining := data

	for len(remaining) > 0 {
		value, newRemaining, err := Deserialize(remaining)
		if err != nil {
			return results, err
		}
		results = append(results, value)
		remaining = newRemaining
	}

	return results, nil
}
