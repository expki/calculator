package compression_test

import (
	"reflect"
	"testing"

	"github.com/expki/calculator/lib/compression"
)

// Test_Compress ...
func Test_Compress(t *testing.T) {
	var input []byte = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	var output []byte
	out := compression.Compress(input)
	original, err := compression.Decompress(out)
	if err != nil {
		t.Fatalf("Decompress(<null>) err = %v", err)
	}
	if !reflect.DeepEqual(input, original) {
		t.Fatalf("Decompress(Compress(<null>)) = \nwant: %v, \ngot:  %v", input, output)
	}
}
