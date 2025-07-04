package serialize

import (
	"bytes"
	"ebe/types"
	"encoding/binary"
	"math"
)

func SerializeFloat(value float64, data *bytes.Buffer) {

	// Figure out what size of float is needed for the data
	switch {

	case value >= -math.SmallestNonzeroFloat32 && value <= math.MaxFloat32:
		// Set the header for the type in the data buffer
		data.WriteByte(types.CreateHeader(types.Float, 4))

		// Write the value
		_ = binary.Write(data, binary.LittleEndian, float32(value))

	default:
		// Set the header for the type in the data buffer
		data.WriteByte(types.CreateHeader(types.Float, 8))

		// Write the value
		_ = binary.Write(data, binary.LittleEndian, float64(value))
	}
}

func DeserializeFloat(data []byte) (float64, []byte) {

	if len(data) == 0 {
		return 0, data
	}

	var header = data[0]
	data = data[1:]

	var headerType = types.TypeFromHeader(header)
	var length = types.ValueFromHeader(header)

	// Make sure the data is a valid integer value
	if headerType != types.Float {
		return 0, data
	}

	if length != 4 && length != 8 {
		return 0, data
	}

	// Convert byte slice into a reader
	buf := bytes.NewReader(data)

	// Read the value
	// If the value is a float32 then read into float32 then copy to float64
	var value float64 = 0
	if length == 4 {
		var float32Value float32 = 0
		_ = binary.Read(buf, binary.LittleEndian, &float32Value)
		data = data[4:]
		value = float64(float32Value)
	} else {
		_ = binary.Read(buf, binary.LittleEndian, &value)
		data = data[8:]
	}

	return value, data
}
