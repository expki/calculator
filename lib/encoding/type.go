package encoding

type Type uint8

const (
	Type_null Type = iota
	Type_bool
	Type_uint8
	Type_int32
	Type_int64
	Type_float32
	Type_buffer
	Type_string
	Type_array
	Type_object
)
