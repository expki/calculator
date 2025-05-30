package schema

type Input struct {
	// mouse
	X int
	Y int
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
