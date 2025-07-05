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
	Pointer Types = 10
	Slice   Types = 11
)

var TypeNames = map[Types]string{
	UNibble: "UNibble",
	SNibble: "SNibble",
	SInt:    "SInt",
	UInt:    "UInt",
	Float:   "Float",
	Complex: "Complex",
	Boolean: "Boolean",
	Pointer: "Pointer",
	Slice:   "Slice",
	String:  "String",
	Buffer:  "Buffer",
	Array:   "Array",
}

func TypeName(typeValue Types) string {
	return TypeNames[typeValue]
}
