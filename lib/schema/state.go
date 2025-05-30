package schema

type State struct {
	Calculator Calculator
	Members    []Member
}

type StateState struct {
	Calculator Calculator
	Members    []MemberState
}

type Member struct {
	Id int
	X  int
	Y  int
}

type MemberState struct {
	Member
	// row 1
	Clear      bool
	Bracket    bool
	Percentage bool
	Divide     bool
	// row 2
	Seven bool
	Eight bool
	Nine  bool
	Times bool
	// row 3
	Four  bool
	Five  bool
	Six   bool
	Minus bool
	// row 4
	One   bool
	Two   bool
	Three bool
	Plus  bool
	// row 5
	Negate  bool
	Zero    bool
	Decimal bool
	Equals  bool
}

type Calculator struct {
	// row 0
	Equation string
	Display  string
	// row 1
	Clear      bool
	Bracket    bool
	Percentage bool
	Divide     bool
	// row 2
	Seven bool
	Eight bool
	Nine  bool
	Times bool
	// row 3
	Four  bool
	Five  bool
	Six   bool
	Minus bool
	// row 4
	One   bool
	Two   bool
	Three bool
	Plus  bool
	// row 5
	Negate  bool
	Zero    bool
	Decimal bool
	Equals  bool
}

func (s *StateState) State() State {
	members := make([]Member, len(s.Members))
	for idx, member := range s.Members {
		members[idx] = member.Member
	}
	return State{
		Calculator: s.Calculator,
		Members:    members,
	}
}
