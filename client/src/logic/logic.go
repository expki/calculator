package logic

import (
	"bytes"
	"calculator/src/types"
	"sync"
	"syscall/js"

	"github.com/expki/calculator/lib/encoding"
	"github.com/expki/calculator/lib/schema"
)

const maxPanding = 100

type Logic struct {
	lock    sync.RWMutex
	state   *schema.Global
	pending chan struct{}
}

func New() *Logic {
	return &Logic{
		state:   &schema.Global{},
		pending: make(chan struct{}, maxPanding),
	}
}

func (l *Logic) GetState() (*schema.Global, func()) {
	l.lock.RLock()
	return l.state, func() {
		l.lock.RUnlock()
	}
}

func (l *Logic) LockState() (*schema.Global, func()) {
	l.lock.Lock()
	return l.state, func() {
		l.pending <- struct{}{}
		l.lock.Unlock()
	}
}

func (l *Logic) LogicLoop(sharedArray js.Value) {
	var lastData []byte
	for n := 0; true; n++ {
		// Wait for pending
		<-l.pending

		// Clear additional pending
		func() {
			for i := 0; i < maxPanding; i++ {
				select {
				case <-l.pending:
				default:
					return
				}
			}
		}()

		// Get state
		global, unlock := l.GetState()
		state := types.State{
			Global: global,
		}

		// Encode state
		data := encoding.Encode(state)
		if !bytes.Equal(lastData, data) {
			js.CopyBytesToJS(sharedArray, data)
		}
		unlock()
	}
}
