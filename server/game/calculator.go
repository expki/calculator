package game

import (
	"fmt"
	"regexp"
	"strings"

	"go/ast"
	"go/parser"
	"go/token"
	"strconv"

	"github.com/expki/calculator/lib/schema"
)

func (s *Session) Calculate() (display string) {
	expr := s.state.Calculator.Equation
	if expr == "" {
		return ""
	}

	// 1. Replace user-friendly symbols with operators
	expr = strings.ReplaceAll(expr, "÷", "/")
	expr = strings.ReplaceAll(expr, "×", "*")

	// 2. Parse the expression into an AST
	node, err := parser.ParseExpr(expr)
	if err != nil {
		return "Error"
	}

	// 3. Evaluate the AST
	result, err := evalAST(node)
	if err != nil {
		return "Error"
	}

	// 4. Convert to string
	return formatFloatClean(result)
}

var regexFinalDecimal = regexp.MustCompile(`\d+(\.\d+)?$`)

func (s *Session) handleUserInput(userIn schema.Input) {
	for idx, member := range s.state.Members {
		if member.Member.Id != s.Id {
			continue
		}

		s.state.Members[idx].Member.X = userIn.X
		s.state.Members[idx].Member.Y = userIn.Y

		s.handleClear(&s.state.Members[idx], userIn)
		s.handleBracket(&s.state.Members[idx], userIn)
		s.handlePercentage(&s.state.Members[idx], userIn)
		s.handleOperator(&s.state.Members[idx], userIn)
		s.handleNumberKeys(&s.state.Members[idx], userIn)
		s.handleNegate(&s.state.Members[idx], userIn)
		s.handleDecimal(&s.state.Members[idx], userIn)
		s.handleEquals(&s.state.Members[idx], userIn)

		break
	}
}

func (s *Session) handleClear(member *schema.MemberState, userIn schema.Input) {
	if member.Clear != userIn.Clear {
		if userIn.Clear {
			s.state.Calculator.Display = ""
			s.state.Calculator.Equation = ""
		}
		member.Clear = userIn.Clear
	}
}

func (s *Session) handleBracket(member *schema.MemberState, userIn schema.Input) {
	if member.Bracket != userIn.Bracket {
		if userIn.Bracket {
			s.state.Calculator.Display = ""
			open := strings.Count(s.state.Calculator.Equation, "(")
			close := strings.Count(s.state.Calculator.Equation, ")")
			if open > close {
				s.state.Calculator.Equation += ")"
			} else {
				s.state.Calculator.Equation += "("
			}
		}
		member.Bracket = userIn.Bracket
	}
}

func (s *Session) handlePercentage(member *schema.MemberState, userIn schema.Input) {
	if member.Percentage != userIn.Percentage {
		if userIn.Percentage {
			s.state.Calculator.Display = ""
			s.state.Calculator.Equation += "×0.01"
		}
		member.Percentage = userIn.Percentage
	}
}

func (s *Session) handleOperator(member *schema.MemberState, userIn schema.Input) {
	if member.Divide != userIn.Divide {
		if userIn.Divide {
			s.state.Calculator.Display = ""
			s.state.Calculator.Equation += "÷"
		}
		member.Divide = userIn.Divide
	}
	if member.Times != userIn.Times {
		if userIn.Times {
			s.state.Calculator.Display = ""
			s.state.Calculator.Equation += "×"
		}
		member.Times = userIn.Times
	}
	if member.Minus != userIn.Minus {
		if userIn.Minus {
			s.state.Calculator.Display = ""
			s.state.Calculator.Equation += "-"
		}
		member.Minus = userIn.Minus
	}
	if member.Plus != userIn.Plus {
		if userIn.Plus {
			s.state.Calculator.Display = ""
			s.state.Calculator.Equation += "+"
		}
		member.Plus = userIn.Plus
	}
}

func (s *Session) handleNumberKeys(member *schema.MemberState, userIn schema.Input) {
	buttons := []struct {
		pressed bool
		val     string
		active  *bool
	}{
		{userIn.Seven, "7", &member.Seven},
		{userIn.Eight, "8", &member.Eight},
		{userIn.Nine, "9", &member.Nine},
		{userIn.Four, "4", &member.Four},
		{userIn.Five, "5", &member.Five},
		{userIn.Six, "6", &member.Six},
		{userIn.One, "1", &member.One},
		{userIn.Two, "2", &member.Two},
		{userIn.Three, "3", &member.Three},
		{userIn.Zero, "0", &member.Zero},
	}

	for _, b := range buttons {
		if *b.active != b.pressed {
			if b.pressed {
				s.state.Calculator.Display += b.val
				s.state.Calculator.Equation += b.val
			}
			*b.active = b.pressed
		}
	}
}

func (s *Session) handleNegate(member *schema.MemberState, userIn schema.Input) {
	if member.Negate != userIn.Negate {
		if userIn.Negate {
			lastNumber := regexFinalDecimal.FindString(s.state.Calculator.Equation)
			if lastNumber != "" {
				s.state.Calculator.Equation = regexFinalDecimal.ReplaceAllString(
					s.state.Calculator.Equation,
					"-"+lastNumber,
				)
			}
		}
		member.Negate = userIn.Negate
	}
}

func (s *Session) handleDecimal(member *schema.MemberState, userIn schema.Input) {
	if member.Decimal != userIn.Decimal {
		if userIn.Decimal {
			lastNumber := regexFinalDecimal.FindString(s.state.Calculator.Equation)
			if lastNumber != "" && !strings.Contains(lastNumber, ".") {
				s.state.Calculator.Display += "."
				s.state.Calculator.Equation += "."
			}
		}
		member.Decimal = userIn.Decimal
	}
}

func (s *Session) handleEquals(member *schema.MemberState, userIn schema.Input) {
	if member.Equals != userIn.Equals {
		if userIn.Equals {
			s.state.Calculator.Display = s.Calculate()
			s.state.Calculator.Equation = ""
		}
		member.Equals = userIn.Equals
	}
}

func evalAST(node ast.Expr) (float64, error) {
	switch n := node.(type) {
	case *ast.BasicLit:
		return strconv.ParseFloat(n.Value, 64)

	case *ast.BinaryExpr:
		x, err := evalAST(n.X)
		if err != nil {
			return 0, err
		}
		y, err := evalAST(n.Y)
		if err != nil {
			return 0, err
		}
		switch n.Op {
		case token.ADD:
			return x + y, nil
		case token.SUB:
			return x - y, nil
		case token.MUL:
			return x * y, nil
		case token.QUO:
			if y == 0 {
				return 0, fmt.Errorf("division by zero")
			}
			return x / y, nil
		default:
			return 0, fmt.Errorf("unsupported operator: %s", n.Op)
		}

	case *ast.ParenExpr:
		return evalAST(n.X)

	case *ast.UnaryExpr:
		val, err := evalAST(n.X)
		if err != nil {
			return 0, err
		}
		switch n.Op {
		case token.SUB:
			return -val, nil
		case token.ADD:
			return val, nil
		default:
			return 0, fmt.Errorf("unsupported unary operator: %s", n.Op)
		}

	default:
		return 0, fmt.Errorf("unsupported expression: %T", n)
	}
}

func formatFloatClean(f float64) string {
	// If it's an integer, show it as an integer
	if f == float64(int(f)) {
		return strconv.Itoa(int(f))
	}

	// Format with precision and strip trailing zeros
	s := strconv.FormatFloat(f, 'f', 5, 64) // up to 5 decimal places
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".") // remove '.' if all decimals were removed

	return s
}
