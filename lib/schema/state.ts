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
    Clear: boolean,
    Bracket: boolean,
    Percentage: boolean,
    Divide: boolean,
    // row 2
    Seven: boolean,
    Eight:  boolean,
    Nine: boolean,
    Times: boolean,
    // row 3
    Four: boolean,
    Five: boolean,
    Six: boolean,
    Minus: boolean,
    // row 4
    One: boolean,
    Two: boolean,
    Three: boolean,
    Plus: boolean,
    // row 5
    Negate: boolean,
    Zero: boolean,
    Decimal: boolean,
    Equals: boolean,
};
