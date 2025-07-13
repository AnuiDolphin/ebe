package types

type Types uint8

const (
	UNibble Types = 0
	SNibble Types = 1
	SInt    Types = 2
	UInt    Types = 3
	Float   Types = 4
	Complex Types = 5
	Boolean Types = 6
	String  Types = 7
	Buffer  Types = 8
	Array   Types = 9
	Json    Types = 10
	Map     Types = 11
	Struct  Types = 12
)

var TypeNames = map[Types]string{
	UNibble: "UNibble",
	SNibble: "SNibble",
	SInt:    "SInt",
	UInt:    "UInt",
	Float:   "Float",
	Complex: "Complex",
	Boolean: "Boolean",
	String:  "String",
	Buffer:  "Buffer",
	Array:   "Array",
	Json:    "Json",
	Map:     "Map",
	Struct:  "Struct",
}

func TypeName(typeValue Types) string {
	return TypeNames[typeValue]
}
