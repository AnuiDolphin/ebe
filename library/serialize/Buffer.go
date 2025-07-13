package serialize

import (
	"bytes"
	"ebe/types"
	"ebe/utils"
	"fmt"
	"io"
)

func serializeBuffer(value []byte, w io.Writer) error {
	// This function appends the serialized buffer to the existing writer
	// Write the length of the buffer as an [UInt]
	var length = len(value)

	// Special case buffers that are less than 8 characters in length by putting the length
	// in the lsb nibble of the header instead of writing a full [UInt].
	// The high bit of the nibble will be 0 if the length is in the nibble and will be 1 if the length is in a following UInt
	// Note: it is legal to have a zero length buffer so zero can't be used as the indicator
	if length <= 0x07 {
		utils.WriteByte(w, types.CreateHeader(types.Buffer, byte(length)))
	} else {
		utils.WriteByte(w, types.CreateHeader(types.Buffer, 0x08))
		if err := serializeUint(uint64(length), w); err != nil {
			return err
		}
	}

	// Write the raw buffer data
	_, err := w.Write(value)
	return err
}

// deserializeBuffer deserializes a buffer with a pre-read header byte
func deserializeBuffer(r io.Reader, header byte) (*bytes.Buffer, error) {
	value := new(bytes.Buffer)
	
	headerType := types.TypeFromHeader(header)
	headerValue := types.ValueFromHeader(header)

	if headerType != types.Buffer {
		return value, fmt.Errorf("expected Buffer type, got %v", types.TypeName(headerType))
	}

	length := uint64(headerValue)

	// If the high bit of the length is set, then the length is in the next data type
	if length&0x08 != 0 {
		l, err := deserializeUintWithHeader(r)
		if err != nil {
			return value, fmt.Errorf("failed to read buffer length: %w", err)
		}
		length = l
	}

	// Read the buffer data
	data := make([]byte, length)
	_, err := io.ReadFull(r, data)
	if err != nil {
		return value, fmt.Errorf("failed to read buffer data: %w", err)
	}
	value.Write(data)
	return value, nil
}
