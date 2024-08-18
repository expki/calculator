package encoding_test

import (
	"testing"

	"github.com/expki/calculator/lib/encoding"
)

// Test_Engrain ...
func Test_Engrain(t *testing.T) {
	type testInner struct {
		C string
	}
	type test struct {
		A int
		B testInner
	}
	var input map[string]any = map[string]any{
		"A": 1,
		"B": map[string]any{
			"C": "hello",
		},
	}
	var output test
	var want test = test{
		A: 1,
		B: testInner{
			C: "hello",
		},
	}
	err := encoding.Engrain(input, &output)
	if err != nil {
		t.Fatalf("Engrain(<map>, <*struct>) = %v", err)
	}
	if output != want {
		t.Fatalf("Engrain(<map>, <*struct>) = %v, want %v", output, want)
	}
}
