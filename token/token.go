package token

// A set of token types.
const (
	Illegal = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers and literals
	Ident  = "IDENT" // add, foo, bar, x, y, etc
	Num    = "NUM"   // numeric literal, only integers for now
	String = "STRING"

	// Operators
	Assign    = "="
	Plus      = "+"
	Minus     = "-"
	Bang      = "!"
	Asterisk  = "*"
	Slash     = "/"
	LT        = "<"
	GT        = ">"
	Equals    = "=="
	NotEquals = "!="

	// Delimiters
	Comma     = ","
	Semicolon = ";"
	LParen    = "("
	RParen    = ")"
	LBrace    = "{"
	RBrace    = "}"
	LBracket  = "["
	RBracket  = "]"

	// Keywords
	Function = "FUNCTION"
	Let      = "LET"
	True     = "TRUE"
	False    = "FALSE"
	If       = "IF"
	Else     = "ELSE"
	Return   = "RETURN"
)

var keywords = map[string]Type{
	"fn":     Function,
	"let":    Let,
	"true":   True,
	"false":  False,
	"if":     If,
	"else":   Else,
	"return": Return,
}

// Type identifies a token type.
type Type string

// Token represents a token.
type Token struct {
	Type    Type
	Literal string
}

// LookupIdent returns an appropriate token type for identifier: if it's a
// reserved keyword, it will return that keyword's type, otherwise Ident.
func LookupIdent(ident string) Type {
	if tokType, ok := keywords[ident]; ok {
		return tokType
	}
	return Ident
}
