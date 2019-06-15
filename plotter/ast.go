package plotter

import (
	"fmt"
	"math"

	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/participle/lexer/ebnf"
)

type function struct {
	Expressions expressions  `@@`
	Ranges      []identRange `"|" @@ ("," @@)*`
}

type expressions struct {
	X expression `@@`
	Y expression `"," @@`
	Z expression `"," @@`
}

type identRange struct {
	Ident string `@Ident`
	Start int    `"<" "-" @Number`
	End   int    `"." "." @Number`
}

type expression struct {
	Left  term     `@@`
	Right []opTerm `{ @@ }`
}

func (ex expression) eval(e env) float64 {
	n := ex.Left.eval(e)
	for _, r := range ex.Right {
		n = r.Operator.eval(n, r.Term.eval(e))
	}
	return n
}

type term struct {
	Left  factor     `@@`
	Right []opFactor `{ @@ }`
}

func (t term) eval(e env) float64 {
	n := t.Left.eval(e)
	for _, r := range t.Right {
		n = r.Operator.eval(n, r.Factor.eval(e))
	}
	return n
}

type factor struct {
	Base     value  `@@`
	Exponent *value `[ "^" @@ ]`
}

func (f factor) eval(e env) float64 {
	n := f.Base.eval(e)
	if f.Exponent != nil {
		return math.Pow(n, f.Exponent.eval(e))
	}
	return n
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

func (o operator) eval(l, r float64) float64 {
	switch o {
	case opMul:
		return l * r
	case opDiv:
		return l / r
	case opAdd:
		return l + r
	case opSub:
		return l - r
	}
	panic("unsupported operator")
}

type opTerm struct {
	Operator operator `@("+" | "-")`
	Term     term     `@@`
}

type opFactor struct {
	Operator operator `@("*" | "/")`
	Factor   factor   `@@`
}

type value struct {
	Number        *float64    `  @Number`
	Call          *call       `| @@`
	Ident         *string     `| @Ident`
	Subexpression *expression `| "(" @@ ")"`
}

func (v value) eval(e env) float64 {
	if v.Number != nil {
		return *v.Number
	}
	if v.Call != nil {
		return v.Call.eval(e)
	}
	if v.Ident != nil {
		value, ok := e[*v.Ident]
		if !ok {
			panic(fmt.Sprint("undefined variable: ", *v.Ident))
		}
		valueNum, ok := value.(float64)
		if !ok {
			panic(fmt.Sprint(*v.Ident, " is not a number"))
		}
		return valueNum
	}
	if v.Subexpression != nil {
		return v.Subexpression.eval(e)
	}
	panic("value has no non-nil fields")
}

type call struct {
	Name string       `@Ident`
	Args []expression `"(" [ @@ { "," @@ } ] ")"`
}

func (c call) eval(e env) float64 {
	value, ok := e[c.Name]
	if !ok {
		panic(fmt.Sprint("undefined function: ", c.Name))
	}
	valueFn, ok := value.((func(...float64) float64))
	if !ok {
		panic(fmt.Sprint(c.Name, " is not a function"))
	}
	args := []float64{}
	for _, a := range c.Args {
		args = append(args, a.eval(e))
	}
	return valueFn(args...)
}

var parser = participle.MustBuild(
	&function{},
	participle.Lexer(lexer.Must(ebnf.New(`
Ident = alpha { alpha } .
Number = ["+" | "-"] digit { digit } .
Whitespace = " " .
Punct = "!"…"/" | ":"…"@" | "["…`+"\"`\""+` | "{"…"~" .
alpha = "a"…"z" | "A"…"Z" .
digit = "0"…"9" .
`))),
	participle.Elide("Whitespace"),
)

func parse(input string) (ast function) {
	parser.ParseString(input, &ast)
	return
}
