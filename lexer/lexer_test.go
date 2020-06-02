package lexer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLexer_Next(t *testing.T) {
	tests := []struct {
		Input          string
		ExpectedTokens []Token
	}{
		{`1 + 2`, []Token{
			NewToken(Number, "1"),
			NewToken(Plus, "+"),
			NewToken(Number, "2"),
		}},
		{`+ / -  *+`, []Token{
			NewToken(Plus, "+"),
			NewToken(Div, "/"),
			NewToken(Minus, "-"),
			NewToken(Mul, "*"),
			NewToken(Plus, "+"),
		}},
		{`< > = !`, []Token{
			NewToken(LT, "<"),
			NewToken(GT, ">"),
			NewToken(Assign, "="),
			NewToken(Bang, "!"),
		}},
		{`<= >= == !=`, []Token{
			NewToken(LTE, "<="),
			NewToken(GTE, ">="),
			NewToken(EQ, "=="),
			NewToken(NEQ, "!="),
		}},
		{`; , .`, []Token{
			NewToken(SemiCol, ";"),
			NewToken(Comma, ","),
			NewToken(Dot, "."),
		}},
		{`(){ } [ ] `, []Token{
			NewToken(OpenParent, "("),
			NewToken(CloseParent, ")"),
			NewToken(OpenBrace, "{"),
			NewToken(CloseBrace, "}"),
			NewToken(OpenBracket, "["),
			NewToken(CloseBracket, "]"),
		}},
		{`a _ _a _a1b`, []Token{
			NewToken(Identifier, "a"),
			NewToken(Identifier, "_"),
			NewToken(Identifier, "_a"),
			NewToken(Identifier, "_a1b"),
		}},
		{`12 2`, []Token{
			NewToken(Number, "12"),
			NewToken(Number, "2"),
		}},
		{`"some kind of text for strings " a1`, []Token{
			NewToken(String, "some kind of text for strings "),
			NewToken(Identifier, "a1"),
		}},
		{`func new return if else a for var true false`, []Token{
			NewToken(Func, "func"),
			NewToken(New, "new"),
			NewToken(Return, "return"),
			NewToken(If, "if"),
			NewToken(Else, "else"),
			NewToken(Identifier, "a"),
			NewToken(For, "for"),
			NewToken(Var, "var"),
			NewToken(True, "true"),
			NewToken(False, "false"),
		}},
		{`func main(a, b) {
	var a = a[0]
	if a > 0 {
		return 1
	} else {
		var f = 1
		for var i = 1; i < b; i = i + 1 {
			f = f * i
 		}
		return f
	}
}
`, []Token{
			NewToken(Func, "func"),
			NewToken(Identifier, "main"),
			NewToken(OpenParent, "("),
			NewToken(Identifier, "a"),
			NewToken(Comma, ","),
			NewToken(Identifier, "b"),
			NewToken(CloseParent, ")"),
			NewToken(OpenBrace, "{"),
			NewToken(Var, "var"),
			NewToken(Identifier, "a"),
			NewToken(Assign, "="),
			NewToken(Identifier, "a"),
			NewToken(OpenBracket, "["),
			NewToken(Number, "0"),
			NewToken(CloseBracket, "]"),
			NewToken(If, "if"),
			NewToken(Identifier, "a"),
			NewToken(GT, ">"),
			NewToken(Number, "0"),
			NewToken(OpenBrace, "{"),
			NewToken(Return, "return"),
			NewToken(Number, "1"),
			NewToken(CloseBrace, "}"),
			NewToken(Else, "else"),
			NewToken(OpenBrace, "{"),
			NewToken(Var, "var"),
			NewToken(Identifier, "f"),
			NewToken(Assign, "="),
			NewToken(Number, "1"),
			NewToken(For, "for"),
			NewToken(Var, "var"),
			NewToken(Identifier, "i"),
			NewToken(Assign, "="),
			NewToken(Number, "1"),
			NewToken(SemiCol, ";"),
			NewToken(Identifier, "i"),
			NewToken(LT, "<"),
			NewToken(Identifier, "b"),
			NewToken(SemiCol, ";"),
			NewToken(Identifier, "i"),
			NewToken(Assign, "="),
			NewToken(Identifier, "i"),
			NewToken(Plus, "+"),
			NewToken(Number, "1"),
			NewToken(OpenBrace, "{"),
			NewToken(Identifier, "f"),
			NewToken(Assign, "="),
			NewToken(Identifier, "f"),
			NewToken(Mul, "*"),
			NewToken(Identifier, "i"),
			NewToken(CloseBrace, "}"),
			NewToken(Return, "return"),
			NewToken(Identifier, "f"),

			NewToken(CloseBrace, "}"),
			NewToken(CloseBrace, "}"),
		}},
	}

	for _, test := range tests {
		gotTokens := consumeLexer(NewLexer(test.Input))

		assert.Equal(t, len(test.ExpectedTokens), len(gotTokens))

		for i, token := range test.ExpectedTokens {
			assert.Equal(t, token, gotTokens[i])
		}
	}
}

func consumeLexer(l *Lexer) []Token {
	tokens := make([]Token, 0)
	for currentToken := l.Next(); currentToken.Type != EOF; currentToken = l.Next() {
		tokens = append(tokens, currentToken)
	}
	return tokens
}
