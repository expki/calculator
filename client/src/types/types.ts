export type Payload<T> = {
    kind: number, 
    payload: T,
};

export type PayloadWasm = {
    wasm: ArrayBuffer,
    pipe: SharedArrayBuffer,
};

export type PayloadInput = Partial<{
    mouseleft: boolean,
    mousex: number,
    mousey: number,
    width: number,
    height: number,
    x: number,
    y: number,
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
}>;

export type ButtonMap = {
    // row 1
    clear: ButtonState,
    bracket: ButtonState,
    percentage: ButtonState,
    divide: ButtonState,
    // row 2
    seven: ButtonState,
    eight: ButtonState,
    nine: ButtonState,
    times: ButtonState,
    // row 3
    four: ButtonState,
    five: ButtonState,
    six: ButtonState,
    minus: ButtonState,
    // row 4
    one: ButtonState,
    two: ButtonState,
    three: ButtonState,
    plus: ButtonState,
    // row 5
    negate: ButtonState,
    zero: ButtonState,
    decimal: ButtonState,
    equals: ButtonState,
};
export type ButtonProperties = {
    text: string,
    col: number,
    row: number,
};
export type ButtonLayout = ButtonProperties & {
    x: number,
    y: number,
    width: number,
    height: number,
};
export type ButtonState = ButtonLayout & {
    pressed: boolean,
    colourPrimary: string,
    colourSecondary: string,
    textColor: string,
    borderColor: string,
};
export type ButtonStyle = {
    pressed: {
        colourPrimary: string,
        colourSecondary: string,
    },
    depressed: {
        colourPrimary: string,
        colourSecondary: string,
    },
    textColor: string,
    borderColor: string,
};
