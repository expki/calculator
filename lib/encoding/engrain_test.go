package encoding_test

import (
	"reflect"
	"testing"

	"github.com/expki/calculator/lib/encoding"
)

// Test_Engrain_Basic tests basic struct mapping
func Test_Engrain_Basic(t *testing.T) {
	type testInner struct {
		C string
	}
	type test struct {
		A int
		B testInner
	}

	input := map[string]any{
		"A": 1,
		"B": map[string]any{
			"C": "hello",
		},
	}

	var output test
	want := test{
		A: 1,
		B: testInner{
			C: "hello",
		},
	}

	err := encoding.Engrain(input, &output)
	if err != nil {
		t.Fatalf("Engrain() error = %v", err)
	}
	if output != want {
		t.Fatalf("Engrain() = %v, want %v", output, want)
	}
}

// Test_Engrain_NilData tests with nil input data
func Test_Engrain_NilData(t *testing.T) {
	type test struct {
		A int
	}

	var output test
	err := encoding.Engrain(nil, &output)
	if err != nil {
		t.Fatalf("Engrain(nil, &struct) error = %v", err)
	}

	if output.A != 0 {
		t.Fatalf("Expected zero value struct, got %v", output)
	}
}

// Test_Engrain_InvalidDestination tests error cases for destination
func Test_Engrain_InvalidDestination(t *testing.T) {
	tests := []struct {
		name string
		dst  any
		want string
	}{
		{
			name: "non-pointer",
			dst:  struct{}{},
			want: "engrain dst must be a non-nil pointer",
		},
		{
			name: "nil pointer",
			dst:  (*struct{})(nil),
			want: "engrain dst must not be nil",
		},
		{
			name: "pointer to non-struct",
			dst:  new(int),
			want: "engrain dst must be a pointer to a struct",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := encoding.Engrain(map[string]any{"key": "value"}, tt.dst)
			if err == nil {
				t.Fatalf("Expected error, got nil")
			}
			if err.Error() != tt.want {
				t.Fatalf("Error = %q, want %q", err.Error(), tt.want)
			}
		})
	}
}

// Test_Engrain_Slices tests slice handling
func Test_Engrain_Slices(t *testing.T) {
	type Item struct {
		Name  string
		Value int
	}

	type Container struct {
		SimpleSlice  []int
		StructSlice  []Item
		PointerSlice []*Item
		StringSlice  []string
	}

	input := map[string]any{
		"SimpleSlice": []any{1, 2, 3},
		"StructSlice": []any{
			map[string]any{"Name": "item1", "Value": 10},
			map[string]any{"Name": "item2", "Value": 20},
		},
		"PointerSlice": []any{
			map[string]any{"Name": "ptr1", "Value": 100},
			map[string]any{"Name": "ptr2", "Value": 200},
		},
		"StringSlice": []any{"a", "b", "c"},
	}

	var output Container
	err := encoding.Engrain(input, &output)
	if err != nil {
		t.Fatalf("Engrain() error = %v", err)
	}

	// Check SimpleSlice
	if !reflect.DeepEqual(output.SimpleSlice, []int{1, 2, 3}) {
		t.Errorf("SimpleSlice = %v, want [1, 2, 3]", output.SimpleSlice)
	}

	// Check StructSlice
	if len(output.StructSlice) != 2 {
		t.Errorf("StructSlice length = %d, want 2", len(output.StructSlice))
	}
	if output.StructSlice[0].Name != "item1" || output.StructSlice[0].Value != 10 {
		t.Errorf("StructSlice[0] = %+v, want {Name:item1 Value:10}", output.StructSlice[0])
	}

	// Check PointerSlice
	if len(output.PointerSlice) != 2 {
		t.Errorf("PointerSlice length = %d, want 2", len(output.PointerSlice))
	}
	if output.PointerSlice[0].Name != "ptr1" || output.PointerSlice[0].Value != 100 {
		t.Errorf("PointerSlice[0] = %+v, want {Name:ptr1 Value:100}", output.PointerSlice[0])
	}

	// Check StringSlice
	if !reflect.DeepEqual(output.StringSlice, []string{"a", "b", "c"}) {
		t.Errorf("StringSlice = %v, want [a, b, c]", output.StringSlice)
	}
}

// Test_Engrain_Maps tests map handling
func Test_Engrain_Maps(t *testing.T) {
	type Item struct {
		Name  string
		Value int
	}

	type Container struct {
		SimpleMap    map[string]int
		StructMap    map[string]Item
		PointerMap   map[string]*Item
		InterfaceMap map[string]interface{}
	}

	input := map[string]any{
		"SimpleMap": map[string]any{
			"one": 1,
			"two": 2,
		},
		"StructMap": map[string]any{
			"first":  map[string]any{"Name": "item1", "Value": 10},
			"second": map[string]any{"Name": "item2", "Value": 20},
		},
		"PointerMap": map[string]any{
			"ptr1": map[string]any{"Name": "ptrItem1", "Value": 100},
			"ptr2": map[string]any{"Name": "ptrItem2", "Value": 200},
		},
		"InterfaceMap": map[string]any{
			"str":  "hello",
			"num":  42,
			"bool": true,
		},
	}

	var output Container
	err := encoding.Engrain(input, &output)
	if err != nil {
		t.Fatalf("Engrain() error = %v", err)
	}

	// Check SimpleMap
	if output.SimpleMap["one"] != 1 || output.SimpleMap["two"] != 2 {
		t.Errorf("SimpleMap = %v, want map[one:1 two:2]", output.SimpleMap)
	}

	// Check StructMap
	if output.StructMap["first"].Name != "item1" || output.StructMap["first"].Value != 10 {
		t.Errorf("StructMap[first] = %+v, want {Name:item1 Value:10}", output.StructMap["first"])
	}

	// Check PointerMap
	if output.PointerMap["ptr1"].Name != "ptrItem1" || output.PointerMap["ptr1"].Value != 100 {
		t.Errorf("PointerMap[ptr1] = %+v, want {Name:ptrItem1 Value:100}", output.PointerMap["ptr1"])
	}

	// Check InterfaceMap
	if output.InterfaceMap["str"] != "hello" {
		t.Errorf("InterfaceMap[str] = %v, want hello", output.InterfaceMap["str"])
	}
}

// Test_Engrain_Pointers tests pointer field handling
func Test_Engrain_Pointers(t *testing.T) {
	type Inner struct {
		Value string
	}

	type Container struct {
		IntPtr    *int
		StringPtr *string
		StructPtr *Inner
		NilPtr    *int
	}

	input := map[string]any{
		"IntPtr":    42,
		"StringPtr": "hello",
		"StructPtr": map[string]any{
			"Value": "inner value",
		},
		"NilPtr": nil,
	}

	var output Container
	err := encoding.Engrain(input, &output)
	if err != nil {
		t.Fatalf("Engrain() error = %v", err)
	}

	// Check IntPtr
	if output.IntPtr == nil || *output.IntPtr != 42 {
		t.Errorf("IntPtr = %v, want 42", output.IntPtr)
	}

	// Check StringPtr
	if output.StringPtr == nil || *output.StringPtr != "hello" {
		t.Errorf("StringPtr = %v, want hello", output.StringPtr)
	}

	// Check StructPtr
	if output.StructPtr == nil || output.StructPtr.Value != "inner value" {
		t.Errorf("StructPtr = %+v, want {Value:inner value}", output.StructPtr)
	}

	// Check NilPtr
	if output.NilPtr != nil {
		t.Errorf("NilPtr = %v, want nil", output.NilPtr)
	}
}

// Test_Engrain_TypeConversion tests type conversions
func Test_Engrain_TypeConversion(t *testing.T) {
	type Container struct {
		Int8    int8
		Int16   int16
		Int32   int32
		Int64   int64
		Uint8   uint8
		Uint16  uint16
		Uint32  uint32
		Uint64  uint64
		Float32 float32
		Float64 float64
		Bool    bool
		String  string
	}

	input := map[string]any{
		"Int8":    int(8),
		"Int16":   int(16),
		"Int32":   int(32),
		"Int64":   int(64),
		"Uint8":   int(8),
		"Uint16":  int(16),
		"Uint32":  int(32),
		"Uint64":  int(64),
		"Float32": float64(3.14),
		"Float64": float64(2.718),
		"Bool":    true,
		"String":  "test",
	}

	var output Container
	err := encoding.Engrain(input, &output)
	if err != nil {
		t.Fatalf("Engrain() error = %v", err)
	}

	if output.Int8 != 8 {
		t.Errorf("Int8 = %d, want 8", output.Int8)
	}
	if output.Float32 != 3.14 {
		t.Errorf("Float32 = %f, want 3.14", output.Float32)
	}
	if output.Bool != true {
		t.Errorf("Bool = %v, want true", output.Bool)
	}
	if output.String != "test" {
		t.Errorf("String = %s, want test", output.String)
	}
}

// Test_Engrain_UnknownFields tests that unknown fields are ignored
func Test_Engrain_UnknownFields(t *testing.T) {
	type Container struct {
		Known string
	}

	input := map[string]any{
		"Known":   "value",
		"Unknown": "should be ignored",
		"Another": 123,
	}

	var output Container
	err := encoding.Engrain(input, &output)
	if err != nil {
		t.Fatalf("Engrain() error = %v", err)
	}

	if output.Known != "value" {
		t.Errorf("Known = %s, want value", output.Known)
	}
}

// Test_Engrain_UnexportedFields tests that unexported fields are ignored
func Test_Engrain_UnexportedFields(t *testing.T) {
	type Container struct {
		Exported   string
		unexported string
	}

	input := map[string]any{
		"Exported":   "visible",
		"unexported": "should be ignored",
	}

	var output Container
	err := encoding.Engrain(input, &output)
	if err != nil {
		t.Fatalf("Engrain() error = %v", err)
	}

	if output.Exported != "visible" {
		t.Errorf("Exported = %s, want visible", output.Exported)
	}
	if output.unexported != "" {
		t.Errorf("unexported = %s, want empty", output.unexported)
	}
}

// Test_Engrain_Interface tests interface{} field handling
func Test_Engrain_Interface(t *testing.T) {
	type Container struct {
		Any  interface{}
		Data interface{}
	}

	input := map[string]any{
		"Any":  "string value",
		"Data": 42,
	}

	var output Container
	err := encoding.Engrain(input, &output)
	if err != nil {
		t.Fatalf("Engrain() error = %v", err)
	}

	if output.Any != "string value" {
		t.Errorf("Any = %v, want string value", output.Any)
	}
	if output.Data != 42 {
		t.Errorf("Data = %v, want 42", output.Data)
	}
}

// Test_Engrain_ComplexNested tests deeply nested structures
func Test_Engrain_ComplexNested(t *testing.T) {
	type Level3 struct {
		Deep string
	}

	type Level2 struct {
		Mid   int
		Next  Level3
		Items []Level3
	}

	type Level1 struct {
		Top     string
		Nested  Level2
		Pointer *Level2
		Map     map[string]Level2
	}

	input := map[string]any{
		"Top": "top level",
		"Nested": map[string]any{
			"Mid": 100,
			"Next": map[string]any{
				"Deep": "deep value",
			},
			"Items": []any{
				map[string]any{"Deep": "item1"},
				map[string]any{"Deep": "item2"},
			},
		},
		"Pointer": map[string]any{
			"Mid": 200,
			"Next": map[string]any{
				"Deep": "pointer deep",
			},
			"Items": []any{
				map[string]any{"Deep": "ptr item"},
			},
		},
		"Map": map[string]any{
			"key1": map[string]any{
				"Mid": 300,
				"Next": map[string]any{
					"Deep": "map deep",
				},
				"Items": []any{},
			},
		},
	}

	var output Level1
	err := encoding.Engrain(input, &output)
	if err != nil {
		t.Fatalf("Engrain() error = %v", err)
	}

	// Verify nested structure
	if output.Top != "top level" {
		t.Errorf("Top = %s, want top level", output.Top)
	}
	if output.Nested.Mid != 100 {
		t.Errorf("Nested.Mid = %d, want 100", output.Nested.Mid)
	}
	if output.Nested.Next.Deep != "deep value" {
		t.Errorf("Nested.Next.Deep = %s, want deep value", output.Nested.Next.Deep)
	}
	if len(output.Nested.Items) != 2 {
		t.Errorf("Nested.Items length = %d, want 2", len(output.Nested.Items))
	}

	// Verify pointer
	if output.Pointer == nil {
		t.Errorf("Pointer is nil, want non-nil")
	} else if output.Pointer.Mid != 200 {
		t.Errorf("Pointer.Mid = %d, want 200", output.Pointer.Mid)
	}

	// Verify map
	if val, ok := output.Map["key1"]; !ok {
		t.Errorf("Map[key1] not found")
	} else if val.Mid != 300 {
		t.Errorf("Map[key1].Mid = %d, want 300", val.Mid)
	}
}

// Test_Engrain_EmptyCollections tests empty slices and maps
func Test_Engrain_EmptyCollections(t *testing.T) {
	type Container struct {
		EmptySlice []int
		EmptyMap   map[string]string
		NilSlice   []int
		NilMap     map[string]string
	}

	input := map[string]any{
		"EmptySlice": []any{},
		"EmptyMap":   map[string]any{},
		"NilSlice":   nil,
		"NilMap":     nil,
	}

	var output Container
	err := encoding.Engrain(input, &output)
	if err != nil {
		t.Fatalf("Engrain() error = %v", err)
	}

	if output.EmptySlice == nil {
		t.Errorf("EmptySlice is nil, want empty slice")
	}
	if len(output.EmptySlice) != 0 {
		t.Errorf("EmptySlice length = %d, want 0", len(output.EmptySlice))
	}

	if output.EmptyMap == nil {
		t.Errorf("EmptyMap is nil, want empty map")
	}
	if len(output.EmptyMap) != 0 {
		t.Errorf("EmptyMap length = %d, want 0", len(output.EmptyMap))
	}

	if output.NilSlice != nil {
		t.Errorf("NilSlice = %v, want nil", output.NilSlice)
	}
	if output.NilMap != nil {
		t.Errorf("NilMap = %v, want nil", output.NilMap)
	}
}

// Test_Engrain_ComprehensivePointerSupport tests comprehensive pointer/non-pointer support
func Test_Engrain_ComprehensivePointerSupport(t *testing.T) {
	type NestedStruct struct {
		Value string
		Count int
	}

	type ComplexStruct struct {
		// Basic non-pointer types
		BoolVal    bool
		IntVal     int
		Int8Val    int8
		Int16Val   int16
		Int32Val   int32
		Int64Val   int64
		UintVal    uint
		Uint8Val   uint8
		Uint16Val  uint16
		Uint32Val  uint32
		Uint64Val  uint64
		Float32Val float32
		Float64Val float64
		StringVal  string

		// Basic pointer types
		BoolPtr    *bool
		IntPtr     *int
		Int8Ptr    *int8
		Int16Ptr   *int16
		Int32Ptr   *int32
		Int64Ptr   *int64
		UintPtr    *uint
		Uint8Ptr   *uint8
		Uint16Ptr  *uint16
		Uint32Ptr  *uint32
		Uint64Ptr  *uint64
		Float32Ptr *float32
		Float64Ptr *float64
		StringPtr  *string

		// Struct types
		StructVal NestedStruct
		StructPtr *NestedStruct

		// Slice types
		SliceInt          []int
		SliceIntPtr       []*int
		SliceStruct       []NestedStruct
		SliceStructPtr    []*NestedStruct
		SliceOfSlice      [][]int
		SliceOfSlicePtr   [][]*int
		SliceOfPtrSlice   []*[]int
		SliceOfPtrSlicePtr []*[]*int

		// Map types
		MapStringInt       map[string]int
		MapStringIntPtr    map[string]*int
		MapStringStruct    map[string]NestedStruct
		MapStringStructPtr map[string]*NestedStruct
		MapOfMap           map[string]map[string]int
		MapOfMapPtr        map[string]map[string]*int
		MapOfPtrMap        map[string]*map[string]int

		// Mixed nested types
		MapOfSlice    map[string][]int
		MapOfSlicePtr map[string][]*int
		SliceOfMap    []map[string]int
		SliceOfMapPtr []*map[string]int

		// Interface types
		InterfaceVal interface{}
		SliceInterface []interface{}
		MapInterface   map[string]interface{}
	}

	// Create test data
	boolTrue := true
	intVal := 42
	int8Val := int8(8)
	float32Val := float32(3.14)
	stringVal := "test"

	input := map[string]any{
		// Basic non-pointer values
		"BoolVal":    true,
		"IntVal":     42,
		"Int8Val":    int8(8),
		"Int16Val":   int16(16),
		"Int32Val":   int32(32),
		"Int64Val":   int64(64),
		"UintVal":    uint(42),
		"Uint8Val":   uint8(8),
		"Uint16Val":  uint16(16),
		"Uint32Val":  uint32(32),
		"Uint64Val":  uint64(64),
		"Float32Val": float32(3.14),
		"Float64Val": 2.718,
		"StringVal":  "test",

		// Basic pointer values (source as non-pointer, should create pointer)
		"BoolPtr":    true,
		"IntPtr":     42,
		"Int8Ptr":    int8(8),
		"Int16Ptr":   int16(16),
		"Int32Ptr":   int32(32),
		"Int64Ptr":   int64(64),
		"UintPtr":    uint(42),
		"Uint8Ptr":   uint8(8),
		"Uint16Ptr":  uint16(16),
		"Uint32Ptr":  uint32(32),
		"Uint64Ptr":  uint64(64),
		"Float32Ptr": float32(3.14),
		"Float64Ptr": 2.718,
		"StringPtr":  "test",

		// Struct values
		"StructVal": map[string]any{
			"Value": "nested",
			"Count": 100,
		},
		"StructPtr": map[string]any{
			"Value": "nested ptr",
			"Count": 200,
		},

		// Slice values
		"SliceInt":    []any{1, 2, 3},
		"SliceIntPtr": []any{10, 20, 30},
		"SliceStruct": []any{
			map[string]any{"Value": "s1", "Count": 1},
			map[string]any{"Value": "s2", "Count": 2},
		},
		"SliceStructPtr": []any{
			map[string]any{"Value": "sp1", "Count": 10},
			map[string]any{"Value": "sp2", "Count": 20},
		},
		"SliceOfSlice": []any{
			[]any{1, 2},
			[]any{3, 4},
		},
		"SliceOfSlicePtr": []any{
			[]any{11, 12},
			[]any{13, 14},
		},
		"SliceOfPtrSlice": []any{
			[]any{21, 22},
			[]any{23, 24},
		},
		"SliceOfPtrSlicePtr": []any{
			[]any{31, 32},
			[]any{33, 34},
		},

		// Map values
		"MapStringInt": map[string]any{
			"one": 1,
			"two": 2,
		},
		"MapStringIntPtr": map[string]any{
			"ten":    10,
			"twenty": 20,
		},
		"MapStringStruct": map[string]any{
			"first": map[string]any{"Value": "m1", "Count": 100},
		},
		"MapStringStructPtr": map[string]any{
			"ptr": map[string]any{"Value": "mp1", "Count": 200},
		},
		"MapOfMap": map[string]any{
			"outer": map[string]any{
				"inner": 42,
			},
		},
		"MapOfMapPtr": map[string]any{
			"outer": map[string]any{
				"inner": 84,
			},
		},
		"MapOfPtrMap": map[string]any{
			"outer": map[string]any{
				"inner": 168,
			},
		},

		// Mixed nested types
		"MapOfSlice": map[string]any{
			"list": []any{1, 2, 3},
		},
		"MapOfSlicePtr": map[string]any{
			"list": []any{10, 20, 30},
		},
		"SliceOfMap": []any{
			map[string]any{"a": 1},
			map[string]any{"b": 2},
		},
		"SliceOfMapPtr": []any{
			map[string]any{"x": 10},
			map[string]any{"y": 20},
		},

		// Interface types
		"InterfaceVal":   "interface string",
		"SliceInterface": []any{1, "two", true},
		"MapInterface": map[string]any{
			"num":  42,
			"str":  "hello",
			"bool": false,
		},
	}

	var output ComplexStruct
	err := encoding.Engrain(input, &output)
	if err != nil {
		t.Fatalf("Engrain() error = %v", err)
	}

	// Verify basic non-pointer types
	if output.BoolVal != true {
		t.Errorf("BoolVal = %v, want true", output.BoolVal)
	}
	if output.IntVal != 42 {
		t.Errorf("IntVal = %v, want 42", output.IntVal)
	}
	if output.StringVal != "test" {
		t.Errorf("StringVal = %v, want test", output.StringVal)
	}

	// Verify basic pointer types
	if output.BoolPtr == nil || *output.BoolPtr != boolTrue {
		t.Errorf("BoolPtr = %v, want &true", output.BoolPtr)
	}
	if output.IntPtr == nil || *output.IntPtr != intVal {
		t.Errorf("IntPtr = %v, want &42", output.IntPtr)
	}
	if output.StringPtr == nil || *output.StringPtr != stringVal {
		t.Errorf("StringPtr = %v, want &test", output.StringPtr)
	}

	// Verify numeric pointer types
	if output.Int8Ptr == nil || *output.Int8Ptr != int8Val {
		t.Errorf("Int8Ptr = %v, want &8", output.Int8Ptr)
	}
	if output.Float32Ptr == nil || *output.Float32Ptr != float32Val {
		t.Errorf("Float32Ptr = %v, want &3.14", output.Float32Ptr)
	}

	// Verify struct types
	if output.StructVal.Value != "nested" || output.StructVal.Count != 100 {
		t.Errorf("StructVal = %+v, want {Value:nested Count:100}", output.StructVal)
	}
	if output.StructPtr == nil || output.StructPtr.Value != "nested ptr" || output.StructPtr.Count != 200 {
		t.Errorf("StructPtr = %+v, want &{Value:nested ptr Count:200}", output.StructPtr)
	}

	// Verify slice types
	if len(output.SliceInt) != 3 || output.SliceInt[0] != 1 {
		t.Errorf("SliceInt = %v, want [1 2 3]", output.SliceInt)
	}
	if len(output.SliceIntPtr) != 3 || output.SliceIntPtr[0] == nil || *output.SliceIntPtr[0] != 10 {
		t.Errorf("SliceIntPtr[0] = %v, want &10", output.SliceIntPtr[0])
	}
	if len(output.SliceStruct) != 2 || output.SliceStruct[0].Value != "s1" {
		t.Errorf("SliceStruct[0] = %+v, want {Value:s1 Count:1}", output.SliceStruct[0])
	}
	if len(output.SliceStructPtr) != 2 || output.SliceStructPtr[0] == nil || output.SliceStructPtr[0].Value != "sp1" {
		t.Errorf("SliceStructPtr[0] = %+v, want &{Value:sp1 Count:10}", output.SliceStructPtr[0])
	}

	// Verify nested slice types
	if len(output.SliceOfSlice) != 2 || len(output.SliceOfSlice[0]) != 2 || output.SliceOfSlice[0][0] != 1 {
		t.Errorf("SliceOfSlice[0][0] = %v, want 1", output.SliceOfSlice[0][0])
	}
	if len(output.SliceOfSlicePtr) != 2 || len(output.SliceOfSlicePtr[0]) != 2 || 
		output.SliceOfSlicePtr[0][0] == nil || *output.SliceOfSlicePtr[0][0] != 11 {
		t.Errorf("SliceOfSlicePtr[0][0] = %v, want &11", output.SliceOfSlicePtr[0][0])
	}

	// Verify map types
	if output.MapStringInt["one"] != 1 {
		t.Errorf("MapStringInt[one] = %v, want 1", output.MapStringInt["one"])
	}
	if output.MapStringIntPtr["ten"] == nil || *output.MapStringIntPtr["ten"] != 10 {
		t.Errorf("MapStringIntPtr[ten] = %v, want &10", output.MapStringIntPtr["ten"])
	}
	if output.MapStringStruct["first"].Value != "m1" {
		t.Errorf("MapStringStruct[first].Value = %v, want m1", output.MapStringStruct["first"].Value)
	}
	if output.MapStringStructPtr["ptr"] == nil || output.MapStringStructPtr["ptr"].Value != "mp1" {
		t.Errorf("MapStringStructPtr[ptr].Value = %v, want mp1", output.MapStringStructPtr["ptr"].Value)
	}

	// Verify nested map types
	if output.MapOfMap["outer"]["inner"] != 42 {
		t.Errorf("MapOfMap[outer][inner] = %v, want 42", output.MapOfMap["outer"]["inner"])
	}
	if output.MapOfMapPtr["outer"]["inner"] == nil || *output.MapOfMapPtr["outer"]["inner"] != 84 {
		t.Errorf("MapOfMapPtr[outer][inner] = %v, want &84", output.MapOfMapPtr["outer"]["inner"])
	}

	// Verify mixed types
	if len(output.MapOfSlice["list"]) != 3 || output.MapOfSlice["list"][0] != 1 {
		t.Errorf("MapOfSlice[list][0] = %v, want 1", output.MapOfSlice["list"][0])
	}
	if len(output.SliceOfMap) != 2 || output.SliceOfMap[0]["a"] != 1 {
		t.Errorf("SliceOfMap[0][a] = %v, want 1", output.SliceOfMap[0]["a"])
	}

	// Verify interface types
	if output.InterfaceVal != "interface string" {
		t.Errorf("InterfaceVal = %v, want interface string", output.InterfaceVal)
	}
	if len(output.SliceInterface) != 3 {
		t.Errorf("SliceInterface length = %v, want 3", len(output.SliceInterface))
	}
	if output.MapInterface["num"] != 42 {
		t.Errorf("MapInterface[num] = %v, want 42", output.MapInterface["num"])
	}
}

// Test_Engrain_StringToBoolConversion tests string to bool conversion
func Test_Engrain_StringToBoolConversion(t *testing.T) {
	type Container struct {
		BoolTrue1  bool
		BoolTrue2  bool
		BoolFalse1 bool
		BoolFalse2 bool
		BoolPtr    *bool
	}

	input := map[string]any{
		"BoolTrue1":  "true",
		"BoolTrue2":  "1",
		"BoolFalse1": "false",
		"BoolFalse2": "0",
		"BoolPtr":    "true",
	}

	var output Container
	err := encoding.Engrain(input, &output)
	if err != nil {
		t.Fatalf("Engrain() error = %v", err)
	}

	if output.BoolTrue1 != true {
		t.Errorf("BoolTrue1 = %v, want true", output.BoolTrue1)
	}
	if output.BoolTrue2 != true {
		t.Errorf("BoolTrue2 = %v, want true", output.BoolTrue2)
	}
	if output.BoolFalse1 != false {
		t.Errorf("BoolFalse1 = %v, want false", output.BoolFalse1)
	}
	if output.BoolFalse2 != false {
		t.Errorf("BoolFalse2 = %v, want false", output.BoolFalse2)
	}
	if output.BoolPtr == nil || *output.BoolPtr != true {
		t.Errorf("BoolPtr = %v, want &true", output.BoolPtr)
	}
}

// Test_Engrain_NilHandling tests comprehensive nil handling
func Test_Engrain_NilHandling(t *testing.T) {
	type Inner struct {
		Value string
	}

	type Container struct {
		StringPtr  *string
		IntPtr     *int
		StructPtr  *Inner
		SlicePtr   *[]int
		MapPtr     *map[string]int
		Interface  interface{}
		SliceOfPtr []*int
		MapOfPtr   map[string]*int
	}

	input := map[string]any{
		"StringPtr":  nil,
		"IntPtr":     nil,
		"StructPtr":  nil,
		"SlicePtr":   nil,
		"MapPtr":     nil,
		"Interface":  nil,
		"SliceOfPtr": []any{nil, nil},
		"MapOfPtr": map[string]any{
			"a": nil,
			"b": nil,
		},
	}

	var output Container
	err := encoding.Engrain(input, &output)
	if err != nil {
		t.Fatalf("Engrain() error = %v", err)
	}

	if output.StringPtr != nil {
		t.Errorf("StringPtr = %v, want nil", output.StringPtr)
	}
	if output.IntPtr != nil {
		t.Errorf("IntPtr = %v, want nil", output.IntPtr)
	}
	if output.StructPtr != nil {
		t.Errorf("StructPtr = %v, want nil", output.StructPtr)
	}
	if output.SlicePtr != nil {
		t.Errorf("SlicePtr = %v, want nil", output.SlicePtr)
	}
	if output.MapPtr != nil {
		t.Errorf("MapPtr = %v, want nil", output.MapPtr)
	}
	if output.Interface != nil {
		t.Errorf("Interface = %v, want nil", output.Interface)
	}
	if len(output.SliceOfPtr) != 2 || output.SliceOfPtr[0] != nil || output.SliceOfPtr[1] != nil {
		t.Errorf("SliceOfPtr = %v, want [nil nil]", output.SliceOfPtr)
	}
	if output.MapOfPtr["a"] != nil || output.MapOfPtr["b"] != nil {
		t.Errorf("MapOfPtr values should all be nil")
	}
}
