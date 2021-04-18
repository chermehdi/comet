package parser

import (
	"fmt"
	"github.com/chermehdi/comet/pkg/lexer"
	"strconv"
	"strings"
)

// Lower binds stronger
const (
	MINIMUM = iota
	LOG
	ADD
	MUL
	DOT
	PARENT
	Index
)

var precedences = map[lexer.TokenType]int{
	lexer.Plus:        ADD,
	lexer.Minus:       ADD,
	lexer.Mul:         MUL,
	lexer.Div:         MUL,
	lexer.LT:          LOG,
	lexer.LTE:         LOG,
	lexer.GT:          LOG,
	lexer.GTE:         LOG,
	lexer.EQ:          LOG,
	lexer.NEQ:         LOG,
	lexer.Dot:         DOT,
	lexer.DotDot:      PARENT,
	lexer.OpenBracket: Index,
}

func getPrecedence(token lexer.Token) int {
	val, has := precedences[token.Type]
	if !has {
		return MINIMUM
	}
	return val
}

type ParseError struct {
	Message string
	Token   lexer.Token
}

func (p *ParseError) Error() string {
	return ""
}

// Container for errors specific to comet.
type ErrorBag struct {
	Errors []*ParseError
}

func (b *ErrorBag) String() string {
	var sb strings.Builder
	for _, err := range b.Errors {
		sb.WriteString(err.Message)
		sb.WriteRune('\n')
	}
	return sb.String()
}

func (b *ErrorBag) Report(token lexer.Token, message string, params ...interface{}) {
	b.Errors = append(b.Errors, &ParseError{
		Message: fmt.Sprintf(message, params...),
		Token:   token,
	})
}

func (b *ErrorBag) HasAny() bool {
	return len(b.Errors) > 0
}

func newErrorBag() *ErrorBag {
	return &ErrorBag{
		make([]*ParseError, 0),
	}
}

// Functions of this type are going to be used to parse binary operations such as addition subtraction ...
// The first parameters is the already parsed left side of the operator and the function should parse

// The right side, and merge both of them and return them as a BinaryExpression
type binaryParseFunction func(Expression) Expression

// Function of this type are going to be use to parse unary operations such as ! -a +12
// The return value is Prefix Expression representing the parsed expression
type prefixParseFunction func() Expression

type Parser struct {
	lex *lexer.Lexer

	CurrentToken lexer.Token
	NextToken    lexer.Token

	Errors      *ErrorBag
	prefixFuncs map[lexer.TokenType]prefixParseFunction
	binaryFuncs map[lexer.TokenType]binaryParseFunction
}

func New(src string) *Parser {
	lex := lexer.NewLexer(src)
	parser := &Parser{
		lex:    lex,
		Errors: newErrorBag(),
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

	// Register functions to parse all operators that are of the form `op expresion`
	p.registerPrefixFunc(p.parseNumberLiteral, lexer.Number)
	p.registerPrefixFunc(p.parsePrefixExpression, lexer.Minus, lexer.Bang)
	p.registerPrefixFunc(p.parseIdentifier, lexer.Identifier)
	p.registerPrefixFunc(p.parseBoolean, lexer.True, lexer.False)
	p.registerPrefixFunc(p.parseParenthesisedExpression, lexer.OpenParent)
	p.registerPrefixFunc(p.parseStringLiteral, lexer.String)
	p.registerPrefixFunc(p.parseArrayLiteral, lexer.OpenBracket)
	p.registerPrefixFunc(p.parseNewCall, lexer.New)

	// Register functions to parse all operators that are of the form `expression op expresion`
	p.registerPrefixFunc(p.parseNumberLiteral, lexer.Number)
	p.registerBinaryFunc(p.parseArrayAccess, lexer.OpenBracket)
	p.registerBinaryFunc(p.parseBinaryExpression, lexer.Plus, lexer.Mul, lexer.Minus, lexer.Div,
		lexer.GT, lexer.GTE, lexer.LT, lexer.LTE, lexer.EQ, lexer.NEQ, lexer.Dot, lexer.DotDot)
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
	p.NextToken = p.lex.Next()
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
	case lexer.For:
		return p.parseForStatement()
	case lexer.Struct:
		return p.parseStructDeclaration()
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
		p.Errors.Report(p.CurrentToken, "Could not parse integer value %s", p.CurrentToken.Literal)
		return &NumberLiteral{0}
	}
	return &NumberLiteral{ActualValue: val}
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
	} else if p.NextToken.Type == lexer.Assign {
		assignExpression := &AssignExpression{
			VarName: p.CurrentToken.Literal,
		}
		p.expectNext(lexer.Assign)
		p.advance()
		assignExpression.Value = p.parseExpression()
		return assignExpression
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
		p.Errors.Report(p.CurrentToken, "Expected ')' got %s", p.CurrentToken.Literal)
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
		p.Errors.Report(p.CurrentToken, "No parsing function found for %s", p.CurrentToken.Literal)
		return nil
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
	for p.CurrentToken.Type != lexer.CloseBrace {
		if p.CurrentToken.Type == lexer.EOF {
			p.Errors.Report(p.CurrentToken, "Unexpected EOF")
			break
		}
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

func (p *Parser) parseForStatement() Statement {
	forStatement := &ForStatement{
		Value: &IdentifierExpression{
			Name: "__empty__",
		},
	}
	p.expectNext(lexer.Identifier)
	forStatement.Key = &IdentifierExpression{Name: p.CurrentToken.Literal}
	// If the next token is a comma, that means that there is a value identifier
	if p.NextToken.Type == lexer.Comma {
		p.advance()                    // at comma
		p.expectNext(lexer.Identifier) // at identifier
		forStatement.Value = &IdentifierExpression{Name: p.CurrentToken.Literal}
	}
	p.expectNext(lexer.In)
	p.advance()
	forStatement.Range = p.parseExpression()
	p.expectNext(lexer.OpenBrace)
	forStatement.Body = p.parseBlockStatement()
	return forStatement
}

func (p *Parser) advanceExpect(expected lexer.TokenType) {
	if p.CurrentToken.Type != expected {
		p.Errors.Report(p.CurrentToken, "Expected %s got %s instead", expected, p.CurrentToken.Literal)
	}
	p.advance()
}

func (p *Parser) expectNext(expected lexer.TokenType) {
	if p.NextToken.Type != expected {
		p.Errors.Report(p.CurrentToken, "Expected %s got %s instead", expected, p.CurrentToken.Literal)
	}
	p.advance()
}

func (p *Parser) parseStringLiteral() Expression {
	return &StringLiteral{Value: p.CurrentToken.Literal}
}

func (p *Parser) parseArrayLiteral() Expression {
	array := &ArrayLiteral{
		make([]Expression, 0),
	}
	p.advanceExpect(lexer.OpenBracket) // consume the first open bracket
	// In case of an empty array
	if p.CurrentToken.Type == lexer.CloseBracket {
		return array
	}
	// exp1, exp2, exp3]
	//     ^
	for p.CurrentToken.Type != lexer.CloseBracket {
		array.Elements = append(array.Elements, p.parseExpression())
		p.advance()
		if p.CurrentToken.Type == lexer.CloseBracket {
			break
		}
		p.advanceExpect(lexer.Comma)
	}
	// make sure that we consume the current token
	return array
}

func (p *Parser) parseArrayAccess(left Expression) Expression {
	indexAccess := &IndexAccess{Identifier: left}
	p.advance()
	indexAccess.Index = p.parseExpression()
	p.expectNext(lexer.CloseBracket)
	return indexAccess
}

func (p *Parser) parseStructDeclaration() Statement {
	structDec := &StructDeclarationStatement{}
	p.advance() // skip the struct keyword
	structDec.Name = p.CurrentToken.Literal
	p.advance() // skip the literal
	functions := make([]*FunctionStatement, 0)
	p.advanceExpect(lexer.OpenBrace) // skip the opening brace
	for {
		if p.CurrentToken.Type == lexer.CloseBrace {
			break
		}
		funcStatement := p.parseFunctionStatement()
		funcCasted, ok := funcStatement.(*FunctionStatement)
		if !ok {
			p.Errors.Report(p.CurrentToken, "Expected a function declaration")
			return structDec
		}
		functions = append(functions, funcCasted)
		p.advanceExpect(lexer.CloseBrace)
	}
	structDec.Methods = functions
	return structDec
}

func (p *Parser) parseNewCall() Expression {
	p.expectNext(lexer.Identifier)
	callExpr := &NewCallExpr{Type: p.CurrentToken.Literal}
	p.advance() // skip the type declaration
	callExpr.Args = p.parseCallArguments()
	return callExpr
}
