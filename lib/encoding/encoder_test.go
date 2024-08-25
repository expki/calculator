package encoding_test

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/expki/calculator/lib/encoding"
)

// Test_Encode_Null ...
func Test_Encode_Null(t *testing.T) {
	var input any = nil
	var want []byte = []byte{byte(encoding.Type_null)}
	var got []byte = encoding.Encode(input)
	var output any
	output, _ = encoding.Decode(got)

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("Encode(<null>) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(want), want, reflect.TypeOf(got), got)
	}
	if !reflect.DeepEqual(input, output) {
		t.Fatalf("Decode(Encode(<null>)) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(input), input, reflect.TypeOf(output), output)
	}
}

// Test_Encode_Bool ...
func Test_Encode_Bool(t *testing.T) {
	var input bool = true
	var want []byte = []byte{byte(encoding.Type_bool), 1}
	var got []byte = encoding.Encode(input)
	var output any
	output, _ = encoding.Decode(got)

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("Encode(<bool>) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(want), want, reflect.TypeOf(got), got)
	}
	if !reflect.DeepEqual(input, output) {
		t.Fatalf("Decode(Encode(<bool>)) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(input), input, reflect.TypeOf(output), output)
	}
}

// Test_Encode_Uint8 ...
func Test_Encode_Uint8(t *testing.T) {
	var input uint8 = 42
	var want []byte = []byte{byte(encoding.Type_uint8), 42}
	var got []byte = encoding.Encode(input)
	var output any
	output, _ = encoding.Decode(got)

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("Encode(<uint8>) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(want), want, reflect.TypeOf(got), got)
	}
	if !reflect.DeepEqual(int(input), output) {
		t.Fatalf("Decode(Encode(<uint8>)) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(input), input, reflect.TypeOf(output), output)
	}
}

// Test_Encode_Int32 ...
func Test_Encode_Int32(t *testing.T) {
	var input int32 = -42
	var want []byte = []byte{byte(encoding.Type_int32), 0, 0, 0, 0}
	binary.LittleEndian.PutUint32(want[1:], uint32(42))
	want[1] |= byte(128)
	var got []byte = encoding.Encode(input)
	var output any
	output, _ = encoding.Decode(got)

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("Encode(<int32>) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(want), want, reflect.TypeOf(got), got)
	}
	if !reflect.DeepEqual(int(input), output) {
		t.Fatalf("Decode(Encode(<int32>)) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(input), input, reflect.TypeOf(output), output)
	}
}

// Test_Encode_Int64 ...
func Test_Encode_Int64(t *testing.T) {
	var input int64 = -(math.MaxInt32 + 42)
	var want []byte = []byte{byte(encoding.Type_int64), 0, 0, 0, 0, 0, 0, 0, 0}
	binary.LittleEndian.PutUint64(want[1:], uint64(math.MaxInt32+42))
	want[1] |= byte(128)
	var got []byte = encoding.Encode(input)
	var output any
	output, _ = encoding.Decode(got)

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("Encode(<int64>) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(want), want, reflect.TypeOf(got), got)
	}
	if !reflect.DeepEqual(int(input), output) {
		t.Fatalf("Decode(Encode(<int64>)) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(input), input, reflect.TypeOf(output), output)
	}
}

// Test_Encode_Float32 ...
func Test_Encode_Float32(t *testing.T) {
	var input float32 = 42.42
	var want []byte = []byte{byte(encoding.Type_float32), 0, 0, 0, 0}
	binary.LittleEndian.PutUint32(want[1:], math.Float32bits(42.42))
	var got []byte = encoding.Encode(input)
	var output any
	output, _ = encoding.Decode(got)

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("Encode(<float32>) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(want), want, reflect.TypeOf(got), got)
	}
	if !reflect.DeepEqual(float32(input), output) {
		t.Fatalf("Decode(Encode(<float32>)) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(input), input, reflect.TypeOf(output), output)
	}
}

// Test_Encode_Buffer ...
func Test_Encode_Buffer(t *testing.T) {
	var input []byte = []byte{0, 1, 2, 3, 4}
	var want []byte = []byte{byte(encoding.Type_buffer), 0, 0, 0, 0, 0, 1, 2, 3, 4}
	binary.LittleEndian.PutUint32(want[1:5], uint32(len(input)))
	var got []byte = encoding.Encode(input)
	var output any
	output, _ = encoding.Decode(got)

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("Encode(<buffer>) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(want), want, reflect.TypeOf(got), got)
	}
	if !reflect.DeepEqual(input, output) {
		t.Fatalf("Decode(Encode(<buffer>)) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(input), input, reflect.TypeOf(output), output)
	}
}

// Test_Encode_String ...
func Test_Encode_String(t *testing.T) {
	var input string = "hello"
	var want []byte = []byte{byte(encoding.Type_string), 0, 0, 0, 0, 'h', 'e', 'l', 'l', 'o'}
	binary.LittleEndian.PutUint32(want[1:5], uint32(len(input)))
	var got []byte = encoding.Encode(input)
	var output any
	output, _ = encoding.Decode(got)

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("Encode(<string>) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(want), want, reflect.TypeOf(got), got)
	}
	if !reflect.DeepEqual(input, output) {
		t.Fatalf("Decode(Encode(<string>)) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(input), input, reflect.TypeOf(output), output)
	}
}

// Test_Encode_Array ...
func Test_Encode_Array(t *testing.T) {
	var input []float32 = []float32{42, -42, 42.42, 0}
	var want []byte = []byte{
		byte(encoding.Type_array), 0, 0, 0, 0,
		byte(encoding.Type_float32), 0, 0, 0, 0,
		byte(encoding.Type_float32), 0, 0, 0, 0,
		byte(encoding.Type_float32), 0, 0, 0, 0,
		byte(encoding.Type_float32), 0, 0, 0, 0,
	}
	binary.LittleEndian.PutUint32(want[1:5], uint32(len(input)))
	binary.LittleEndian.PutUint32(want[6:10], math.Float32bits(input[0]))
	binary.LittleEndian.PutUint32(want[11:15], math.Float32bits(input[1]))
	binary.LittleEndian.PutUint32(want[16:20], math.Float32bits(input[2]))
	binary.LittleEndian.PutUint32(want[21:25], math.Float32bits(input[3]))
	var got []byte = encoding.Encode(input)
	var output any
	output, _ = encoding.Decode(got)

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("Encode(<array>) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(want), want, reflect.TypeOf(got), got)
	}
	if !reflect.DeepEqual(input, output) {
		t.Fatalf("Decode(Encode(<array>)) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(input), input, reflect.TypeOf(output), output)
	}
}

// Test_Encode_Map ...
func Test_Encode_Map(t *testing.T) {
	var input map[string]any = map[string]any{
		"A": 42,
		"B": []int{42, -42},
		"C": map[string]any{
			"D": 42.42,
			"E": "hello",
		},
	}
	var want []byte = []byte{
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
	binary.LittleEndian.PutUint32(want[25:29], 42)
	want[25] |= byte(128)
	binary.LittleEndian.PutUint32(want[45:49], math.Float32bits(42.42))

	var got []byte = encoding.Encode(input)
	var output any
	output, _ = encoding.Decode(got)

	wantCount := make(map[uint8]int)
	gotCount := make(map[uint8]int)
	for i := 0; i < len(want); i++ {
		wantCount[want[i]]++
	}
	for i := 0; i < len(got); i++ {
		gotCount[got[i]]++
	}
	if !reflect.DeepEqual(wantCount, gotCount) {
		t.Fatalf("Encode(<map>) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(want), want, reflect.TypeOf(got), got)
	}
	inputJson, _ := json.Marshal(input)
	outputJson, _ := json.Marshal(output)
	if !reflect.DeepEqual(string(inputJson), string(outputJson)) {
		t.Fatalf("Decode(Encode(<map>)) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(input), input, reflect.TypeOf(output), output)
	}
}

// Test_Encode_Struct ...
func Test_Encode_Struct(t *testing.T) {
	type object struct {
		a bool
		A int
		B []int
		C struct {
			D float32
			E string
		}
	}
	var input object = object{
		a: true,
		A: 42,
		B: []int{42, -42},
		C: struct {
			D float32
			E string
		}{
			D: 42.42,
			E: "hello",
		},
	}
	var want []byte = []byte{
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
	binary.LittleEndian.PutUint32(want[25:29], 42)
	want[25] |= byte(128)
	binary.LittleEndian.PutUint32(want[45:49], math.Float32bits(42.42))

	var got []byte = encoding.Encode(input)
	var output any
	output, _ = encoding.Decode(got)

	wantCount := make(map[uint8]int)
	gotCount := make(map[uint8]int)
	for i := 0; i < len(want); i++ {
		wantCount[want[i]]++
	}
	for i := 0; i < len(got); i++ {
		gotCount[got[i]]++
	}
	if !reflect.DeepEqual(wantCount, gotCount) {
		t.Fatalf("Encode(<struct>) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(want), want, reflect.TypeOf(got), got)
	}
	inputJson, _ := json.Marshal(input)
	outputJson, _ := json.Marshal(output)
	if !reflect.DeepEqual(inputJson, outputJson) {
		t.Fatalf("Decode(Encode(<struct>)) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(input), input, reflect.TypeOf(output), output)
	}
}

// Test_Encode_Struct2 ...
func Test_Encode_Struct2(t *testing.T) {
	type object struct {
		A struct {
			Aa int
		}
		B struct {
			Ba int
		}
	}
	var input object = object{
		A: struct{ Aa int }{
			Aa: 1,
		},
		B: struct{ Ba int }{
			Ba: 2,
		},
	}
	var want []byte = []byte{
		byte(encoding.Type_object), 2, 0, 0, 0,
		// A
		1, 0, 0, 0, 'A', byte(encoding.Type_object), 1, 0, 0, 0,
		2, 0, 0, 0, 'A', 'a', byte(encoding.Type_uint8), 1,
		// B
		1, 0, 0, 0, 'B', byte(encoding.Type_object), 1, 0, 0, 0,
		2, 0, 0, 0, 'B', 'a', byte(encoding.Type_uint8), 2,
	}

	var got []byte = encoding.Encode(input)
	var output any
	output, _ = encoding.Decode(got)

	wantCount := make(map[uint8]int)
	gotCount := make(map[uint8]int)
	for i := 0; i < len(want); i++ {
		wantCount[want[i]]++
	}
	for i := 0; i < len(got); i++ {
		gotCount[got[i]]++
	}
	if !reflect.DeepEqual(wantCount, gotCount) {
		t.Fatalf("Encode(<struct>) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(want), want, reflect.TypeOf(got), got)
	}
	inputJson, _ := json.Marshal(input)
	outputJson, _ := json.Marshal(output)
	if !reflect.DeepEqual(inputJson, outputJson) {
		t.Fatalf("Decode(Encode(<struct>)) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(input), input, reflect.TypeOf(output), output)
	}
}

// Test_Encode_Struct3 ...
func Test_Encode_Struct3(t *testing.T) {
	type object struct {
		A []int
		B []int
	}
	var input object = object{
		B: []int{0},
	}
	var want []byte = []byte{
		byte(encoding.Type_object), 2, 0, 0, 0,
		// A
		1, 0, 0, 0, 'A', byte(encoding.Type_null),
		// B
		1, 0, 0, 0, 'B', byte(encoding.Type_array), 1, 0, 0, 0,
		byte(encoding.Type_uint8), 0,
	}
	var got []byte = encoding.Encode(input)
	var output any
	output, _ = encoding.Decode(got)

	wantCount := make(map[uint8]int)
	gotCount := make(map[uint8]int)
	for i := 0; i < len(want); i++ {
		wantCount[want[i]]++
	}
	for i := 0; i < len(got); i++ {
		gotCount[got[i]]++
	}
	if !reflect.DeepEqual(wantCount, gotCount) {
		t.Fatalf("Encode(<struct>) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(want), want, reflect.TypeOf(got), got)
	}
	inputJson, _ := json.Marshal(input)
	outputJson, _ := json.Marshal(output)
	if !reflect.DeepEqual(inputJson, outputJson) {
		t.Fatalf("Decode(Encode(<struct>)) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(input), input, reflect.TypeOf(output), output)
	}
}

// Test_Encode_Struct_Compression ...
func Test_Encode_Struct_Compression(t *testing.T) {
	type object struct {
		a bool
		A int
		B []int
		C struct {
			D float32
			E string
		}
	}
	var input object = object{
		a: true,
		A: 42,
		B: []int{42, -42},
		C: struct {
			D float32
			E string
		}{
			D: 42.42,
			E: "hello",
		},
	}

	var got []byte = encoding.EncodeWithCompression(input)
	fmt.Println("compress len:", len(got))
	output, err := encoding.DecodeWithCompression(got)
	if err != nil {
		t.Fatalf("Decode(Encode(<struct>)) err = %v", err)
	}

	inputJson, _ := json.Marshal(input)
	outputJson, _ := json.Marshal(output)
	if !reflect.DeepEqual(inputJson, outputJson) {
		t.Fatalf("Decode(Encode(<struct>)) = \nwant: %v(%v), \ngot:  %v(%v)", reflect.TypeOf(input), input, reflect.TypeOf(output), output)
	}
}
