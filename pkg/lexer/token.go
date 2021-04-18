package lexer

type TokenType string

type Token struct {
	Type         TokenType
	Literal      string
	LineNumber   int
	ColumnNumber int
}

// Creates a new token from the given literal and type.
func NewToken(tokenType TokenType, literal string) Token {
	return Token{
		tokenType,
		literal,
		0,
		0,
	}
}

// Creates a new token from the given literal and type.
func NewTokenWithMeta(tokenType TokenType, literal string, line, column int) Token {
	return Token{
		tokenType,
		literal,
		line,
		column,
	}
}

// Token types
const (
	// Special tokens
	EOF = "EOF"

	// Operators
	Plus  = "+"
	Minus = "-"
	Mul   = "*"
	Div   = "/"
	Bang  = "!"

	// Logical operators
	GT     = ">"
	GTE    = ">="
	LT     = "<"
	Assign = "="
	LTE    = "<="
	EQ     = "=="
	NEQ    = "!="

	// Bitwise operators
	RSHIFT = ">>"
	LSHIFT = "<<"
	OR     = "|"
	AND    = "&"
	XOR    = "^"
	NOT    = "~"

	OROR   = "||"
	ANDAND = "&&"

	// Structural tokens
	OpenParent   = "("
	CloseParent  = ")"
	OpenBracket  = "["
	CloseBracket = "]"
	OpenBrace    = "{"
	CloseBrace   = "}"

	// Keywords
	Func   = "func"
	New    = "new"
	Struct = "struct"
	Return = "return"
	Var    = "var"
	True   = "true"
	False  = "false"
	If     = "if"
	Else   = "else"
	For    = "for"
	In     = "in"

	// Seperators
	Comma   = ","
	Dot     = "."
	DotDot  = ".."
	SemiCol = ";"

	// Identifier
	Identifier = "Identifier"
	Number     = "[0-9]"
	String     = "String"
)

// All keywords recognized by comet.
var Keywords = map[string]TokenType{
	"func":   Func,
	"new":    New,
	"struct": Struct,
	"return": Return,
	"var":    Var,
	"true":   True,
	"false":  False,
	"if":     If,
	"else":   Else,
	"for":    For,
	"in":     In,
}
