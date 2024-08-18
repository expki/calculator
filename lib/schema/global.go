package schema

import "sync"

type Global struct {
	lock       sync.RWMutex
	Calculator Calculator
	Members    []Member
	members    map[string]*Member
}

func (g *Global) AddMember(id string, x, y int) {
	g.lock.Lock()
	if g.members == nil {
		g.members = make(map[string]*Member)
	}
	g.Members = append(g.Members, Member{Id: id, X: x, Y: y})
	g.members[id] = &g.Members[len(g.Members)-1]
	g.lock.Unlock()
}

func (g *Global) RemoveMember(id string) {
	g.lock.Lock()
	for i := range g.Members {
		if g.Members[i].Id == id {
			g.Members = append(g.Members[:i], g.Members[i+1:]...)
			delete(g.members, id)
			break
		}
	}
	g.lock.Unlock()
}

func (g *Global) GetMember(id string) *Member {
	g.lock.RLock()
	member, ok := g.members[id]
	g.lock.RUnlock()
	if ok {
		return member
	}
	g.lock.Lock()
	defer g.lock.Unlock()
	for i := range g.Members {
		if g.Members[i].Id == id {
			g.members[id] = &g.Members[i]
			return &g.Members[i]
		}
	}
	return nil
}
