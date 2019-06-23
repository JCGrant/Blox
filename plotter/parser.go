package plotter

import (
	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/participle/lexer/ebnf"
)

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
