package schema_test

import (
	"testing"

	"github.com/expki/calculator/lib/encoding"
	"github.com/expki/calculator/lib/schema"
)

// Test_Global_Encode ...
func Test_Global_Encode(t *testing.T) {
	input := schema.Global{
		Calculator: schema.Calculator{},
		Members:    nil,
	}
	data := encoding.Encode(&input)
	generic, _ := encoding.Decode(data)
	var output schema.Global
	err := encoding.Engrain(generic.(map[string]any), &output)
	if err != nil {
		t.Fatalf("Engrain(Decode(Encode(<*global>))) = %v", err)
	}
}
