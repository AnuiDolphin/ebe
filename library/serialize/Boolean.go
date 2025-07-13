package serialize

import (
	"ebe-library/types"
	"ebe-library/utils"
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
	return utils.WriteByte(w, b)
}

func deserializeBoolean(r io.Reader) (bool, error) {
	// Read the header using utils.ReadHeader
	headerType, headerValue, err := utils.ReadHeader(r)
	if err != nil {
		return false, fmt.Errorf("failed to read boolean header: %w", err)
	}

	if headerType != types.Boolean {
		return false, fmt.Errorf("expected Boolean type, got %v", headerType)
	}

	return headerValue != 0, nil
}
