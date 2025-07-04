package serialize

import (
	"bytes"
	"ebe/types"
	"ebe/utils"
)

func SerializeSint(value int64, data *bytes.Buffer) {

	// Get the negative sign and the abs of the data since we will store the value as
	// an unsigned integer with the high bit used as the negative sign
	var negative bool = value < 0
	var v uint64 = utils.Abs(value)

	// Figure out what size of integer is needed for the data
	// Note since we need the high bit for the negative sign any value that uses the high
	// bit needs to add an extra byte so we need to check for values without the high bit set
	var length uint8 = 0
	switch {

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

	// Set the header for the type in the data buffer
	data.WriteByte(types.CreateHeader(types.SInt, length))

	// Move the data into the data buffer
	for i := uint8(length - 1); ; i-- {

		var mask uint64 = 0xff << (i * 8)            // Create the mask to clear all bits other than they byte at position i
		var maskedValue = v & mask                   // Clear all bits other than the byte at position i
		var byteValue = byte(maskedValue >> (i * 8)) // Get the byte at position i

		// For the first byte, if the value is negative then set the high bit
		if i == uint8(length-1) && negative {
			byteValue = byteValue | 0x80
		}

		data.WriteByte(byteValue)

		if i == 0 {
			break
		}
	}
}

func DeserializeSint(data []byte) (int64, []byte) {

	if len(data) == 0 {
		return 0, data
	}

	var header = data[0]
	data = data[1:]

	var headerType = types.TypeFromHeader(header)
	var length = types.ValueFromHeader(header)

	if headerType != types.SInt {
		return 0, data
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
	data = data[length:]

	if negative {
		value = -value
	}

	return value, data
}
