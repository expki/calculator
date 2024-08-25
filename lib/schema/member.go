package schema

type Member struct {
	Id string
	X  int
	Y  int
}

func (m *Member) SetId(id string) {
	m.Id = id
}

func (m *Member) SetX(x int) {
	m.X = x
}

func (m *Member) SetY(y int) {
	m.Y = y
}

func (m *Member) GetId() string {
	return m.Id
}

func (m *Member) GetX() int {
	return m.X
}

func (m *Member) GetY() int {
	return m.Y
}
