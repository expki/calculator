package schema

import "sync"

type Calculator struct {
	lock    sync.RWMutex
	Display string
	X       int
	Y       int
}

func (c *Calculator) SetDisplay(display string) {
	c.lock.Lock()
	c.Display = display
	c.lock.Unlock()
}

func (c *Calculator) GetDisplay() string {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.Display
}

func (c *Calculator) SetX(x int) {
	c.lock.Lock()
	c.X = x
	c.lock.Unlock()
}

func (c *Calculator) GetX() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.X
}

func (c *Calculator) SetY(y int) {
	c.lock.Lock()
	c.Y = y
	c.lock.Unlock()
}

func (c *Calculator) GetY() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.Y
}
