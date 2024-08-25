package encoding

import (
	"encoding/binary"
	"math"
	"reflect"

	"github.com/expki/calculator/lib/compression"
)

func EncodeWithCompression(data any) (encoded []byte) {
	payload := Encode(data)
	payloadCompressed := compression.Compress(payload)
	if len(payloadCompressed) < len(payload) || true {
		encoded = make([]byte, 1+len(payloadCompressed))
		encoded[0] = 1
		copy(encoded[1:], payloadCompressed)
	} else {
		encoded = make([]byte, 1+len(payload))
		encoded[0] = 0
		copy(encoded[1:], payload)
	}
	return encoded
}

// Encode implements a custom encoding scheme for the wasm worker
func Encode(data any) (encoded []byte) {
	value := reflect.ValueOf(data)
	switch value.Kind() {
	case reflect.Invalid:
		return []byte{byte(Type_null)}
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
		var v uint64
		var isNegative bool
		vs := value.Int()
		if vs < 0 {
			v = uint64(-vs)
			isNegative = true
		} else {
			v = uint64(vs)
		}
		if v > math.MaxInt32 {
			// encode int64
			encoded = make([]byte, 9)
			encoded[0] = byte(Type_int64)
			binary.LittleEndian.PutUint64(encoded[1:], v)
		} else if isNegative || v > math.MaxUint8 {
			// encode int32
			encoded = make([]byte, 5)
			encoded[0] = byte(Type_int32)
			binary.LittleEndian.PutUint32(encoded[1:], uint32(v))
		} else {
			// encode uint8
			encoded = []byte{byte(Type_uint8), byte(v)}
		}
		if isNegative {
			encoded[1] |= byte(128)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v := value.Uint()
		if v > math.MaxInt32 {
			// encode int64
			encoded = make([]byte, 9)
			encoded[0] = byte(Type_int64)
			binary.LittleEndian.PutUint64(encoded[1:], v)
		} else if v > math.MaxUint8 {
			// encode int32
			encoded = make([]byte, 5)
			encoded[0] = byte(Type_int32)
			binary.LittleEndian.PutUint32(encoded[1:], uint32(v))
		} else {
			// encode uint8
			encoded = []byte{byte(Type_uint8), byte(v)}
		}
	case reflect.Float32, reflect.Float64:
		encoded = make([]byte, 5)
		encoded[0] = byte(Type_float32)
		binary.LittleEndian.PutUint32(encoded[1:], math.Float32bits(float32(value.Float())))
	case reflect.String:
		v := value.String()
		encoded = make([]byte, 5+len(v))
		encoded[0] = byte(Type_string)
		itemLen := uint32(len(v))
		binary.LittleEndian.PutUint32(encoded[1:5], itemLen)
		copy(encoded[5:], v)
	case reflect.Slice:
		if value.IsNil() {
			return []byte{byte(Type_null)}
		}
		fallthrough
	case reflect.Array:
		length := value.Len()
		if length == 0 {
			// unfortunately the decoder can't identify the golang type of an empty slice
			return []byte{byte(Type_null)}
		}
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
		binary.LittleEndian.PutUint32(encoded[1:5], uint32(length))
		for i := 0; i < length; i++ {
			value := value.Index(i).Interface()
			encodedValue := Encode(value)
			encoded = append(encoded, encodedValue...)
		}
	case reflect.Struct:
		objType := value.Type()
		data := make([]byte, 0)
		var count uint32 = 0
		for idx := 0; idx < value.NumField(); idx++ {
			info := objType.Field(idx)
			if !info.IsExported() {
				continue
			}
			fieldValue := value.Field(idx).Interface()
			data = append(data, objValue(info.Name, fieldValue)...)
			count++
		}
		encoded = make([]byte, 5+len(data))
		encoded[0] = byte(Type_object)
		binary.LittleEndian.PutUint32(encoded[1:5], count)
		copy(encoded[5:], data)
	case reflect.Map:
		data := make([]byte, 0)
		var count uint32
		for _, key := range value.MapKeys() {
			data = append(data, objValue(key.String(), value.MapIndex(key).Interface())...)
			count++
		}
		encoded = make([]byte, 5+len(data))
		encoded[0] = byte(Type_object)
		binary.LittleEndian.PutUint32(encoded[1:5], count)
		copy(encoded[5:], data)
	}
	return encoded
}

func objValue(key string, value any) []byte {
	data := Encode(value)
	dataLen := uint32(len(data))
	keyLen := uint32(len(key))
	bytes := make([]byte, 4+keyLen+dataLen)
	// key
	binary.LittleEndian.PutUint32(bytes[0:4], keyLen)
	copy(bytes[4:4+keyLen], key)
	// value
	copy(bytes[4+keyLen:4+keyLen+dataLen], data)
	return bytes
}
