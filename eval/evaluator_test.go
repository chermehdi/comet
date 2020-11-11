package eval

import (
	"fmt"
	"math"
	"testing"

	"github.com/chermehdi/comet/parser"
	"github.com/chermehdi/comet/std"
	"github.com/stretchr/testify/assert"
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
		Src        string
		AssertFunc func(*Evaluator)
	}{
		{
			Src: `var a = 1
				`,
			AssertFunc: func(ev *Evaluator) {
				a := assertFoundInScope(t, ev, "a", std.IntType)
				aValue := a.(*std.CometInt)
				assert.Equal(t, int64(1), aValue.Value)
			},
		},
		{
			Src: `var a = 1 * 2 + 1
				`,
			AssertFunc: func(ev *Evaluator) {
				a := assertFoundInScope(t, ev, "a", std.IntType)
				aValue := a.(*std.CometInt)
				assert.Equal(t, int64(3), aValue.Value)
			},
		},

		{
			Src: `var a = 1 * 2 + 1
                 var c = 10
				 var d = a * c
				`,
			AssertFunc: func(ev *Evaluator) {
				a := assertFoundInScope(t, ev, "a", std.IntType)
				aValue := a.(*std.CometInt)
				assert.Equal(t, int64(3), aValue.Value)

				c := assertFoundInScope(t, ev, "c", std.IntType)
				cValue := c.(*std.CometInt)
				assert.Equal(t, int64(10), cValue.Value)

				d := assertFoundInScope(t, ev, "d", std.IntType)
				dValue := d.(*std.CometInt)
				assert.Equal(t, int64(30), dValue.Value)
			},
		},
		{
			Src: `var a = "Hello world!"
				`,
			AssertFunc: func(ev *Evaluator) {
				a := assertFoundInScope(t, ev, "a", std.StrType)
				aValue := a.(*std.CometStr)
				assert.Equal(t, "Hello world!", aValue.Value)
				assert.Equal(t, 12, aValue.Size)
			},
		},
	}

	for _, test := range tests {
		evaluator := NewEvaluator()
		rootNode := parseOrDie(test.Src)
		evaluator.Eval(rootNode)
		test.AssertFunc(evaluator)
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
		{
			`c = 10
				`,
			"Identifier (c) is not bounded to any value, have you tried declaring it?",
		},
	}

	evaluator := NewEvaluator()
	for _, test := range tests {
		rootNode := parseOrDie(test.Src)
		v := evaluator.Eval(rootNode)
		assertError(t, v, test.ExpectedMessage)
	}
}

func TestEvaluator_Eval_StringOperations(t *testing.T) {

	tests := []struct {
		Src        string
		AssertFunc func(*Evaluator)
	}{
		{
			Src: `var a = "Hello " + "world!"`,
			AssertFunc: func(ev *Evaluator) {
				a := assertFoundInScope(t, ev, "a", std.StrType)
				aValue := a.(*std.CometStr)
				assert.Equal(t, "Hello world!", aValue.Value)
				assert.Equal(t, 12, aValue.Size)
			},
		},
		{
			Src: `
				var a = "Hello" * 3
				var b = 3 * "Hello"
			`,
			AssertFunc: func(ev *Evaluator) {
				a := assertFoundInScope(t, ev, "a", std.StrType)
				aValue := a.(*std.CometStr)
				assert.Equal(t, "HelloHelloHello", aValue.Value)
				assert.Equal(t, 15, aValue.Size)

				b := assertFoundInScope(t, ev, "b", std.StrType)
				bValue := b.(*std.CometStr)
				assert.Equal(t, "HelloHelloHello", bValue.Value)
				assert.Equal(t, 15, bValue.Size)
			},
		},
		{
			Src: `
				var a = "Hello" + true
				var b = true + "Hello"
				var c = false + "Hello"
				var d = "Hello" + false
			`,
			AssertFunc: func(ev *Evaluator) {
				a := assertFoundInScope(t, ev, "a", std.StrType)
				aValue := a.(*std.CometStr)
				assert.Equal(t, "Hellotrue", aValue.Value)
				assert.Equal(t, 9, aValue.Size)

				b := assertFoundInScope(t, ev, "b", std.StrType)
				bValue := b.(*std.CometStr)
				assert.Equal(t, "trueHello", bValue.Value)
				assert.Equal(t, 9, bValue.Size)

				c := assertFoundInScope(t, ev, "c", std.StrType)
				cValue := c.(*std.CometStr)
				assert.Equal(t, "falseHello", cValue.Value)
				assert.Equal(t, 9, cValue.Size)

				d := assertFoundInScope(t, ev, "d", std.StrType)
				dValue := d.(*std.CometStr)
				assert.Equal(t, "Hellofalse", dValue.Value)
				assert.Equal(t, 9, dValue.Size)
			},
		},
		{
			Src: `
				var a = "Hello" + 42 
				var b = 42 + "Hello"
				`,
			AssertFunc: func(ev *Evaluator) {
				a := assertFoundInScope(t, ev, "a", std.StrType)
				aValue := a.(*std.CometStr)
				assert.Equal(t, "Hello42", aValue.Value)
				assert.Equal(t, 7, aValue.Size)

				b := assertFoundInScope(t, ev, "b", std.StrType)
				bValue := b.(*std.CometStr)
				assert.Equal(t, "42Hello", bValue.Value)
				assert.Equal(t, 7, bValue.Size)
			},
		},
	}

	for _, test := range tests {
		evaluator := NewEvaluator()
		rootNode := parseOrDie(test.Src)
		evaluator.Eval(rootNode)
		test.AssertFunc(evaluator)
	}
}

func TestEvaluator_Eval_FunctionDeclarationTest(t *testing.T) {
	tests := []struct {
		Src        string
		AssertFunc func(*Evaluator)
	}{
		{
			Src: `func a() { return 1} 
				var c = a()	
            `,
			AssertFunc: func(evaluator *Evaluator) {
				obj := assertFoundInScope(t, evaluator, "a", std.FuncType)
				function, _ := obj.(*std.CometFunc)
				assert.Len(t, function.Params, 0)
			},
		},
		{
			Src: `func a(p1, p2) {}`,
			AssertFunc: func(evaluator *Evaluator) {
				obj := assertFoundInScope(t, evaluator, "a", std.FuncType)
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

func TestEvaluator_Eval_FunctionCallTest(t *testing.T) {
	tests := []struct {
		Src        string
		AssertFunc func(*Evaluator)
	}{
		{
			Src: `func a() { return 1} 
				var c = a()	
            `,
			AssertFunc: func(evaluator *Evaluator) {
				assertFoundInScope(t, evaluator, "a", std.FuncType)
				c := assertFoundInScope(t, evaluator, "c", std.IntType)
				value := c.(*std.CometInt)
				assert.Equal(t, int64(1), value.Value)
			},
		},
		{
			Src: `func a() { return 1} 
				func b(v, f) { return v * f() }
                var c = b(2, a)
            `,
			AssertFunc: func(evaluator *Evaluator) {
				assertFoundInScope(t, evaluator, "a", std.FuncType)
				assertFoundInScope(t, evaluator, "b", std.FuncType)
				c := assertFoundInScope(t, evaluator, "c", std.IntType)
				value := c.(*std.CometInt)
				assert.Equal(t, int64(2), value.Value)
			},
		},
		{
			Src: `func a() { return 1} 
				func b(v, f) { return v * f() }
                var c1 = b(2, a)
                var c2 = b(2, a)
				var comp = c1 == c2
            `,
			AssertFunc: func(evaluator *Evaluator) {
				assertFoundInScope(t, evaluator, "a", std.FuncType)
				assertFoundInScope(t, evaluator, "b", std.FuncType)
				assertFoundInScope(t, evaluator, "c1", std.IntType)
				assertFoundInScope(t, evaluator, "c2", std.IntType)
				comp := assertFoundInScope(t, evaluator, "comp", std.BoolType)
				value := comp.(*std.CometBool)
				assert.Equal(t, true, value.Value)
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

func TestEvaluator_Eval_EvaluateForStatement(t *testing.T) {
	tests := []struct {
		Src        string
		AssertFunc func(*Evaluator)
	}{
		{
			Src: `	
				var a = 10
				for i in 0..2 { 
                  for j in 0..2 {
					a = a + i * j
                  }
				}
            `,
			AssertFunc: func(evaluator *Evaluator) {
				a := assertFoundInScope(t, evaluator, "a", std.IntType)
				value := a.(*std.CometInt)
				assert.Equal(t, int64(19), value.Value)
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

func assertFoundInScope(t *testing.T, ev *Evaluator, name string, expectedType std.CometType) std.CometObject {
	obj, found := ev.Scope.Lookup(name)
	assert.True(t, found)
	assert.True(t, expectedType == obj.Type())
	return obj
}

func parseOrDie(s string) parser.Node {
	return parser.New(s).Parse()
}

func ExampleBuiltinPrintf() {
	// (anouard24): I'm not pround of this
	// but i'm still learning Golang
	// its enough for now
	tests := []string{
		`printf("Hi")`,
		`println()
		printf("Hi %d", 2021)`,
		`println()
		printf("% 7d", 5)`,
		`println()
		printf("%s", "Test")`,
		`println()
		printf("%8s", "Test")`,
		`println()
		printf("%t", true)`,
		`println()
		printf("%t", false)`,
		`println()
		printf("%04d%9s%t", 7, "Comet ", true)`,
	}
	evaluator := NewEvaluator()
	for _, test := range tests {
		rootNode := parseOrDie(test)
		evaluator.Eval(rootNode)
	}

	// Output:
	// Hi
	// Hi 2021
	//       5
	// Test
	//     Test
	// true
	// false
	// 0007   Comet true
}
