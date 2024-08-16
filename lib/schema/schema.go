package schema

type StateLocal struct {
	Cpu float32 // cpu utilization ratio (0.0 - 1.0)
}

type StateUser struct {
	X int
	Y int
}

type StateMember struct {
	X      int
	Y      int
	Name   string
	Colour string
}

type StateCalculator struct {
	Display string
}

type StateGlobal struct {
	Members    []StateMember
	Calculator StateCalculator
}

type State struct {
	Local  StateLocal  // stays on local instance
	User   StateUser   // sends to server
	Global StateGlobal // receives from server
}
