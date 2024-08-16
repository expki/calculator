package userinput

import (
	"sync"
	"syscall/js"
)

type UserInput struct {
	lock      sync.RWMutex
	handle    js.Func
	width     int
	height    int
	key       map[string]bool
	mouseleft bool
	mousex    int
	mousey    int
}

func New() *UserInput {
	window := js.Global()
	u := UserInput{
		key: make(map[string]bool),
	}
	u.handle = js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) < 1 {
			return nil
		}
		arg := args[0]
		if value := arg.Get("width"); !value.IsUndefined() {
			u.setWidth(value.Int())
		}
		if value := arg.Get("height"); !value.IsUndefined() {
			u.setHeight(value.Int())
		}
		if value := arg.Get("keyup"); !value.IsUndefined() {
			u.setKey(value.String(), false)
		}
		if value := arg.Get("keydown"); !value.IsUndefined() {
			u.setKey(value.String(), true)
		}
		if value := arg.Get("mouseleft"); !value.IsUndefined() {
			u.setMouseLeft(value.Bool())
		}
		if value := arg.Get("mousex"); !value.IsUndefined() {
			u.setMouseX(value.Int())
		}
		if value := arg.Get("mousey"); !value.IsUndefined() {
			u.setMouseY(value.Int())
		}
		return nil
	})
	window.Set("handleInput", u.handle)
	return &u
}

func (u *UserInput) Close() {
	u.lock.Lock()
	defer u.lock.Unlock()
	u.handle.Release()
}
