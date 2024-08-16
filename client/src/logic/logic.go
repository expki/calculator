package logic

import (
	"syscall/js"
	"time"

	"github.com/expki/calculator/lib/encoding"
	"github.com/expki/calculator/lib/schema"
)

const target_tick_rate = time.Millisecond * 1000 / 60

func LogicLoop(sharedArray js.Value) {
	pref := [60]float32{0.0}
	state := schema.State{}
	for n := 0; true; n++ {
		// Measure start time
		start := time.Now()

		// Update state
		state.Local.Cpu = 0.0
		for _, value := range pref {
			state.Local.Cpu += value
		}
		state.Local.Cpu /= 5

		// Encode state
		data := encoding.Encode(state)
		js.CopyBytesToJS(sharedArray, data)

		// Wait to hit target tick rate
		end := time.Since(start)
		if end < target_tick_rate {
			time.Sleep(target_tick_rate - end)
		}

		// Record tick rate
		pref[n%5] = float32(end) / float32(target_tick_rate)
	}
}
