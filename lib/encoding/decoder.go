package encoding

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"

	"github.com/expki/calculator/lib/compression"
)

func DecodeWithCompression(encoded []byte) (data any, err error) {
	var msg []byte
	if encoded[0] == 1 {
		msg, err = compression.Decompress(encoded[1:])
		if err != nil {
			return nil, err
		}
	} else {
		msg = encoded[1:]
	}
	data, _ = Decode(msg)
	return data, nil
}

// Decode ...
func Decode(encoded []byte) (data any, length int) {
	typeByte := encoded[0]
	switch Type(typeByte) {
	case Type_null:
		return nil, 1
	case Type_bool:
		return encoded[1] == 1, 2
	case Type_uint8:
		return int(encoded[1]), 2
	case Type_int32:
		if encoded[1] >= 128 {
			clone := make([]byte, 4)
			copy(clone, encoded[1:5])
			clone[0] ^= 128
			return -int(binary.LittleEndian.Uint32(clone)), 5
		}
		return int(binary.LittleEndian.Uint32(encoded[1:5])), 5
	case Type_int64:
		if encoded[1] >= 128 {
			clone := make([]byte, 8)
			copy(clone, encoded[1:9])
			clone[0] ^= 128
			return -int(binary.LittleEndian.Uint64(clone)), 9
		}
		return int(binary.LittleEndian.Uint64(encoded[1:9])), 9
	case Type_float32:
		return math.Float32frombits(binary.LittleEndian.Uint32(encoded[1:5])), 5
	case Type_buffer:
		length := int(binary.LittleEndian.Uint32(encoded[1:5]))
		return encoded[5 : 5+length], 5 + length
	case Type_string:
		length := int(binary.LittleEndian.Uint32(encoded[1:5]))
		return string(encoded[5 : 5+length]), 5 + length
	case Type_array:
		length := int(binary.LittleEndian.Uint32(encoded[1:5]))
		array := make([]any, length)
		offset := 5
		var valueLength int
		for idx := 0; idx < length; idx++ {
			array[idx], valueLength = Decode(encoded[offset:])
			if valueLength == 0 {
				panic(fmt.Errorf("unable to decode array item %d", idx))
			}
			offset += valueLength
		}
		return setSliceType(array), offset
	case Type_object:
		obj := make(map[string]any)
		count := int(binary.LittleEndian.Uint32(encoded[1:5]))
		offset := 5
		var valueLength int
		for idx := 0; idx < count; idx++ {
			keyLength := int(binary.LittleEndian.Uint32(encoded[offset : offset+4]))
			offset += 4
			key := string(encoded[offset : offset+keyLength])
			offset += keyLength
			obj[key], valueLength = Decode(encoded[offset:])
			if valueLength == 0 {
				panic(fmt.Errorf("unable to decode object item %s", key))
			}
			offset += valueLength
		}
		return obj, offset
	default:
		return nil, 0
	}
}

func setSliceType(unknown []any) (known any) {
	// Compatible with any slice type
	if len(unknown) == 0 {
		return nil
	}
	// Confirm all has the same type
	firstType := reflect.TypeOf(unknown[0])
	for idx := 1; idx < len(unknown); idx++ {
		if reflect.TypeOf(unknown[idx]) != firstType {
			return unknown
		}
	}
	// Assign all as first type
	output := reflect.MakeSlice(reflect.SliceOf(firstType), 0, len(unknown))
	for _, item := range unknown {
		output = reflect.Append(output, reflect.ValueOf(item))
	}
	return output.Interface()
}

// DecodeIO ...
func DecodeIO(in io.Reader) (data any, length int) {
	encoded := make([]byte, 1)
	if _, err := in.Read(encoded); err != nil {
		return nil, 0
	}
	typeByte := encoded[0]
	switch Type(typeByte) {
	case Type_null:
		return nil, 1
	case Type_bool:
		encoded := make([]byte, 1)
		if _, err := in.Read(encoded); err != nil {
			return nil, 0
		}
		return encoded[0] == 1, 2
	case Type_uint8:
		encoded := make([]byte, 1)
		if _, err := in.Read(encoded); err != nil {
			return nil, 0
		}
		return int(encoded[0]), 2
	case Type_int32:
		encoded := make([]byte, 4)
		if _, err := in.Read(encoded); err != nil {
			return nil, 0
		}
		if encoded[0] >= 128 {
			clone := make([]byte, 4)
			copy(clone, encoded[0:4])
			clone[0] ^= 128
			return -int(binary.LittleEndian.Uint32(clone)), 5
		}
		return int(binary.LittleEndian.Uint32(encoded[1:5])), 5
	case Type_int64:
		encoded := make([]byte, 8)
		if _, err := in.Read(encoded); err != nil {
			return nil, 0
		}
		if encoded[0] >= 128 {
			clone := make([]byte, 8)
			copy(clone, encoded[0:8])
			clone[0] ^= 128
			return -int(binary.LittleEndian.Uint64(clone)), 9
		}
		return int(binary.LittleEndian.Uint64(encoded[0:8])), 9
	case Type_float32:
		encoded := make([]byte, 4)
		if _, err := in.Read(encoded); err != nil {
			return nil, 0
		}
		return math.Float32frombits(binary.LittleEndian.Uint32(encoded[0:4])), 5
	case Type_buffer:
		encoded := make([]byte, 4)
		if _, err := in.Read(encoded); err != nil {
			return nil, 0
		}
		length := int(binary.LittleEndian.Uint32(encoded[0:4]))
		encoded = make([]byte, length)
		if _, err := in.Read(encoded); err != nil {
			return nil, 0
		}
		return encoded[0:length], 5 + length
	case Type_string:
		encoded := make([]byte, 4)
		if _, err := in.Read(encoded); err != nil {
			return nil, 0
		}
		length := int(binary.LittleEndian.Uint32(encoded[0:4]))
		encoded = make([]byte, length)
		if _, err := in.Read(encoded); err != nil {
			return nil, 0
		}
		return string(encoded[0:length]), 5 + length
	case Type_array:
		encoded := make([]byte, 4)
		if _, err := in.Read(encoded); err != nil {
			return nil, 0
		}
		length := int(binary.LittleEndian.Uint32(encoded[0:4]))
		array := make([]any, length)
		offset := 5
		var valueLength int
		for idx := 0; idx < length; idx++ {
			array[idx], valueLength = DecodeIO(in)
			if valueLength == 0 {
				panic(fmt.Errorf("unable to decode array item %d", idx))
			}
			offset += valueLength
		}
		return setSliceType(array), offset
	case Type_object:
		obj := make(map[string]any)
		encoded := make([]byte, 4)
		if _, err := in.Read(encoded); err != nil {
			return nil, 0
		}
		count := int(binary.LittleEndian.Uint32(encoded[0:4]))
		offset := 5
		var valueLength int
		for idx := 0; idx < count; idx++ {
			encoded := make([]byte, 4)
			if _, err := in.Read(encoded); err != nil {
				return nil, 0
			}
			keyLength := int(binary.LittleEndian.Uint32(encoded[0:4]))
			offset += 4
			encoded = make([]byte, keyLength)
			if _, err := in.Read(encoded); err != nil {
				return nil, 0
			}
			key := string(encoded)
			offset += keyLength
			obj[key], valueLength = DecodeIO(in)
			if valueLength == 0 {
				panic(fmt.Errorf("unable to decode object item %s", key))
			}
			offset += valueLength
		}
		return obj, offset
	default:
		return nil, 0
	}
}
