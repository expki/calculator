package schema

import "sync"

type Member struct {
	lock sync.RWMutex
	Id   string
	X    int
	Y    int
}

func (m *Member) SetId(id string) {
	m.lock.Lock()
	m.Id = id
	m.lock.Unlock()
}

func (m *Member) SetX(x int) {
	m.lock.Lock()
	m.X = x
	m.lock.Unlock()
}

func (m *Member) SetY(y int) {
	m.lock.Lock()
	m.Y = y
	m.lock.Unlock()
}

func (m *Member) GetId() string {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.Id
}

func (m *Member) GetX() int {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.X
}

func (m *Member) GetY() int {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.Y
}
