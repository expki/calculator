package logic

import (
	"calculator/src/types"
	"sync"
	"syscall/js"
	"time"

	"github.com/expki/calculator/lib/encoding"
	"github.com/expki/calculator/lib/schema"
)

const target_tick_rate = time.Millisecond * 1000 / 60

type Logic struct {
	lock  sync.RWMutex
	state *schema.Global
}

func New() *Logic {
	return &Logic{
		state: &schema.Global{},
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
		l.lock.Unlock()
	}
}

func (l *Logic) LogicLoop(sharedArray js.Value) {
	pref := [60]float32{0.0}
	for n := 0; true; n++ {
		// Measure start time
		start := time.Now()

		// Get state
		global, unlock := l.LockState()
		state := types.State{
			Global: global,
		}
		for _, value := range pref {
			state.CpuLogic += value
		}
		state.CpuLogic /= 5

		// Encode state
		data := encoding.Encode(state)
		js.CopyBytesToJS(sharedArray, data)
		unlock()

		// Wait to hit target tick rate
		end := time.Since(start)
		if end < target_tick_rate {
			time.Sleep(target_tick_rate - end)
		}

		// Record tick rate
		pref[n%5] = float32(end) / float32(target_tick_rate)
	}
}
