package eval

import (
	"fmt"
	"github.com/chermehdi/comet/parser"
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

	evaluator := New()
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

	evaluator := New()
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

	evaluator := New()
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

	evaluator := New()
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

	evaluator := New()
	for _, test := range tests {
		rootNode := parseOrDie(test.Src)
		v := evaluator.Eval(rootNode)
		assertError(t, v, test.ExpectedErrorMsg)
	}
}

func assertError(t *testing.T, v CometObject, ExpectedErrorMsg string) {
	error, ok := v.(*CometError)
	assert.True(t, ok)
	assert.Equal(t, ExpectedErrorMsg, error.Message)
}

func assertBoolean(t *testing.T, v CometObject, expected bool) {
	boolean, ok := v.(*CometBool)
	assert.True(t, ok)
	assert.Equal(t, expected, boolean.Value)
}

func assertInteger(t *testing.T, v CometObject, expected int64) {
	integer, ok := v.(*CometInt)
	assert.True(t, ok)
	assert.Equal(t, expected, integer.Value)
}

func parseOrDie(s string) parser.Node {
	return parser.New(s).Parse()
}
