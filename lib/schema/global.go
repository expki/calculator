package schema

type Global struct {
	Calculator Calculator
	Members    map[string]*Member
}

func (g *Global) AddMember(id string, x, y int) {
	g.Members[id] = &Member{Id: id, X: x, Y: y}
}

func (g *Global) RemoveMember(id string) {
	delete(g.Members, id)
}

func (g *Global) GetMember(id string) *Member {
	return g.Members[id]
}
