package serialize

import (
	"bytes"
	"ebe/types"
)

func SerializeString(value string, data *bytes.Buffer) {

	// Write the length of the string as an [UInt]
	var length = len(value)

	// Special case strings that are less than 8 characters in length by putting the length
	// in the lsb nibble of the header instead of writing a full [UInt].
	// The high bit of the nibble will be 0 if the length is in the nibble and will be 1 if the length is in a following UInt
	// Note: it is legal to have a zero length string so zero can't be used as the indicator
	if length <= 0x07 {
		data.WriteByte(types.CreateHeader(types.String, byte(length)))
	} else {
		data.WriteByte(types.CreateHeader(types.String, 0x08))
		SerializeUint(uint64(length), data)
	}

	// Write the raw string data
	data.WriteString(value)
}

func DeserializeString(data []byte) (string, []byte) {

	if len(data) == 0 {
		return "", data
	}

	var header = data[0]
	data = data[1:]

	var headerType = types.TypeFromHeader(header)

	if headerType != types.String {
		return "", data
	}

	var length = uint64(types.ValueFromHeader(header))

	// If the 4th bit of the nibble is not set then this is a special cased string whose length fits in the length nibble
	// Otherwise get the string length from integer in the next data type
	if length == 8 {
		length, data = DeserializeUint(data[0:])
	}

	var value = string(data[0:length])
	data = data[length:]

	return value, data
}
