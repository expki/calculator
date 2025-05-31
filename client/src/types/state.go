package types

import "github.com/expki/calculator/lib/schema"

type LocalState struct {
	schema.State
	Id      int
	CpuLoad float32
}
