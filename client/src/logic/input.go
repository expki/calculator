package logic

import (
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
