package types

import "github.com/expki/calculator/lib/schema"

type LocalState struct {
	schema.State
	Id      int
	X       int
	Y       int
	CpuLoad float32
}
