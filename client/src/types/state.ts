export type StateLocalExt = StateLocal & {
    RenderCpu?: number, // render cpu utilization ratio (0.0 - 1.0)
}

export type StateLocal = {
	Cpu: number, // cpu utilization ratio (0.0 - 1.0)
}

export type StateUser = {
	X: number,
	Y: number,
}

export type StateMember = {
	X:      number,
	Y:      number,
	Name:   string,
	Colour: string,
}

export type StateCalculator = {
	Display: string,
}

export type StateGlobal = {
	Members:    StateMember[],
	Calculator: StateCalculator,
}

export type State = {
	Local:  StateLocal,  // stays on local instance
	User:   StateUser,   // sends to server
	Global: StateGlobal, // receives from server
}

export type StateExt = {
	Local:  StateLocalExt,  // stays on local instance
	User:   StateUser,   // sends to server
	Global: StateGlobal, // receives from server
}
