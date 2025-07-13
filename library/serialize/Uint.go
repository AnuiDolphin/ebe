package serialize

import (
	"ebe/types"
	"ebe/utils"
	"fmt"
	"io"
	"math"
)

func serializeUint(value uint64, writer io.Writer) error {

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
		if err := utils.WriteByte(writer, byte(value>>(i*8))); err != nil {
			return err
		}
		if i == 0 {
			break
		}
	}

	return nil
}

func deserializeUint(r io.Reader) (uint64, error) {

	// Read the header using utils.ReadHeader
	headerType, headerValue, err := utils.ReadHeader(r)
	if err != nil {
		return 0, fmt.Errorf("failed to read UInt64 header: %w", err)
	}

	length := headerValue

	// For UNibble, the value is stored directly in the header (no negative values)
	if headerType == types.UNibble {
		return uint64(length), nil
	}

	// Make sure the data is a valid integer value
	if headerType != types.UInt {
		return 0, fmt.Errorf("expected UInt type, got %v", headerType)
	}

	// Read the data bytes
	data := make([]byte, length)
	_, err = io.ReadFull(r, data)
	if err != nil {
		return 0, fmt.Errorf("failed to read UInt64 data: %w", err)
	}

	var value uint64 = 0
	for i := 0; i < int(length); i++ {
		var dataByte = data[i]
		value = value << 8
		value = value | uint64(dataByte)
	}

	return value, nil
}
