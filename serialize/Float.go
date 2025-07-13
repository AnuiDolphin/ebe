package serialize

import (
	"ebe-library/types"
	"ebe-library/utils"
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

func serializeFloat(value float64, writer io.Writer) error {

	// If the value fits into a float32, then serialize as a float32
	if value >= -math.SmallestNonzeroFloat32 && value <= math.MaxFloat32 {

		// Write the header as float32
		if err := utils.WriteByte(writer, types.CreateHeader(types.Float, 4)); err != nil {
			return err
		}

		// Write the value
		return binary.Write(writer, binary.LittleEndian, float32(value))
	}

	// Write the header as float64
	if err := utils.WriteByte(writer, types.CreateHeader(types.Float, 8)); err != nil {
		return err
	}

	// Write the value
	return binary.Write(writer, binary.LittleEndian, float64(value))
}

func deserializeFloat(r io.Reader) (float64, error) {
	
	// Read the header using utils.ReadHeader
	headerType, headerValue, err := utils.ReadHeader(r)
	if err != nil {
		return 0, fmt.Errorf("failed to read float header: %w", err)
	}

	length := headerValue

	// Make sure the data is a valid float value
	if headerType != types.Float {
		return 0, fmt.Errorf("expected Float type, got %v", headerType)
	}

	if length != 4 && length != 8 {
		return 0, fmt.Errorf("invalid float length: expected 4 or 8, got %d", length)
	}

	// Read the value
	// If the value is a float32 then read into float32 then copy to float64
	var value float64 = 0
	if length == 4 {
		var float32Value float32 = 0
		err := binary.Read(r, binary.LittleEndian, &float32Value)
		if err != nil {
			return 0, fmt.Errorf("failed to read float32: %w", err)
		}
		value = float64(float32Value)
	} else {
		err := binary.Read(r, binary.LittleEndian, &value)
		if err != nil {
			return 0, fmt.Errorf("failed to read float64: %w", err)
		}
	}

	return value, nil
}
