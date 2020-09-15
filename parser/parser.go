package parser

import (
	"fmt"
	"github.com/chermehdi/comet/lexer"
	"strconv"
)

const (
	MINIMUM = 0
	LOG     = 1
	ADD     = 1
	MUL     = 2
	PARENT  = 3
)

var precedences = map[lexer.TokenType]int{
	lexer.Plus:  ADD,
	lexer.Minus: ADD,
	lexer.Mul:   MUL,
	lexer.Div:   MUL,
	lexer.LT:    LOG,
	lexer.LTE:   LOG,
	lexer.GT:    LOG,
	lexer.GTE:   LOG,
	lexer.EQ:    LOG,
	lexer.NEQ:   LOG,
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
	p.registerPrefixFunc(p.parsePrefixExpression, lexer.Minus, lexer.Bang)
	p.registerPrefixFunc(p.parseIdentifier, lexer.Identifier)
	p.registerPrefixFunc(p.parseBoolean, lexer.True, lexer.False)
	p.registerPrefixFunc(p.parseParenthesisedExpression, lexer.OpenParent)
	p.registerPrefixFunc(p.parseStringLiteral, lexer.String)

	p.registerBinaryFunc(p.parseBinaryExpression, lexer.Plus, lexer.Mul, lexer.Minus, lexer.Div,
		lexer.GT, lexer.GTE, lexer.LT, lexer.LTE, lexer.EQ, lexer.NEQ)
}

// Utility method to enable prefix function registration for given token types.
func (p *Parser) registerPrefixFunc(fun prefixParseFunction, tokenTypes ...lexer.TokenType) {
	for _, t := range tokenTypes {
		p.prefixFuncs[t] = fun
	}
}

// Utility method to enable binary function registration for given token types.
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

// Try to parse a statement, it's possible just by knowing the current token type because
// the Grammar of the language allows us to. Otherwise fallback to try and parse an expression.
func (p *Parser) parseStatement() Statement {
	switch p.CurrentToken.Type {
	case lexer.Var:
		return p.parseDeclaration()
	case lexer.Return:
		return p.parseReturnStatement()
	case lexer.OpenBrace:
		return p.parseBlockStatement()
	case lexer.If:
		return p.parseIfStatement()
	case lexer.Func:
		return p.parseFunctionStatement()
	default:
		return p.parseExpression()
	}
}

// A declaration operation is anything of this form: var name = expression.
func (p *Parser) parseDeclaration() Statement {
	declarationStatement := &DeclarationStatement{
		varToken: p.CurrentToken,
	}
	p.advanceExpect(lexer.Var)
	declarationStatement.Identifier = p.CurrentToken
	p.advanceExpect(lexer.Identifier)
	p.advanceExpect(lexer.Assign)
	declarationStatement.Expression = p.parseExpression()
	return declarationStatement
}

// A return statement is anything of the form: return expression
func (p *Parser) parseReturnStatement() Statement {
	returnStatement := &ReturnStatement{
		returnToken: p.CurrentToken,
	}
	p.advanceExpect(lexer.Return)
	returnStatement.Expression = p.parseExpression()
	return returnStatement
}

// This will initiate try parsing an expression with the Minimum precedence.
func (p *Parser) parseExpression() Expression {
	return p.parseInternal(MINIMUM)
}

func (p *Parser) parsePrefixExpression() Expression {
	expression := &PrefixExpression{
		Op: p.CurrentToken,
	}
	p.advance()
	expression.Right = p.parseExpression()
	return expression
}

// A Number Literal is an expression that represents a number.
func (p *Parser) parseNumberLiteral() Expression {
	val, err := strconv.ParseInt(p.CurrentToken.Literal, 10, 64)
	if err != nil {
		panic("Could not parse integer value")
	}
	return &NumberLiteralExpression{ActualValue: val}
}

// an identifier is an expression that represents the name of a variable.
func (p *Parser) parseIdentifier() Expression {
	if p.NextToken.Type == lexer.OpenParent {
		// This is a function call
		callExpression := &CallExpression{
			Name: p.CurrentToken.Literal,
		}
		p.advance()
		callExpression.Arguments = p.parseCallArguments()
		return callExpression
	} else {
		// This is an identifier
		return &IdentifierExpression{Name: p.CurrentToken.Literal}
	}
}

func (p *Parser) parseCallArguments() []Expression {
	args := []Expression{}
	if p.NextToken.Type == lexer.CloseParent {
		p.advance()
		return args
	}
	p.advance()
	// parse first argument
	args = append(args, p.parseExpression())
	for p.NextToken.Type == lexer.Comma {
		p.advance() // Skip last token of current expression
		p.advance() // Skip the comma
		args = append(args, p.parseExpression())
	}
	p.advance()
	if p.CurrentToken.Type != lexer.CloseParent {
		panic(fmt.Sprintf("Expected %s, got %s", lexer.CloseParent, p.CurrentToken.Literal))
	}
	return args
}

// any expression of the form ( expression )
func (p *Parser) parseParenthesisedExpression() Expression {
	// (expression)
	p.advanceExpect(lexer.OpenParent)
	expression := p.parseExpression()
	parenthesised := &ParenthesisedExpression{
		Expression: expression,
	}
	p.expectNext(lexer.CloseParent)
	return parenthesised
}

// A binary expression is an expression of the form: expression operator expression
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

// Tries to parse as long as the currentPrecedence is smaller than the precedence of the next operator.
// This is an implementation of the idea of a Pratt Parser.
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

func (p *Parser) parseBlockStatement() *BlockStatement {
	blockStatement := &BlockStatement{}
	statements := make([]Statement, 0)
	p.advanceExpect(lexer.OpenBrace)
	for p.CurrentToken.Type != lexer.CloseBrace && p.CurrentToken.Type != lexer.EOF {
		curStatement := p.parseStatement()
		if curStatement == nil {
			// TODO: probably an error, fix when error handling is added.
			panic("current statement is nil")
		}
		statements = append(statements, curStatement)
		p.advance()
	}
	blockStatement.Statements = statements
	return blockStatement
}

func (p *Parser) parseBoolean() Expression {
	return &BooleanLiteral{
		ActualValue: p.CurrentToken.Type == lexer.True,
		Token:       p.CurrentToken,
	}
}

func (p *Parser) parseIfStatement() Statement {
	ifStatement := newIfStatement()

	p.advanceExpect(lexer.If)
	ifStatement.Test = p.parseExpression()
	p.expectNext(lexer.OpenBrace)

	ifStatement.Then = *p.parseBlockStatement()
	p.advanceExpect(lexer.CloseBrace)

	if p.CurrentToken.Type == lexer.Else {
		p.advanceExpect(lexer.Else)
		ifStatement.Else = *p.parseBlockStatement()
	}
	return ifStatement
}

func (p *Parser) parseFunctionStatement() Statement {
	funcStatement := newFunctionStatement()
	p.advanceExpect(lexer.Func)

	funcStatement.Name = p.CurrentToken.Literal
	p.advanceExpect(lexer.Identifier)

	p.advanceExpect(lexer.OpenParent)
	// if there are parameters
	if p.CurrentToken.Type != lexer.CloseParent {
		for {
			if p.CurrentToken.Type == lexer.EOF || p.CurrentToken.Type == lexer.CloseParent {
				break
			}
			parameterName := p.parseIdentifier()
			parameterExpression, _ := parameterName.(*IdentifierExpression)
			funcStatement.Parameters = append(funcStatement.Parameters, parameterExpression)
			p.advance()
			if p.CurrentToken.Type == lexer.Comma {
				p.advance()
			}
		}
	}
	p.advanceExpect(lexer.CloseParent)

	funcStatement.Block = p.parseBlockStatement()
	return funcStatement
}

func (p *Parser) advanceExpect(expected lexer.TokenType) {
	if p.CurrentToken.Type != expected {
		panic(fmt.Sprintf("Expected %s got %s", expected, p.CurrentToken.Literal))
	}
	p.advance()
}

func (p *Parser) expectNext(expected lexer.TokenType) {
	if p.NextToken.Type != expected {
		panic(fmt.Sprintf("Expected %s got %s", expected, p.CurrentToken.Literal))
	}
	p.advance()
}

func (p *Parser) parseStringLiteral() Expression {
	return &StringLiteral{Value: p.CurrentToken.Literal}
}
