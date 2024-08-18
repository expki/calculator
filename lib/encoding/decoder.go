package encoding

import (
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
)

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
		fmt.Println("Type_array")
		length := int(binary.LittleEndian.Uint32(encoded[1:5]))
		fmt.Println("Array length:", length)
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
