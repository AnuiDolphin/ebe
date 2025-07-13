package serialize

import (
	"ebe/types"
	"ebe/utils"
	"fmt"
	"io"
)

func serializeSint(value int64, writer io.Writer) error {

	// Get the negative sign and the abs of the data since we will store the value as
	// an unsigned integer with the high bit used as the negative sign
	var negative bool = value < 0
	var v uint64 = utils.Abs(value)

	// Figure out what size of integer is needed for the data
	// Note since we need the high bit for the negative sign any value that uses the high
	// bit needs to add an extra byte so we need to check for values without the high bit set
	var length uint8 = 0
	switch {

	case v <= 0x07:
		// For SNibble, we encode sign in bit 3 and magnitude in bits 0-2
		// This gives us range -7 to +7 (magnitude 0-7)
		var nibble byte
		if negative {
			nibble = 0x8 | byte(v) // Set bit 3 for negative, store magnitude in bits 0-2
		} else {
			nibble = byte(v) // Just use the magnitude for positive values
		}
		return utils.WriteByte(writer, types.CreateHeader(types.SNibble, nibble))

	case v <= 0x7f:
		length = 1

	case v <= 0x7fff:
		length = 2

	case v <= 0x7fffff:
		length = 3

	case v <= 0x7fffffff:
		length = 4

	case v <= 0x7fffffffff:
		length = 5

	case v <= 0x7fffffffffff:
		length = 6

	case v <= 0x7fffffffffffff:
		length = 7

	default:
		length = 8
	}

	// Set the header for the type
	if err := utils.WriteByte(writer, types.CreateHeader(types.SInt, length)); err != nil {
		return err
	}

	// Write the data bytes
	return writeValueBytes(writer, v, length, negative)
}

func deserializeSint(r io.Reader) (int64, error) {

	// Read the header using utils.ReadHeader
	headerType, headerValue, err := utils.ReadHeader(r)
	if err != nil {
		return 0, fmt.Errorf("failed to read SInt64 header: %w", err)
	}

	length := headerValue

	// If the header is a SNibble, we can return the value directly
	// For SNibble: bit 3 is sign, bits 0-2 are magnitude
	if headerType == types.SNibble {
		var negative = (length & 0x08) != 0
		var magnitude = length & 0x7 // Get bits 0-2

		if negative {
			return -int64(magnitude), nil
		} else {
			return int64(magnitude), nil
		}
	}

	if headerType != types.SInt {
		return 0, fmt.Errorf("expected SInt type, got %v", headerType)
	}

	// Read the data bytes
	data := make([]byte, length)
	_, err = io.ReadFull(r, data)
	if err != nil {
		return 0, fmt.Errorf("failed to read SInt64 data: %w", err)
	}

	// Get the value of the first byte of the data making sure to skip the negative bit
	var value = int64(data[0] & 0x7f)

	// If the high bit of the first byte is set then the value is negative
	var negative = (data[0] & 0x80) != 0

	// Add the rest of the data bytes to the value
	for i := uint8(1); i < length; i++ {
		value = value << 8
		var dataByte = int64(data[i])
		value = value | dataByte
	}

	if negative {
		value = -value
	}

	return value, nil
}

// Helper function to write value bytes in reverse order (big-endian)
// Sets the high bit of the first byte if negative is true
func writeValueBytes(writer io.Writer, value uint64, length uint8, negative bool) error {
	for i := uint8(length - 1); ; i-- {
		var byteValue = byte(value >> (i * 8))
		if i == uint8(length-1) && negative {
			byteValue = byteValue | 0x80
		}
		if err := utils.WriteByte(writer, byteValue); err != nil {
			return err
		}
		if i == 0 {
			break
		}
	}
	return nil
}
