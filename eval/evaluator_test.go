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
	}

	evaluator := New()
	for _, test := range tests {
		rootNode := parseOrDie(test.Token)
		v := evaluator.Eval(rootNode)
		assertBoolean(t, v, test.Expected)
	}
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
