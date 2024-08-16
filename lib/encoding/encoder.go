package encoding

import (
	"encoding/binary"
	"math"
	"reflect"
)

// Encode implements a custom encoding scheme for the wasm worker
func Encode(data any) (encoded []byte) {
	value := reflect.ValueOf(data)
	switch value.Kind() {
	case reflect.Pointer:
		if value.IsNil() {
			return []byte{byte(Type_null)}
		}
		return Encode(value.Elem().Interface())
	case reflect.Bool:
		if value.Bool() {
			return []byte{byte(Type_bool), 1}
		}
		return []byte{byte(Type_bool), 0}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v := value.Int()
		if v > math.MaxInt32 || v < math.MinInt32 {
			// encode int64
			encoded = make([]byte, 9)
			encoded[0] = byte(Type_int64)
			binary.LittleEndian.PutUint64(encoded[1:], uint64(v))
		} else if v > math.MaxUint8 || v < 0 {
			// encode int32
			encoded = make([]byte, 5)
			encoded[0] = byte(Type_int32)
			binary.LittleEndian.PutUint32(encoded[1:], uint32(v))
		} else {
			// encode uint8
			encoded = []byte{byte(Type_uint8), byte(v)}
		}
		return encoded
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v := value.Uint()
		if v > math.MaxInt32 {
			// encode int64
			encoded = make([]byte, 9)
			encoded[0] = byte(Type_int64)
			binary.LittleEndian.PutUint64(encoded[1:], uint64(v))
		} else if v > math.MaxUint8 {
			// encode int32
			encoded = make([]byte, 5)
			encoded[0] = byte(Type_int32)
			binary.LittleEndian.PutUint32(encoded[1:], uint32(v))
		} else {
			// encode uint8
			encoded = []byte{byte(Type_uint8), byte(v)}
		}
		return encoded
	case reflect.Float32, reflect.Float64:
		encoded = make([]byte, 5)
		encoded[0] = byte(Type_float32)
		binary.LittleEndian.PutUint32(encoded[1:], math.Float32bits(float32(value.Float())))
		return encoded
	case reflect.String:
		v := value.String()
		encoded = make([]byte, 5+len(v))
		encoded[0] = byte(Type_string)
		itemLen := uint32(len(v))
		binary.LittleEndian.PutUint32(encoded[1:5], itemLen)
		copy(encoded[5:], v)
		return encoded
	case reflect.Slice:
		if value.IsNil() {
			return []byte{byte(Type_null)}
		}
		fallthrough
	case reflect.Array:
		if value.Type().Elem().Kind() == reflect.Uint8 { // Handle []byte
			v := value.Bytes()
			encoded = make([]byte, 5+len(v))
			encoded[0] = byte(Type_buffer)
			itemLen := uint32(len(v))
			binary.LittleEndian.PutUint32(encoded[1:5], itemLen)
			copy(encoded[5:], v)
			return encoded
		}
		encoded = make([]byte, 5)
		encoded[0] = byte(Type_array)
		binary.LittleEndian.PutUint32(encoded[1:5], uint32(value.Len()))
		for i := 0; i < value.Len(); i++ {
			encoded = append(encoded, arrayValue(value.Index(i).Interface())...)
		}
		return encoded
	case reflect.Struct:
		objType := value.Type()
		data := make([]byte, 0)
		for i := 0; i < value.NumField(); i++ {
			info := objType.Field(i)
			if !info.IsExported() {
				continue
			}
			fieldValue := value.Field(i)
			data = append(data, objValue(info.Name, fieldValue.Interface())...)
		}
		encoded = make([]byte, 5+len(data))
		encoded[0] = byte(Type_object)
		binary.LittleEndian.PutUint32(encoded[1:5], uint32(len(data)))
		copy(encoded[5:], data)
		return encoded
	case reflect.Map:
		data := make([]byte, 0)
		for _, key := range value.MapKeys() {
			data = append(data, objValue(key.String(), value.MapIndex(key).Interface())...)
		}
		encoded = make([]byte, 5+len(data))
		encoded[0] = byte(Type_object)
		binary.LittleEndian.PutUint32(encoded[1:5], uint32(len(data)))
		copy(encoded[5:], data)
		return encoded
	}
	return encoded
}

func arrayValue(item any) []byte {
	itemBytes := Encode(item)
	itemLen := uint32(len(itemBytes))
	itemLenBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(itemLenBytes, itemLen)
	return append(itemLenBytes, itemBytes...)
}

func objValue(key string, value any) []byte {
	keyLen := uint32(len(key))
	keyBytes := make([]byte, 5+keyLen)
	keyBytes[0] = byte(Type_int32)
	binary.LittleEndian.PutUint32(keyBytes[1:5], keyLen)
	copy(keyBytes[5:], key)
	valueBytes := Encode(value)
	valueLen := uint32(len(valueBytes))
	valueLenBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(valueLenBytes, valueLen)
	return append(append(keyBytes, valueLenBytes...), valueBytes...)
}
