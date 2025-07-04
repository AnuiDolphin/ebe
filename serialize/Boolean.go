package serialize

import (
	"bytes"
	"ebe/types"
)

func SerializeBoolean(value bool, data *bytes.Buffer) {

	// Set the header for the type and put the boolean value in the value nibble
	if value {
		data.WriteByte(types.CreateHeader(types.Boolean, 1))
	} else {
		data.WriteByte(types.CreateHeader(types.Boolean, 0))
	}
}

func DeserializeBoolean(data []byte) (bool, []byte) {

	if len(data) == 0 {
		return false, data
	}

	var header = data[0]
	data = data[1:]

	var headerType = types.TypeFromHeader(header)

	if headerType != types.Boolean {
		return false, data
	}

	var value = types.ValueFromHeader(header)
	return value != 0, data
}
