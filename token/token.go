package token

// A set of token types.
const (
	Illegal = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers and literals
	Ident = "IDENT" // add, foo, bar, x, y, etc
	Num   = "NUM"   // numeric literal, only integers for now

	// Operators
	Assign = "="
	Plus   = "+"

	// Delimiters
	Comma     = ","
	Semicolon = ";"
	LParen    = "("
	RParen    = ")"
	LBrace    = "{"
	RBrace    = "}"

	// Keywords
	Function = "FUNCTION"
	Let      = "LET"
)

var keywords = map[string]Type{
	"fn":  Function,
	"let": Let,
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
