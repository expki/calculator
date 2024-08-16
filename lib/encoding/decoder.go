package encoding

import (
	"encoding/binary"
	"errors"
)

// Decode implements a custom decoding scheme for the wasm worker
func Decode(encoded []byte) (data any, err error) {
	if len(encoded) == 0 {
		return nil, errors.New("empty input")
	}

	typeByte := encoded[0]
	switch typeByte {
	case byte(Type_null):
		return nil, nil
	case byte(Type_bool):
		if len(encoded) < 2 {
			return nil, errors.New("invalid bool encoding")
		}
		return encoded[1] == 1, nil
	case byte(Type_int32):
		if len(encoded) < 5 {
			return nil, errors.New("invalid int32 encoding")
		}
		v := int32(binary.LittleEndian.Uint32(encoded[1:5]))
		return v, nil
	case byte(Type_int64):
		if len(encoded) < 9 {
			return nil, errors.New("invalid int64 encoding")
		}
		v := int64(binary.LittleEndian.Uint64(encoded[1:9]))
		return v, nil
	default:
		return nil, errors.New("unknown type")
	}
}
