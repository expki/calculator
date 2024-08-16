package userinput

func (u *UserInput) setWidth(w int) {
	u.lock.Lock()
	defer u.lock.Unlock()
	u.width = w
}

func (u *UserInput) setHeight(h int) {
	u.lock.Lock()
	defer u.lock.Unlock()
	u.height = h
}

func (u *UserInput) setKey(k string, v bool) {
	u.lock.Lock()
	defer u.lock.Unlock()
	u.key[k] = v
}

func (u *UserInput) setMouseLeft(b bool) {
	u.lock.Lock()
	defer u.lock.Unlock()
	u.mouseleft = b
}

func (u *UserInput) setMouseX(x int) {
	u.lock.Lock()
	defer u.lock.Unlock()
	u.mousex = x
}

func (u *UserInput) setMouseY(y int) {
	u.lock.Lock()
	defer u.lock.Unlock()
	u.mousey = y
}
