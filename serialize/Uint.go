package serialize

import (
	"bytes"
	"ebe/types"
	"ebe/utils"
	"fmt"
	"io"
	"math"
)

func DeserializeUNibble(data []byte) (uint8, []byte, error) {

	if len(data) == 0 {
		return 0, data, fmt.Errorf("no data to deserialize")
	}

	var header = data[0]
	var headerType = types.TypeFromHeader(header)
	var length = types.ValueFromHeader(header)
	data = data[1:]

	// If the header is a UNibble, we can return the value directly
	if headerType != types.UNibble {
		return 0, data, fmt.Errorf("expected UNibble type, got %v", types.TypeNameFromHeader(header))
	}

	// For UNibble, the value is stored directly in the header (no negative values)
	return uint8(length), data, nil
}

func DeserializeUint8(data []byte) (uint8, []byte, error) {
	value, remainingData, err := DeserializeUint64(data)
	if err != nil {
		return 0, remainingData, err
	}
	if value > math.MaxUint8 {
		return 0, remainingData, fmt.Errorf("value %d does not fit in uint8 range [0, %d]", value, math.MaxUint8)
	}
	return uint8(value), remainingData, nil
}

func DeserializeUint16(data []byte) (uint16, []byte, error) {
	value, remainingData, err := DeserializeUint64(data)
	if err != nil {
		return 0, remainingData, err
	}
	if value > math.MaxUint16 {
		return 0, remainingData, fmt.Errorf("value %d does not fit in uint16 range [0, %d]", value, math.MaxUint16)
	}
	return uint16(value), remainingData, nil
}

func DeserializeUint32(data []byte) (uint32, []byte, error) {
	value, remainingData, err := DeserializeUint64(data)
	if err != nil {
		return 0, remainingData, err
	}
	if value > math.MaxUint32 {
		return 0, remainingData, fmt.Errorf("value %d does not fit in uint32 range [0, %d]", value, math.MaxUint32)
	}
	return uint32(value), remainingData, nil
}

func DeserializeUint64(data []byte) (uint64, []byte, error) {

	if len(data) == 0 {
		return 0, data, fmt.Errorf("no data to deserialize")
	}

	var header = data[0]
	var headerType = types.TypeFromHeader(header)
	var length = types.ValueFromHeader(header)
	data = data[1:]

	if len(data) == 0 {
		return 0, data, fmt.Errorf("no data to deserialize")
	}

	// For UNibble, the value is stored directly in the header (no negative values)
	if headerType == types.UNibble {
		return uint64(length), data, nil
	}

	// Make sure the data is a valid integer value
	if headerType != types.UInt {
		return 0, data, fmt.Errorf("expected UInt type, got %v", headerType)
	}

	if len(data) < int(length) {
		return 0, data, fmt.Errorf("insufficient data: need %d bytes, got %d", length, len(data))
	}

	var value uint64 = 0
	for i := 0; i < int(length); i++ {
		var dataByte = data[i]
		value = value << 8
		value = value | uint64(dataByte)
	}
	data = data[length:]

	return value, data, nil
}

func DeserializeUint(data []byte) (uint64, []byte, error) {
	return DeserializeUint64(data)
}

//
// Serialization functions
// These functions serialize unsigned integers into a byte buffer.
// All methods write to the writer.
//

func SerializeUint8(value uint8, writer io.Writer) error {
	// Serialize uint8 as uint64
	return SerializeUint64(uint64(value), writer)
}

func SerializeUint16(value uint16, writer io.Writer) error {
	// Serialize uint16 as uint64
	return SerializeUint64(uint64(value), writer)
}

func SerializeUint32(value uint32, writer io.Writer) error {
	// Serialize uint32 as uint64
	return SerializeUint64(uint64(value), writer)
}

func SerializeUint64(value uint64, writer io.Writer) error {
	// This function writes the serialized unsigned integer to the writer

	// Figure out what size of integer is needed for the data
	var length uint8 = 0
	switch {

	// Special case integer values less than or equal to 15
	// Instead of using additional bytes for the value,
	// put the value in the lsb nibble of the header
	case value <= 0x0f:
		if value == 0x00 {
			return utils.WriteByte(writer, types.CreateHeader(types.SNibble, byte(value)))
		} else {
			return utils.WriteByte(writer, types.CreateHeader(types.UNibble, byte(value)))
		}

	case value <= math.MaxUint8:
		length = 1

	case value <= math.MaxUint16:
		length = 2

	case value <= 0xffffff:
		length = 3

	case value <= math.MaxUint32:
		length = 4

	case value <= 0xffffffffff:
		length = 5

	case value <= 0xffffffffffff:
		length = 6

	case value <= 0xffffffffffffff:
		length = 7

	default:
		length = 8
	}

	// Set the header for the type
	if err := utils.WriteByte(writer, types.CreateHeader(types.UInt, length)); err != nil {
		return err
	}

	// Write the data bytes
	for i := uint8(length - 1); ; i-- {
		if err := utils.WriteByte(writer, byte(value >> (i * 8))); err != nil {
			return err
		}
		if i == 0 {
			break
		}
	}
	
	return nil
}

// Wrapper functions that preserve the *bytes.Buffer interface for backward compatibility

func SerializeUint8Buffer(value uint8, data *bytes.Buffer) error {
	return SerializeUint8(value, data)
}

func SerializeUint16Buffer(value uint16, data *bytes.Buffer) error {
	return SerializeUint16(value, data)
}

func SerializeUint32Buffer(value uint32, data *bytes.Buffer) error {
	return SerializeUint32(value, data)
}

func SerializeUint64Buffer(value uint64, data *bytes.Buffer) error {
	return SerializeUint64(value, data)
}
