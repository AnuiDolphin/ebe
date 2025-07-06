package serialize

import (
	"ebe/types"
	"fmt"
	"io"
)

func serializeBoolean(value bool, w io.Writer) error {
	// This function appends the serialized boolean to the existing writer
	// Set the header for the type and put the boolean value in the value nibble
	var b byte
	if value {
		b = types.CreateHeader(types.Boolean, 1)
	} else {
		b = types.CreateHeader(types.Boolean, 0)
	}
	_, err := w.Write([]byte{b})
	return err
}

func deserializeBoolean(data []byte) (bool, []byte, error) {

	if len(data) == 0 {
		return false, data, fmt.Errorf("no data to deserialize")
	}

	var header = data[0]
	data = data[1:]

	var headerType = types.TypeFromHeader(header)

	if headerType != types.Boolean {
		return false, data, fmt.Errorf("expected Boolean type, got %v", headerType)
	}

	var value = types.ValueFromHeader(header)
	return value != 0, data, nil
}
