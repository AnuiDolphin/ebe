package serialize

import (
	"ebe/types"
	"ebe/utils"
	"fmt"
	"io"
)

// serializeString serializes a string value to the writer
func serializeString(value string, w io.Writer) error {

	// Write the string length header
	var length = len(value)

	// Strings under a certain length can use the shorter format
	if length <= 0x07 {
		utils.WriteByte(w, types.CreateHeader(types.String, byte(length)))
	} else {
		utils.WriteByte(w, types.CreateHeader(types.String, 0x08))
		if err := serializeUint(uint64(length), w); err != nil {
			return err
		}
	}

	// Write the raw string data
	_, err := w.Write([]byte(value))
	return err
}

// deserializeString deserializes a string with a pre-read header byte
func deserializeString(r io.Reader, header byte) (string, error) {
	headerType := types.TypeFromHeader(header)
	headerValue := types.ValueFromHeader(header)

	if headerType != types.String {
		return "", fmt.Errorf("expected String type, got %v", types.TypeName(headerType))
	}

	length := uint64(headerValue)

	// If the high bit of the length is set, then the length is in the next data type
	if length&0x08 != 0 {
		actualLength, err := deserializeUintWithHeader(r)
		if err != nil {
			return "", fmt.Errorf("failed to deserialize string length: %w", err)
		}
		length = actualLength
	}

	// Read the actual string data
	data := make([]byte, length)
	n, err := io.ReadFull(r, data)
	if err != nil {
		return "", fmt.Errorf("failed to read string data: %w", err)
	}
	if n != int(length) {
		return "", fmt.Errorf("expected to read %d bytes, got %d", length, n)
	}

	return string(data), nil
}
