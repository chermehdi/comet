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
	line      int
	column    int
}

// Creates an initializes a new lexer from the given input source.
func NewLexer(src string) *Lexer {
	return &Lexer{
		src:       src,
		pos:       0,
		current:   src[0],
		inputSize: len(src),
		line:      1,
		column:    1,
	}
}

func (l *Lexer) Next() Token {
	var result Token
	l.ignoreWhiteSpace()
	switch l.current {
	case '+':
		result = NewTokenWithMeta(Plus, "+", l.line, l.column)
	case '-':
		result = NewTokenWithMeta(Minus, "-", l.line, l.column)
	case '*':
		result = NewTokenWithMeta(Mul, "*", l.line, l.column)
	case '/':
		result = NewTokenWithMeta(Div, "/", l.line, l.column)
	case '^':
		result = NewTokenWithMeta(XOR, "^", l.line, l.column)
	case '~':
		result = NewTokenWithMeta(NOT, "~", l.line, l.column)
	case '>':
		if l.peek() == '=' {
			l.advance()
			result = NewTokenWithMeta(GTE, ">=", l.line, l.column)
		} else if l.peek() == '>' {
			l.advance()
			result = NewTokenWithMeta(RSHIFT, ">>", l.line, l.column)
		} else {
			result = NewTokenWithMeta(GT, ">", l.line, l.column)
		}
	case '<':
		if l.peek() == '=' {
			l.advance()
			result = NewTokenWithMeta(LTE, "<=", l.line, l.column)
		} else if l.peek() == '<' {
			l.advance()
			result = NewTokenWithMeta(LSHIFT, "<<", l.line, l.column)
		} else {
			result = NewTokenWithMeta(LT, "<", l.line, l.column)
		}
	case '=':
		if l.peek() == '=' {
			l.advance()
			result = NewTokenWithMeta(EQ, "==", l.line, l.column)
		} else {
			result = NewTokenWithMeta(Assign, "=", l.line, l.column)
		}
	case '!':
		if l.peek() == '=' {
			l.advance()
			result = NewTokenWithMeta(NEQ, "!=", l.line, l.column)
		} else {
			result = NewTokenWithMeta(Bang, "!", l.line, l.column)
		}
	case '&':
		if l.peek() == '&' {
			l.advance()
			result = NewTokenWithMeta(ANDAND, "&&", l.line, l.column)
		} else {
			result = NewTokenWithMeta(AND, "&", l.line, l.column)
		}
	case '|':
		if l.peek() == '|' {
			l.advance()
			result = NewTokenWithMeta(OROR, "||", l.line, l.column)
		} else {
			result = NewTokenWithMeta(OR, "|", l.line, l.column)
		}
	case '(':
		result = NewTokenWithMeta(OpenParent, "(", l.line, l.column)
	case ')':
		result = NewTokenWithMeta(CloseParent, ")", l.line, l.column)
	case '[':
		result = NewTokenWithMeta(OpenBracket, "[", l.line, l.column)
	case ']':
		result = NewTokenWithMeta(CloseBracket, "]", l.line, l.column)
	case '{':
		result = NewTokenWithMeta(OpenBrace, "{", l.line, l.column)
	case '}':
		result = NewTokenWithMeta(CloseBrace, "}", l.line, l.column)
	case '.':
		if l.peek() == '.' {
			l.advance()
			result = NewTokenWithMeta(DotDot, "..", l.line, l.column)
		} else {
			result = NewTokenWithMeta(Dot, ".", l.line, l.column)
		}
	case ';':
		result = NewTokenWithMeta(SemiCol, ";", l.line, l.column)
	case ',':
		result = NewTokenWithMeta(Comma, ",", l.line, l.column)
	case 0:
		result = NewTokenWithMeta(EOF, "EOF", l.line, l.column)
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
		if l.current == '\n' {
			l.line++
			l.column = 1
		}
		l.advance()
	}
}

func (l *Lexer) advance() {
	l.pos += 1
	l.column += 1
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
	return NewTokenWithMeta(Number, l.src[start:l.pos+1], l.line, l.column)
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
	return NewTokenWithMeta(String, l.src[start:l.pos], l.line, l.column)
}

func isWhiteSpace(c byte) bool {
	return unicode.IsSpace(rune(c))
}
