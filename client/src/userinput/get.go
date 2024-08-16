package userinput

func (u *UserInput) GetWidth() int {
	u.lock.RLock()
	defer u.lock.RUnlock()
	return u.width
}

func (u *UserInput) GetHeight() int {
	u.lock.RLock()
	defer u.lock.RUnlock()
	return u.height
}

func (u *UserInput) GetKey(k string) (v bool) {
	u.lock.RLock()
	defer u.lock.RUnlock()
	return u.key[k]
}

func (u *UserInput) GetMouseLeft() bool {
	u.lock.RLock()
	defer u.lock.RUnlock()
	return u.mouseleft
}

func (u *UserInput) GetMouseX() int {
	u.lock.RLock()
	defer u.lock.RUnlock()
	return u.mousex
}

func (u *UserInput) GetMouseY() int {
	u.lock.RLock()
	defer u.lock.RUnlock()
	return u.mousey
}
