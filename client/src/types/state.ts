export type StateExt = State & {
    CpuRender?: number, // render cpu utilization ratio (0.0 - 1.0)
}

export type State = {
	Global: Global,
	CpuLogic: number,
}

export type Global = {
	Calculator: Calculator,
	Members: null | {[id: string]: Member},
}

export type Member = {
    Id: string,
    X: number,
    Y: number,
}

export type Calculator = {
	Display: string,
    X: number,
    Y: number,
}
