package types

import (
	"fmt"
)

func TypeNameFromHeader(header byte) string {
	var headerType = TypeFromHeader(header)
	var typeName = TypeName(headerType)
	return typeName
}

func TypeFromHeader(header byte) Types {
	var vartype Types = Types((header & 0xf0) >> 0x04)
	return vartype
}

func ValueFromHeader(header byte) byte {
	var value = header & 0x0f
	return value
}

func CreateHeader(typev Types, value byte) byte {
	var typeInt = byte(typev)
	var header = (typeInt << 4) | (value & 0x0f)
	return header
}

func PrintHeader(data []byte) {
	var header = data[0]
	var typeName = TypeNameFromHeader(header)
	var value = ValueFromHeader(header)

	fmt.Print("{", typeName, ":", value, "}", data[1:])
}
