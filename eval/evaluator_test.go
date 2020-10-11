package eval

import (
	"fmt"
	"github.com/chermehdi/comet/parser"
	"github.com/chermehdi/comet/std"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestEvaluator_Eval_Integers(t *testing.T) {
	tests := []struct {
		Token    string
		Expected int64
	}{
		{
			"-1",
			-1,
		},
		{
			"10",
			10,
		},
		{
			fmt.Sprintf("%d", math.MaxInt64),
			math.MaxInt64,
		},
		{
			"1 + 1",
			2,
		},
		{
			"1 - 1",
			0,
		},
		{
			"2 * 15",
			30,
		},
		{
			"15 / 3",
			5,
		},
		{
			"1 + 2 * 3",
			7,
		},
		{
			"1 * -2",
			-2,
		},
		{
			"(1)",
			1,
		},
	}

	evaluator := NewEvaluator()
	for _, test := range tests {
		rootNode := parseOrDie(test.Token)
		v := evaluator.Eval(rootNode)
		assertInteger(t, v, test.Expected)
	}
}

func TestEvaluator_Eval_Booleans(t *testing.T) {
	tests := []struct {
		Token    string
		Expected bool
	}{
		{
			"true",
			true,
		},
		{
			"false",
			false,
		},
		{
			"!true",
			false,
		},
		{
			"!!true",
			true,
		},
		{
			"true == true",
			true,
		},
		{
			"true != false",
			true,
		},
		{
			"true == false",
			false,
		},
		{
			"true != false",
			true,
		},
		{
			"1 == true",
			false,
		},
		{
			"1 != true",
			true,
		},
	}

	evaluator := NewEvaluator()
	for _, test := range tests {
		rootNode := parseOrDie(test.Token)
		v := evaluator.Eval(rootNode)
		assertBoolean(t, v, test.Expected)
	}
}

func TestEvaluator_Eval_Conditionals(t *testing.T) {
	tests := []struct {
		Src      string
		Expected bool
	}{
		{
			`if(true) { true }`,
			true,
		},
		{
			`if(false) { true } else { false }`,
			false,
		},
		{
			`if(1 < 2) { true }`,
			true,
		},
		{
			`if(1 == 1) { true }`,
			true,
		},
	}

	evaluator := NewEvaluator()
	for _, test := range tests {
		rootNode := parseOrDie(test.Src)
		v := evaluator.Eval(rootNode)
		assertBoolean(t, v, test.Expected)
	}
}

func TestEvaluator_Eval_ReturnStatement(t *testing.T) {
	tests := []struct {
		Src      string
		Expected int64
	}{
		{
			"return 10",
			10,
		},
		{
			`9 * 9
				return 10`,
			10,
		},
		{
			`9 * 9
				return 10
				8 + 10`,
			10,
		},
		{
			`if(true) {
					if (true) {
						return 10
					}
					return 1
				}`, 10,
		},
	}

	evaluator := NewEvaluator()
	for _, test := range tests {
		rootNode := parseOrDie(test.Src)
		v := evaluator.Eval(rootNode)
		assertInteger(t, v, test.Expected)
	}
}

func TestEvaluator_Eval_Errors(t *testing.T) {
	tests := []struct {
		Src              string
		ExpectedErrorMsg string
	}{
		{
			"1 + true",
			"Cannot apply operator + on given types INTEGER and BOOLEAN",
		},
		{
			"1 * true",
			"Cannot apply operator * on given types INTEGER and BOOLEAN",
		},
		{
			"1 - true",
			"Cannot apply operator - on given types INTEGER and BOOLEAN",
		},
		{
			"true > 1",
			"Cannot apply operator > on given types BOOLEAN and INTEGER",
		},
		{
			"true < 1",
			"Cannot apply operator < on given types BOOLEAN and INTEGER",
		},
		{
			"-true",
			"Cannot apply operator (-) on none INTEGER type BOOLEAN",
		},
		{
			"-false",
			"Cannot apply operator (-) on none INTEGER type BOOLEAN",
		},
		{
			"!1",
			"Cannot apply operator (!) on none BOOLEAN type INTEGER",
		},
		{
			`
				if (true) {
					!1
					false
				}
				`,
			"Cannot apply operator (!) on none BOOLEAN type INTEGER",
		},
	}

	evaluator := NewEvaluator()
	for _, test := range tests {
		rootNode := parseOrDie(test.Src)
		v := evaluator.Eval(rootNode)
		assertError(t, v, test.ExpectedErrorMsg)
	}
}

func TestEvaluator_Eval_Declarations(t *testing.T) {
	tests := []struct {
		Src      string
		Expected int64
	}{
		{
			`var a = 1
				a
				`,
			1,
		},
		{
			`var a = 1 * 2 + 1
				a
				`,
			3,
		},

		{
			`var a = 1 * 2 + 1
                 var c = 10
				 var d = a * c
				 d
				`,
			30,
		},
	}

	evaluator := NewEvaluator()
	for _, test := range tests {
		rootNode := parseOrDie(test.Src)
		v := evaluator.Eval(rootNode)
		assertInteger(t, v, test.Expected)
	}
}

func TestEvaluator_Eval_DeclarationError(t *testing.T) {
	tests := []struct {
		Src             string
		ExpectedMessage string
	}{
		{
			`var a = b * 10 
				`,
			"Identifier (b) is not bounded to any value, have you tried declaring it?",
		},
	}

	evaluator := NewEvaluator()
	for _, test := range tests {
		rootNode := parseOrDie(test.Src)
		v := evaluator.Eval(rootNode)
		assertError(t, v, test.ExpectedMessage)
	}
}

func TestEvaluator_Eval_FunctionDeclarationTest(t *testing.T) {
	tests := []struct {
		Src        string
		AssertFunc func(*Evaluator)
	}{
		{
			Src: `func a() {}`,
			AssertFunc: func(evaluator *Evaluator) {
				obj, found := evaluator.Scope.Lookup("a")
				assert.True(t, found)
				assert.True(t, std.FuncType == obj.Type())
				function, _ := obj.(*std.CometFunc)
				assert.Len(t, function.Params, 0)
			},
		},
		{
			Src: `func a(p1, p2) {}`,
			AssertFunc: func(evaluator *Evaluator) {
				obj, found := evaluator.Scope.Lookup("a")
				assert.True(t, found)
				assert.True(t, std.FuncType == obj.Type())
				function, _ := obj.(*std.CometFunc)
				assert.Len(t, function.Params, 2)
				assert.Equal(t, "p1", function.Params[0].Name)
				assert.Equal(t, "p2", function.Params[1].Name)
			},
		},
	}
	evaluator := NewEvaluator()
	for _, test := range tests {
		rootNode := parseOrDie(test.Src)
		evaluator.Eval(rootNode)
		test.AssertFunc(evaluator)
	}
}

func assertError(t *testing.T, v std.CometObject, ExpectedErrorMsg string) {
	err, ok := v.(*std.CometError)
	assert.True(t, ok)
	assert.Equal(t, ExpectedErrorMsg, err.Message)
}

func assertBoolean(t *testing.T, v std.CometObject, expected bool) {
	boolean, ok := v.(*std.CometBool)
	assert.True(t, ok)
	assert.Equal(t, expected, boolean.Value)
}

func assertInteger(t *testing.T, v std.CometObject, expected int64) {
	integer, ok := v.(*std.CometInt)
	assert.True(t, ok)
	assert.Equal(t, expected, integer.Value)
}

func parseOrDie(s string) parser.Node {
	return parser.New(s).Parse()
}
