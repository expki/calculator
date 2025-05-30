export type State = Partial<{
	Calculator: Calculator,
	Members: null | Member[],
}>;

export type Member = {
    Id: number,
    X: number,
    Y: number,
};

export type Calculator = {
    // row 0
    Equation: string,
	Display: string,
    // row 1
    clear: boolean,
    bracket: boolean,
    percentage: boolean,
    divide: boolean,
    // row 2
    seven: boolean,
    eight:  boolean,
    nine: boolean,
    times: boolean,
    // row 3
    four: boolean,
    five: boolean,
    six: boolean,
    minus: boolean,
    // row 4
    one: boolean,
    two: boolean,
    three: boolean,
    plus: boolean,
    // row 5
    negate: boolean,
    zero: boolean,
    decimal: boolean,
    equals: boolean,
};
