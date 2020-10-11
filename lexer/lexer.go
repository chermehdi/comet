package lexer

import (
	"fmt"
	"unicode"
)

type Lexer struct {
	src       string
	pos       int
	current   byte
	inputSize int
}

// Creates an initializes a new lexer from the given input source.
func NewLexer(src string) *Lexer {
	return &Lexer{
		src:       src,
		pos:       0,
		current:   src[0],
		inputSize: len(src),
	}
}

func (l *Lexer) Next() Token {
	var result Token
	l.ignoreWhiteSpace()
	switch l.current {
	case '+':
		result = NewToken(Plus, "+")
	case '-':
		result = NewToken(Minus, "-")
	case '*':
		result = NewToken(Mul, "*")
	case '/':
		result = NewToken(Div, "/")
	case '>':
		if l.peek() == '=' {
			l.advance()
			result = NewToken(GTE, ">=")
		} else {
			result = NewToken(GT, ">")
		}
	case '<':
		if l.peek() == '=' {
			l.advance()
			result = NewToken(LTE, "<=")
		} else {
			result = NewToken(LT, "<")
		}
	case '=':
		if l.peek() == '=' {
			l.advance()
			result = NewToken(EQ, "==")
		} else {
			result = NewToken(Assign, "=")
		}
	case '!':
		if l.peek() == '=' {
			l.advance()
			result = NewToken(NEQ, "!=")
		} else {
			result = NewToken(Bang, "!")
		}
	case '(':
		result = NewToken(OpenParent, "(")
	case ')':
		result = NewToken(CloseParent, ")")
	case '[':
		result = NewToken(OpenBracket, "[")
	case ']':
		result = NewToken(CloseBracket, "]")
	case '{':
		result = NewToken(OpenBrace, "{")
	case '}':
		result = NewToken(CloseBrace, "}")
	case '.':
		if l.peek() == '.' {
			l.advance()
			result = NewToken(DotDot, "..")
		} else {
			result = NewToken(Dot, ".")
		}
	case ';':
		result = NewToken(SemiCol, ";")
	case ',':
		result = NewToken(Comma, ",")
	case 0:
		result = NewToken(EOF, "EOF")
	case '"':
		result = l.readString()
	default:
		if unicode.IsDigit(rune(l.current)) {
			result = l.readNumber()
		} else if unicode.IsLetter(rune(l.current)) || l.current == '_' {
			result = l.readIdentifier()
		}
	}
	l.advance()
	return result
}

// Escape white space.
// Whitespace is anything of '\n' '\r' ' ' '\t'
func (l *Lexer) ignoreWhiteSpace() {
	for isWhiteSpace(l.current) {
		l.advance()
	}
}

func (l *Lexer) advance() {
	l.pos += 1
	if l.pos < l.inputSize {
		l.current = l.src[l.pos]
	} else {
		// Indicates EOF
		l.current = 0
	}
}

func (l *Lexer) peek() byte {
	if l.pos+1 < l.inputSize {
		return l.src[l.pos+1]
	}
	return 0
}

func (l *Lexer) readIdentifier() Token {
	start := l.pos
	for {
		if !identifierCharacter(l.peek()) {
			break
		}
		l.advance()
	}
	literal := l.src[start : l.pos+1]
	literalType, has := Keywords[literal]
	if !has {
		return NewToken(Identifier, literal)
	}
	return NewToken(literalType, literal)
}

func identifierCharacter(c byte) bool {
	return c == '_' || unicode.IsDigit(rune(c)) || unicode.IsLetter(rune(c))
}

// TODO add support for other kind of formats
// examples: +1 -2 1.12 1e12 0x16 0777
func (l *Lexer) readNumber() Token {
	start := l.pos
	for {
		if !unicode.IsDigit(rune(l.peek())) {
			break
		}
		l.advance()
	}
	return NewToken(Number, l.src[start:l.pos+1])
}

func (l *Lexer) readString() Token {
	start := l.pos + 1
	// "some string"
	for {
		l.advance()
		if l.current == '"' {
			break
		}
		if l.current == '\n' || l.current == '\r' || l.current == 0 {
			// TODO: panic is not proper error handling, fix it.
			panic(fmt.Sprint("Reached the end of line or end of input without closing the string quote"))
		}
	}
	return NewToken(String, l.src[start:l.pos])
}

func isWhiteSpace(c byte) bool {
	return unicode.IsSpace(rune(c))
}
