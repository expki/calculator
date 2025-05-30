package schema_test

import (
	"testing"

	"github.com/expki/calculator/lib/encoding"
	"github.com/expki/calculator/lib/schema"
)

// Test_Global_Encode ...
func Test_Global_Encode(t *testing.T) {
	input := schema.State{
		Calculator: schema.Calculator{},
		Members:    []schema.Member{{}},
	}
	data := encoding.Encode(&input)
	generic, _ := encoding.Decode(data)
	var output schema.State
	err := encoding.Engrain(generic.(map[string]any), &output)
	if err != nil {
		t.Fatalf("Engrain(Decode(Encode(<*state>))) = %v", err)
	}
}
