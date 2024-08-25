package schema

type Calculator struct {
	Display string
	X       int
	Y       int
}

func (c *Calculator) SetDisplay(display string) {
	c.Display = display
}

func (c *Calculator) GetDisplay() string {
	return c.Display
}

func (c *Calculator) SetX(x int) {
	c.X = x
}

func (c *Calculator) GetX() int {
	return c.X
}

func (c *Calculator) SetY(y int) {
	c.Y = y
}

func (c *Calculator) GetY() int {
	return c.Y
}
