package serialize

import (
	"bytes"
	"ebe/types"
	"fmt"
)

func SerializeBuffer(value []byte, data *bytes.Buffer) {
	// This function appends the serialized buffer to the existing buffer
	// Write the length of the buffer as an [UInt]
	var length = len(value)

	// Special case buffers that are less than 8 characters in length by putting the length
	// in the lsb nibble of the header instead of writing a full [UInt].
	// The high bit of the nibble will be 0 if the length is in the nibble and will be 1 if the length is in a following UInt
	// Note: it is legal to have a zero length buffer so zero can't be used as the indicator
	if length <= 0x07 {
		data.WriteByte(types.CreateHeader(types.Buffer, byte(length)))
	} else {
		data.WriteByte(types.CreateHeader(types.Buffer, 0x08))
		SerializeUint64(uint64(length), data)
	}

	// Write the raw buffer data
	data.Write(value)
}

func DeserializeBuffer(data []byte) (*bytes.Buffer, []byte, error) {

	if len(data) == 0 {
		return new(bytes.Buffer), data, fmt.Errorf("no data to deserialize")
	}

	var header = data[0]
	data = data[1:]

	var headerType = types.TypeFromHeader(header)
	var value = new(bytes.Buffer)

	if headerType != types.Buffer {
		return value, data, fmt.Errorf("expected Buffer type, got %v", headerType)
	}

	var length = uint64(types.ValueFromHeader(header))

	// If the 4th bit of the nibble is not set then this is a special cased buffer whose length fits in the length nibble
	// Otherwise get the buffer length from integer in the next data type
	if length == 8 {
		var err error
		length, data, err = DeserializeUint(data[0:])
		if err != nil {
			return value, data, fmt.Errorf("failed to deserialize buffer length: %w", err)
		}
	}

	if len(data) < int(length) {
		return value, data, fmt.Errorf("insufficient data: need %d bytes, got %d", length, len(data))
	}

	value.Write(data[0:length])
	data = data[length:]

	return value, data, nil
}
