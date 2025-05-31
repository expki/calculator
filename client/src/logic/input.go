package logic

import (
	"calculator/src/types"

	"github.com/expki/calculator/lib/schema"
)

type LocalInput struct {
	mouseleft bool
	mousex    int
	mousey    int
	width     int
	height    int
	x         int
	y         int
	// row 1
	clear      bool
	bracket    bool
	percentage bool
	divide     bool
	// row 2
	seven bool
	eight bool
	nine  bool
	times bool
	// row 3
	four  bool
	five  bool
	six   bool
	minus bool
	// row 4
	one   bool
	two   bool
	three bool
	plus  bool
	// row 5
	negate  bool
	zero    bool
	decimal bool
	equals  bool
}

func (local *LocalInput) Translate() (value schema.Input) {
	value.X = local.x
	value.Y = local.y
	if !local.mouseleft {
		return value
	}
	// row 1
	value.Clear = local.clear
	value.Bracket = local.bracket
	value.Percentage = local.percentage
	value.Divide = local.divide
	// row 2
	value.Seven = local.seven
	value.Eight = local.eight
	value.Nine = local.nine
	value.Times = local.times
	// row 3
	value.Four = local.four
	value.Five = local.five
	value.Six = local.six
	value.Minus = local.minus
	// row 4
	value.One = local.one
	value.Two = local.two
	value.Three = local.three
	value.Plus = local.plus
	// row 5
	value.Negate = local.negate
	value.Zero = local.zero
	value.Decimal = local.decimal
	value.Equals = local.equals
	return value
}

func (local *LocalInput) Meged(state types.LocalState) types.LocalState {
	// row 1
	state.Calculator.Clear = state.Calculator.Clear || local.clear
	state.Calculator.Bracket = state.Calculator.Bracket || local.bracket
	state.Calculator.Percentage = state.Calculator.Percentage || local.percentage
	state.Calculator.Divide = state.Calculator.Divide || local.divide
	// row 2
	state.Calculator.Seven = state.Calculator.Seven || local.seven
	state.Calculator.Eight = state.Calculator.Eight || local.eight
	state.Calculator.Nine = state.Calculator.Nine || local.nine
	state.Calculator.Times = state.Calculator.Times || local.times
	// row 3
	state.Calculator.Four = state.Calculator.Four || local.four
	state.Calculator.Five = state.Calculator.Five || local.five
	state.Calculator.Six = state.Calculator.Six || local.six
	state.Calculator.Minus = state.Calculator.Minus || local.minus
	// row 4
	state.Calculator.One = state.Calculator.One || local.one
	state.Calculator.Two = state.Calculator.Two || local.two
	state.Calculator.Three = state.Calculator.Three || local.three
	state.Calculator.Plus = state.Calculator.Plus || local.plus
	// row 5
	state.Calculator.Negate = state.Calculator.Negate || local.negate
	state.Calculator.Zero = state.Calculator.Zero || local.zero
	state.Calculator.Decimal = state.Calculator.Decimal || local.decimal
	state.Calculator.Equals = state.Calculator.Equals || local.equals
	return state
}
