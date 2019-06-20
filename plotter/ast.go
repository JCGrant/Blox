package plotter

type function struct {
	Expressions *expressions  `@@`
	Ranges      []*identRange `"|" @@ ("," @@)*`
}

type expressions struct {
	X *expression `@@`
	Y *expression `"," @@`
	Z *expression `"," @@`
}

type identRange struct {
	Ident *string `@Ident`
	Start *int    `"<" "-" @Number`
	End   *int    `"." "." @Number`
}

type expression struct {
	Left  *term     `@@`
	Right []*opTerm `{ @@ }`
}

type term struct {
	Left  *factor     `@@`
	Right []*opFactor `{ @@ }`
}

type factor struct {
	Base     *value `@@`
	Exponent *value `[ "^" @@ ]`
}

type operator int

const (
	opMul operator = iota
	opDiv
	opAdd
	opSub
)

var operatorMap = map[string]operator{"+": opAdd, "-": opSub, "*": opMul, "/": opDiv}

func (o *operator) Capture(s []string) error {
	*o = operatorMap[s[0]]
	return nil
}

type opTerm struct {
	Operator operator `@("+" | "-")`
	Term     *term    `@@`
}

type opFactor struct {
	Operator operator `@("*" | "/")`
	Factor   *factor  `@@`
}

type value struct {
	Number        *float64    `  @Number`
	Call          *call       `| @@`
	Ident         *string     `| @Ident`
	Subexpression *expression `| "(" @@ ")"`
}

type call struct {
	Name *string       `@Ident`
	Args []*expression `"(" [ @@ { "," @@ } ] ")"`
}
