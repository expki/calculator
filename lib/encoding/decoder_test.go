package encoding_test

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"math"
	"reflect"
	"testing"

	"github.com/expki/calculator/lib/encoding"
)

// Test_Decode_Null ...
func Test_Decode_Null(t *testing.T) {
	var input []byte = []byte{byte(encoding.Type_null)}
	var want any = nil
	var got any
	got, _ = encoding.Decode(input)
	var output []byte = encoding.Encode(got)

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("Decode(<null>) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(want), want, reflect.TypeOf(got), got)
	}
	if !reflect.DeepEqual(input, output) {
		t.Fatalf("Encode(Decode(<null>)) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(input), input, reflect.TypeOf(output), output)
	}
}

// Test_Decode_Bool ...
func Test_Decode_Bool(t *testing.T) {
	var input []byte = []byte{byte(encoding.Type_bool), 1}
	var want bool = true
	var got any
	got, _ = encoding.Decode(input)
	var output []byte = encoding.Encode(got)

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("Decode(<bool>) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(want), want, reflect.TypeOf(got), got)
	}
	if !reflect.DeepEqual(input, output) {
		t.Fatalf("Encode(Decode(<bool>)) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(input), input, reflect.TypeOf(output), output)
	}
}

// Test_Decode_Uint8 ...
func Test_Decode_Uint8(t *testing.T) {
	var input []byte = []byte{byte(encoding.Type_uint8), math.MaxUint8}
	var want uint8 = math.MaxUint8
	var got any
	got, _ = encoding.Decode(input)
	var output []byte = encoding.Encode(got)

	if !reflect.DeepEqual(int(want), got) {
		t.Fatalf("Decode(<uint8>) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(want), want, reflect.TypeOf(got), got)
	}
	if !reflect.DeepEqual(input, output) {
		t.Fatalf("Encode(Decode(<uint8>)) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(input), input, reflect.TypeOf(output), output)
	}
}

// Test_Decode_Int32 ...
func Test_Decode_Int32(t *testing.T) {
	var input []byte = []byte{byte(encoding.Type_int32), 0, 0, 0, 0}
	binary.LittleEndian.PutUint32(input[1:], uint32(42))
	input[1] |= byte(128)
	var want int32 = -42
	var got any
	got, _ = encoding.Decode(input)
	var output []byte = encoding.Encode(got)

	if !reflect.DeepEqual(int(want), got) {
		t.Fatalf("Decode(<int32>) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(want), want, reflect.TypeOf(got), got)
	}
	if !reflect.DeepEqual(input, output) {
		t.Fatalf("Encode(Decode(<int32>)) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(input), input, reflect.TypeOf(output), output)
	}
}

// Test_Decode_Int64 ...
func Test_Decode_Int64(t *testing.T) {
	var input []byte = []byte{byte(encoding.Type_int64), 0, 0, 0, 0, 0, 0, 0, 0}
	binary.LittleEndian.PutUint64(input[1:], uint64(math.MaxInt32+42))
	input[1] |= byte(128)
	var want int64 = -(math.MaxInt32 + 42)
	var got any
	got, _ = encoding.Decode(input)
	var output []byte = encoding.Encode(got)

	if !reflect.DeepEqual(int(want), got) {
		t.Fatalf("Decode(<int64>) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(want), want, reflect.TypeOf(got), got)
	}
	if !reflect.DeepEqual(input, output) {
		t.Fatalf("Encode(Decode(<int64>)) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(input), input, reflect.TypeOf(output), output)
	}
}

// Test_Decode_Float32 ...
func Test_Decode_Float32(t *testing.T) {
	var input []byte = []byte{byte(encoding.Type_float32), 0, 0, 0, 0}
	binary.LittleEndian.PutUint32(input[1:], math.Float32bits(42))
	var want float32 = 42
	var got any
	got, _ = encoding.Decode(input)
	var output []byte = encoding.Encode(got)

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("Decode(<float32>) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(want), want, reflect.TypeOf(got), got)
	}
	if !reflect.DeepEqual(input, output) {
		t.Fatalf("Encode(Decode(<float32>)) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(input), input, reflect.TypeOf(output), output)
	}
}

// Test_Decode_Buffer ...
func Test_Decode_Buffer(t *testing.T) {
	var input []byte = []byte{byte(encoding.Type_buffer), 0, 0, 0, 0, 0, 1, 2, 3, 4}
	binary.LittleEndian.PutUint32(input[1:5], 5)
	var want []byte = []byte{0, 1, 2, 3, 4}
	var got any
	got, _ = encoding.Decode(input)
	var output []byte = encoding.Encode(got)

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("Decode(<buffer>) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(want), want, reflect.TypeOf(got), got)
	}
	if !reflect.DeepEqual(input, output) {
		t.Fatalf("Encode(Decode(<buffer>)) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(input), input, reflect.TypeOf(output), output)
	}
}

// Test_Decode_String ...
func Test_Decode_String(t *testing.T) {
	var input []byte = []byte{byte(encoding.Type_string), 0, 0, 0, 0, 'h', 'e', 'l', 'l', 'o'}
	binary.LittleEndian.PutUint32(input[1:5], 5)
	var want string = "hello"
	var got any
	got, _ = encoding.Decode(input)
	var output []byte = encoding.Encode(got)

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("Decode(<string>) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(want), want, reflect.TypeOf(got), got)
	}
	if !reflect.DeepEqual(input, output) {
		t.Fatalf("Encode(Decode(<string>)) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(input), input, reflect.TypeOf(output), output)
	}
}

func Test_Decode_Array(t *testing.T) {
	var want []float32 = []float32{42, -42, 42.42, 0}
	var input []byte = []byte{
		byte(encoding.Type_array), 0, 0, 0, 0,
		byte(encoding.Type_float32), 0, 0, 0, 0,
		byte(encoding.Type_float32), 0, 0, 0, 0,
		byte(encoding.Type_float32), 0, 0, 0, 0,
		byte(encoding.Type_float32), 0, 0, 0, 0,
	}
	binary.LittleEndian.PutUint32(input[1:5], uint32(len(want)))
	binary.LittleEndian.PutUint32(input[6:10], math.Float32bits(want[0]))
	binary.LittleEndian.PutUint32(input[11:15], math.Float32bits(want[1]))
	binary.LittleEndian.PutUint32(input[16:20], math.Float32bits(want[2]))
	binary.LittleEndian.PutUint32(input[21:25], math.Float32bits(want[3]))
	var got any
	got, _ = encoding.Decode(input)
	var output []byte = encoding.Encode(got)

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("Decode(<array>) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(want), want, reflect.TypeOf(got), got)
	}
	if !reflect.DeepEqual(input, output) {
		t.Fatalf("Encode(Decode(<array>)) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(input), input, reflect.TypeOf(output), output)
	}
}

// Test_Decode_Map ...
func Test_Decode_Map(t *testing.T) {
	var want map[string]any = map[string]any{
		"A": 42,
		"B": []int{42, -42},
		"C": map[string]any{
			"D": 42.42,
			"E": "hello",
		},
	}
	var input []byte = []byte{
		byte(encoding.Type_object), 3, 0, 0, 0,
		// A
		1, 0, 0, 0, 'A', byte(encoding.Type_uint8), 42,
		// B
		1, 0, 0, 0, 'B', byte(encoding.Type_array), 2, 0, 0, 0,
		byte(encoding.Type_uint8), 42,
		byte(encoding.Type_int32), 0, 0, 0, 0,
		// C
		1, 0, 0, 0, 'C', byte(encoding.Type_object), 2, 0, 0, 0,
		// D
		1, 0, 0, 0, 'D', byte(encoding.Type_float32), 0, 0, 0, 0,
		// E
		1, 0, 0, 0, 'E', byte(encoding.Type_string), 5, 0, 0, 0,
		'h', 'e', 'l', 'l', 'o',
	}
	binary.LittleEndian.PutUint32(input[25:29], 42)
	input[25] |= byte(128)
	binary.LittleEndian.PutUint32(input[45:49], math.Float32bits(42.42))

	var got any
	got, _ = encoding.Decode(input)
	var output []byte = encoding.Encode(got)

	wantJson, _ := json.Marshal(want)
	gotJson, _ := json.Marshal(got)
	if !reflect.DeepEqual(string(wantJson), string(gotJson)) {
		t.Fatalf("Decode(<object>) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(want), want, reflect.TypeOf(got), got)
	}
	inputCount := make(map[uint8]int)
	outputCount := make(map[uint8]int)
	for i := 0; i < len(input); i++ {
		inputCount[input[i]]++
	}
	for i := 0; i < len(output); i++ {
		outputCount[output[i]]++
	}
	if !reflect.DeepEqual(inputCount, outputCount) {
		t.Fatalf("Encode(Decode(<object>)) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(input), input, reflect.TypeOf(output), output)
	}
}

// Test_Decode_Map ...
func Test_DecodeIO_Map(t *testing.T) {
	var want map[string]any = map[string]any{
		"A": 42,
		"B": []int{42, -42},
		"C": map[string]any{
			"D": 42.42,
			"E": "hello",
		},
	}
	var input []byte = []byte{
		byte(encoding.Type_object), 3, 0, 0, 0,
		// A
		1, 0, 0, 0, 'A', byte(encoding.Type_uint8), 42,
		// B
		1, 0, 0, 0, 'B', byte(encoding.Type_array), 2, 0, 0, 0,
		byte(encoding.Type_uint8), 42,
		byte(encoding.Type_int32), 0, 0, 0, 0,
		// C
		1, 0, 0, 0, 'C', byte(encoding.Type_object), 2, 0, 0, 0,
		// D
		1, 0, 0, 0, 'D', byte(encoding.Type_float32), 0, 0, 0, 0,
		// E
		1, 0, 0, 0, 'E', byte(encoding.Type_string), 5, 0, 0, 0,
		'h', 'e', 'l', 'l', 'o',
	}
	binary.LittleEndian.PutUint32(input[25:29], 42)
	input[25] |= byte(128)
	binary.LittleEndian.PutUint32(input[45:49], math.Float32bits(42.42))

	var got any
	got, _ = encoding.DecodeIO(bytes.NewReader(input))
	var output []byte = encoding.Encode(got)

	wantJson, _ := json.Marshal(want)
	gotJson, _ := json.Marshal(got)
	if !reflect.DeepEqual(string(wantJson), string(gotJson)) {
		t.Fatalf("Decode(<object>) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(want), want, reflect.TypeOf(got), got)
	}
	inputCount := make(map[uint8]int)
	outputCount := make(map[uint8]int)
	for i := 0; i < len(input); i++ {
		inputCount[input[i]]++
	}
	for i := 0; i < len(output); i++ {
		outputCount[output[i]]++
	}
	if !reflect.DeepEqual(inputCount, outputCount) {
		t.Fatalf("Encode(Decode(<object>)) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(input), input, reflect.TypeOf(output), output)
	}
}
