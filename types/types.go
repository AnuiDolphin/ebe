package types

type Types uint8

const (
	UIntNibble Types = 0
	SInt       Types = 1
	UInt       Types = 2
	Float      Types = 3
	Complex    Types = 4
	Boolean    Types = 5
	String     Types = 6
	Buffer     Types = 7
	Array      Types = 8
	Pointer    Types = 9
	Slice      Types = 10
)

var TypeNames = map[Types]string{
	UIntNibble: "Nibble",
	SInt:       "SInt",
	UInt:       "UInt",
	Float:      "Float",
	Complex:    "Complex",
	Boolean:    "Boolean",
	Pointer:    "Pointer",
	Slice:      "Slice",
	String:     "String",
	Buffer:     "Buffer",
	Array:      "Array",
}

func TypeName(typeValue Types) string {
	return TypeNames[typeValue]
}
