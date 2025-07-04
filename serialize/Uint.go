package serialize

import (
	"bytes"
	"ebe/types"
	"math"
)

func SerializeUint(value uint64, data *bytes.Buffer) {

	// Figure out what size of integer is needed for the data
	var length uint8 = 0
	switch {

	// Special case integer values less than 15
	// Instead of using additional bytes for the value,
	// put the value in the lsb nibble of the header
	case value <= 0x0f:
		data.WriteByte(types.CreateHeader(types.UIntNibble, byte(value)))
		return

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

	// Set the header for the type in the data buffer
	data.WriteByte(types.CreateHeader(types.UInt, length))

	// Move the data into the data buffer
	for i := uint8(length - 1); ; i-- {
		var mask uint64 = 0xff << (i * 8)            // Create the mask to clear all bits other than they byte at position i
		var maskedValue = value & mask               // Clear all bits other than the byte at position i
		var byteValue = byte(maskedValue >> (i * 8)) // Get the byte at position i

		data.WriteByte(byteValue)

		if i == 0 {
			break
		}
	}
}

func DeserializeUint(data []byte) (uint64, []byte) {

	if len(data) == 0 {
		return 0, data
	}

	var header = data[0]
	data = data[1:]

	var headerType = types.TypeFromHeader(header)
	var length = types.ValueFromHeader(header)

	// Values that fit into the length nibble of the header are special
	// cased by storing the length in the header itself
	if headerType == types.UIntNibble {
		return uint64(length), data
	}

	// Make sure the data is a valid integer value
	if headerType != types.UInt {
		return 0, data
	}

	var value uint64 = 0
	for i := 0; i < int(length); i++ {
		var dataByte = data[i]
		value = value << 8
		value = value | uint64(dataByte)
	}
	data = data[length:]

	return value, data
}
