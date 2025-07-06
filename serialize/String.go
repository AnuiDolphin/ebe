package serialize

import (
	"ebe/types"
	"ebe/utils"
	"fmt"
	"io"
)

func SerializeString(value string, w io.Writer) error {
	// This function appends the serialized string to the existing writer
	// Write the length of the string as an [UInt]
	var length = len(value)

	// Special case strings that are less than 8 characters in length by putting the length
	// in the lsb nibble of the header instead of writing a full [UInt].
	// The high bit of the nibble will be 0 if the length is in the nibble and will be 1 if the length is in a following UInt
	// Note: it is legal to have a zero length string so zero can't be used as the indicator
	if length <= 0x07 {
		utils.WriteByte(w, types.CreateHeader(types.String, byte(length)))
	} else {
		utils.WriteByte(w, types.CreateHeader(types.String, 0x08))
		if err := SerializeUint64(uint64(length), w); err != nil {
			return err
		}
	}

	// Write the raw string data
	_, err := w.Write([]byte(value))
	return err
}

func DeserializeString(data []byte) (string, []byte, error) {

	if len(data) == 0 {
		return "", data, fmt.Errorf("no data to deserialize")
	}

	var header = data[0]
	data = data[1:]

	var headerType = types.TypeFromHeader(header)

	if headerType != types.String {
		return "", data, fmt.Errorf("expected String type, got %v", types.TypeName(headerType))
	}

	var length = uint64(types.ValueFromHeader(header))

	// If the 4th bit of the nibble is not set then this is a special cased string whose length fits in the length nibble
	// Otherwise get the string length from integer in the next data type
	var err error = nil
	if length&0x08 != 0 {
		length, data, err = DeserializeUint64(data)
		if err != nil {
			return "", data, fmt.Errorf("failed to deserialize string length: %w", err)
		}
	}

	if uint64(len(data)) < length {
		return "", data, fmt.Errorf("insufficient data: need %d bytes, have %d", length, len(data))
	}

	var value = string(data[0:length])
	data = data[length:]

	return value, data, err
}
