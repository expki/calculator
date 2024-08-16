package encoding

import (
	"encoding/binary"
	"math"
	"reflect"
)

// Decode implements a custom decoding scheme for the wasm worker
func Decode(encoded []byte) (data any) {
	typeByte := encoded[0]
	switch Type(typeByte) {
	case Type_null:
		return nil
	case Type_bool:
		return encoded[1] == 1
	case Type_uint8:
		return int(encoded[1])
	case Type_int32:
		return int(binary.LittleEndian.Uint32(encoded[1:5]))
	case Type_int64:
		return int(binary.LittleEndian.Uint64(encoded[1:9]))
	case Type_float32:
		return math.Float32frombits(binary.LittleEndian.Uint32(encoded[1:5]))
	case Type_buffer:
		length := int(binary.LittleEndian.Uint32(encoded[1:5]))
		return encoded[5 : 5+length]
	case Type_string:
		length := int(binary.LittleEndian.Uint32(encoded[1:5]))
		return string(encoded[5 : 5+length])
	case Type_array:
		length := int(binary.LittleEndian.Uint32(encoded[1:5]))
		array := make([]any, length)
		offset := 5
		for idx := 0; idx < length; idx++ {
			valueLength := int(binary.LittleEndian.Uint32(encoded[offset : offset+4]))
			array[idx] = Decode(encoded[offset+4 : offset+4+valueLength])
			offset += valueLength + 4
		}
		return array
	case Type_object:
		obj := make(map[string]any)
		maxLength := int(binary.LittleEndian.Uint32(encoded[1:5]))
		offset := 5
		for offset < maxLength+5 {
			keyLength := int(binary.LittleEndian.Uint32(encoded[offset : offset+4]))
			key := string(encoded[offset+4 : offset+4+keyLength])
			offset += keyLength + 4
			valueLength := int(binary.LittleEndian.Uint32(encoded[offset : offset+4]))
			obj[key] = Decode(encoded[offset+4 : offset+4+valueLength])
			offset += 4 + valueLength
		}
		return obj
	default:
		return nil
	}
}

// DecodeIntoStruct implements a custom decoding scheme for the wasm worker
func DecodeIntoStruct(encoded []byte, dst any) (rest any) {
	if reflect.TypeOf(dst).Kind() != reflect.Ptr {
		panic("dst must be a pointer")
	}
	dstStruct := reflect.ValueOf(dst).Elem()
	typeByte := encoded[0]
	switch Type(typeByte) {
	case Type_object:
		obj := make(map[string]any)
		maxLength := int(binary.LittleEndian.Uint32(encoded[1:5]))
		offset := 5
		for offset < maxLength+5 {
			keyLength := int(binary.LittleEndian.Uint32(encoded[offset : offset+4]))
			key := string(encoded[offset+4 : offset+4+keyLength])
			offset += keyLength + 4
			valueLength := int(binary.LittleEndian.Uint32(encoded[offset : offset+4]))
			dstStruct.FieldByName(key).Set(reflect.ValueOf(Decode(encoded[offset+4 : offset+4+valueLength])))
			offset += 4 + valueLength
		}
		return obj
	default:
		return Decode(encoded)
	}
}
