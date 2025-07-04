package main

import (
	"bytes"
	"ebe/serialize"
	"ebe/types"
	"fmt"
	"math"
)

func main() {
	TestUint(0)
	TestUint(1)
	TestUint(0xff)
	TestUint(0xffff)
	TestUint(0xffffff)
	TestUint(0xffffffff)
	TestUint(0xffffffffff)
	TestUint(0xffffffffffff)
	TestUint(0xffffffffffffff)
	TestUint(0xffffffffffffffff)

	TestUint(0xff + 1)
	TestUint(0xffff + 1)
	TestUint(0xffffff + 1)
	TestUint(0xffffffff + 1)
	TestUint(0xffffffffff + 1)
	TestUint(0xffffffffffff + 1)

	TestSint(-0x7f)
	TestSint(-0x7fff)
	TestSint(-0x7fffff)
	TestSint(-0x7fffffff)
	TestSint(-0x7fffffffff)
	TestSint(-0x7fffffffffff)
	TestSint(-0x7fffffffffffff)
	TestSint(-0x7fffffffffffffff)

	TestSint(-0xff)
	TestSint(-0xffff)
	TestSint(-0xffffff)
	TestSint(-0xffffffff)
	TestSint(-0xffffffffff)
	TestSint(-0xffffffffffff)
	TestSint(-0xffffffffffffff)
	TestSint(-0x7fffffffffffffff)

	TestString("")
	TestString("A")
	TestString("AB")
	TestString("Hello")
	TestString("The quick brown fox jumped over the lazy dog")
	TestString("🙂")

	TestBuffer([]byte{})
	TestBuffer([]byte{1})
	TestBuffer([]byte{1, 2, 3})
	TestBuffer([]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17})

	TestBoolean(true)
	TestBoolean(false)

	TestFloat(0)
	TestFloat(0.1)
	TestFloat(-1.1)
	TestFloat(1.17549435e-38) // Min float32
	TestFloat(1.17549435e-39)
	TestFloat(math.MaxFloat32)
	TestFloat(math.MaxFloat64)
}

func TestUint(value uint64) {
	var data bytes.Buffer

	serialize.SerializeUint(value, &data)
	fmt.Print("value = ", value, ", ")
	types.PrintHeader(data.Bytes())

	var readValue, _ = serialize.DeserializeUint(data.Bytes())
	fmt.Print(", Read:", readValue)

	if readValue == value {
		fmt.Println(", Pass")
	} else {
		fmt.Println(", *** FAIL ***")
	}
}

func TestSint(value int64) {
	var data bytes.Buffer

	serialize.SerializeSint(value, &data)
	fmt.Print("value = ", value, ", ")
	types.PrintHeader(data.Bytes())

	var readValue, _ = serialize.DeserializeSint(data.Bytes())
	fmt.Print(", Read:", readValue)

	if readValue == value {
		fmt.Println(", Pass")
	} else {
		fmt.Println(", *** FAIL ***")
	}
}

func TestString(value string) {
	var data bytes.Buffer

	serialize.SerializeString(value, &data)
	fmt.Print("value = '", value, "', ")
	types.PrintHeader(data.Bytes())

	var readValue, _ = serialize.DeserializeString(data.Bytes())
	fmt.Print(", Read:'", readValue, "'")

	if readValue == value {
		fmt.Println(", Pass")
	} else {
		fmt.Println(", *** FAIL ***")
	}
}

func TestBuffer(value []byte) {
	var data bytes.Buffer

	serialize.SerializeBuffer(value, &data)
	fmt.Print("value: '", value, "', ")
	types.PrintHeader(data.Bytes())

	var readValue, _ = serialize.DeserializeBuffer(data.Bytes())
	fmt.Print(", Read:", readValue.Bytes())

	if bytes.Equal(readValue.Bytes(), value) {
		fmt.Println(", Pass")
	} else {
		fmt.Println(", *** FAIL ***")
	}
}

func TestBoolean(value bool) {
	var data bytes.Buffer

	serialize.SerializeBoolean(value, &data)
	fmt.Print("value: '", value, "', ")
	types.PrintHeader(data.Bytes())

	var readValue, _ = serialize.DeserializeBoolean(data.Bytes())
	fmt.Print(", Read:", readValue)

	if value == readValue {
		fmt.Println(", Pass")
	} else {
		fmt.Println(", *** FAIL ***")
	}
}

func TestFloat(value float64) {
	var data bytes.Buffer

	serialize.SerializeFloat(value, &data)
	fmt.Print("value = ", value, ", ")
	types.PrintHeader(data.Bytes())

	var readValue, _ = serialize.DeserializeFloat(data.Bytes())
	fmt.Print(", Read:", readValue)

	if readValue == value {
		fmt.Println(", Pass")
	} else {
		fmt.Println(", *** FAIL ***")
	}
}
