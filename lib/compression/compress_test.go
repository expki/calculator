package compression_test

import (
	"bytes"
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
	if bytes.Equal(input, out) {
		t.Fatalf("Decompress(<data>) = data, no compression occured")
	}
	if err != nil {
		t.Fatalf("Decompress(<data>) err = %v", err)
	}
	if !reflect.DeepEqual(input, original) {
		t.Fatalf("Decompress(Compress(<data>)) = \nwant: %v, \ngot:  %v", input, output)
	}
}
