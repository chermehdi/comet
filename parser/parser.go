package parser

import (
	"fmt"
	"github.com/chermehdi/comet/lexer"
	"strconv"
)

const (
	MINIMUM = 0
	ADD     = 1
	MUL     = 2
)

var precedences = map[lexer.TokenType]int{
	lexer.Plus:  ADD,
	lexer.Minus: ADD,
	lexer.Mul:   MUL,
	lexer.Div:   MUL,
}

func getPrecedence(token lexer.Token) int {
	val, has := precedences[token.Type]
	if !has {
		return MINIMUM
	}
	return val
}

// Functions of this type are going to be used to parse binary operations such as addition subtraction ...
// The first parameters is the already parsed left side of the operator and the function should parse

// The right side, and merge both of them and return them as a BinaryExpression
type binaryParseFunction func(Expression) Expression

// Function of this type are going to be use to parse unary operations such as ! -a +12
// The return value is Prefix Expression representing the parsed expression
type prefixParseFunction func() Expression

type Parser struct {
	lexer *lexer.Lexer

	CurrentToken lexer.Token
	NextToken    lexer.Token

	prefixFuncs map[lexer.TokenType]prefixParseFunction
	binaryFuncs map[lexer.TokenType]binaryParseFunction
}

func New(src string) *Parser {
	lexer := lexer.NewLexer(src)
	parser := &Parser{
		lexer: lexer,
	}
	parser.init()
	return parser
}

// Initialize the state of the parser
// Register the prefix and binary parsing functions and initializes first 2 tokens (current and next)
func (p *Parser) init() {
	p.advance()
	p.advance()
	p.prefixFuncs = make(map[lexer.TokenType]prefixParseFunction)
	p.binaryFuncs = make(map[lexer.TokenType]binaryParseFunction)

	p.registerPrefixFunc(p.parseNumberLiteral, lexer.Number)
	p.registerBinaryFunc(p.parseBinaryExpression, lexer.Plus, lexer.Mul, lexer.Minus, lexer.Div)
}

func (p *Parser) registerPrefixFunc(fun prefixParseFunction, tokenTypes ...lexer.TokenType) {
	for _, t := range tokenTypes {
		p.prefixFuncs[t] = fun
	}
}

func (p *Parser) registerBinaryFunc(fun binaryParseFunction, tokenTypes ...lexer.TokenType) {
	for _, t := range tokenTypes {
		p.binaryFuncs[t] = fun
	}
}

// Changes the current token to the next token.
func (p *Parser) advance() {
	p.CurrentToken = p.NextToken
	p.NextToken = p.lexer.Next()
}

// Parse the program and return a RootNode representing the root of the AST.
func (p *Parser) Parse() *RootNode {
	statements := make([]Statement, 0)
	for p.CurrentToken.Type != lexer.EOF {
		// TODO: function based language is better in this context.
		statement := p.parseStatement()
		if statement != nil {
			statements = append(statements, statement)
		}
		p.advance()
	}
	return &RootNode{
		Statements: statements,
	}
}

func (p *Parser) parseStatement() Statement {
	switch p.CurrentToken {
	// TODO: add statements
	default:
		return p.parseExpression()
	}
}

func (p *Parser) parseExpression() Expression {
	return p.parseInternal(MINIMUM)
}

func (p *Parser) parseNumberLiteral() Expression {
	val, err := strconv.ParseInt(p.CurrentToken.Literal, 10, 64)
	if err != nil {
		panic("Could not parse integer value")
	}
	return &NumberLiteralExpression{ActualValue: val}
}

func (p *Parser) parseBinaryExpression(left Expression) Expression {
	binary := &BinaryExpression{
		Left: left,
		Op:   p.CurrentToken,
	}
	precedence := getPrecedence(p.CurrentToken)
	p.advance()
	right := p.parseInternal(precedence)
	binary.Right = right
	return binary
}

func (p *Parser) parseInternal(currentPrecedence int) Expression {
	prefix, has := p.prefixFuncs[p.CurrentToken.Type]
	if !has {
		panic(fmt.Sprintf("No parsing function found for %s", p.CurrentToken))
	}
	left := prefix()
	for currentPrecedence < getPrecedence(p.NextToken) {
		binary, has := p.binaryFuncs[p.NextToken.Type]
		p.advance()
		if !has {
			return left
		}
		left = binary(left)
	}
	return left
}

func (p *Parser) advanceExpect(expected lexer.TokenType) {
	if p.CurrentToken.Type != expected {
		panic(fmt.Sprintf("Expected %s got %s", expected, p.CurrentToken.Literal))
	}
	p.advance()
}
