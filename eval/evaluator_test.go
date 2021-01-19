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

func TestEvaluator_Eval_EvaluateArrayDeclaration(t *testing.T) {
	tests := []struct {
		Src        string
		AssertFunc func(*Evaluator)
	}{
		{
			Src: `	
				var a = []
            `,
			AssertFunc: func(evaluator *Evaluator) {
				a := assertFoundInScope(t, evaluator, "a", std.ArrayType)
				array := a.(*std.CometArray)
				assert.Equal(t, 0, array.Length)
			},
		},
		{
			Src: `	
				var a = [1, 2, 3]
            `,
			AssertFunc: func(evaluator *Evaluator) {
				a := assertFoundInScope(t, evaluator, "a", std.ArrayType)
				array := a.(*std.CometArray)
				assert.Equal(t, 3, array.Length)

				assertInteger(t, array.Values[0], 1)
				assertInteger(t, array.Values[1], 2)
				assertInteger(t, array.Values[2], 3)
			},
		},
		{
			Src: `	
				var a = [[], [1, 2]]
            `,
			AssertFunc: func(evaluator *Evaluator) {
				a := assertFoundInScope(t, evaluator, "a", std.ArrayType)
				array := a.(*std.CometArray)
				assert.Equal(t, 2, array.Length)

				assert.True(t, array.Values[0].Type() == std.ArrayType)

				arrv0 := array.Values[0].(*std.CometArray)
				assert.Equal(t, 0, arrv0.Length)

				assert.True(t, array.Values[1].Type() == std.ArrayType)

				arrv1 := array.Values[1].(*std.CometArray)
				assert.Equal(t, 2, arrv1.Length)
				assertInteger(t, arrv1.Values[0], 1)
				assertInteger(t, arrv1.Values[1], 2)
			},
		},
		{
			Src: `	
				var a = ["comet", "42"]
            `,
			AssertFunc: func(evaluator *Evaluator) {
				a := assertFoundInScope(t, evaluator, "a", std.ArrayType)
				array := a.(*std.CometArray)
				assert.Equal(t, 2, array.Length)

				assertStr(t, array.Values[0], "comet")
				assertStr(t, array.Values[1], "42")
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

func TestEvaluator_Eval_EvaluateArrayAccess(t *testing.T) {
	tests := []struct {
		Src        string
		AssertFunc func(*Evaluator)
	}{
		{
			Src: `	
				var a = [0, 1]
				var b = a[0]
            `,
			AssertFunc: func(evaluator *Evaluator) {
				b := assertFoundInScope(t, evaluator, "b", std.IntType)
				bValue := b.(*std.CometInt)
				assert.Equal(t, int64(0), bValue.Value)
			},
		},
		{
			Src: `	
				var a = ["12"]
				var b = a[0]
            `,
			AssertFunc: func(evaluator *Evaluator) {
				b := assertFoundInScope(t, evaluator, "b", std.StrType)
				bValue := b.(*std.CometStr)
				assert.Equal(t, "12", bValue.Value)
			},
		},
		{
			Src: `	
				func getArray() {
					return [1, 2, 3]
				}
				var b = getArray()[0]
            `,
			AssertFunc: func(evaluator *Evaluator) {
				b := assertFoundInScope(t, evaluator, "b", std.IntType)
				bValue := b.(*std.CometInt)
				assert.Equal(t, int64(1), bValue.Value)
			},
		},
		{
			Src: `	
				var b = [1, 2, 3][2]
            `,
			AssertFunc: func(evaluator *Evaluator) {
				b := assertFoundInScope(t, evaluator, "b", std.IntType)
				bValue := b.(*std.CometInt)
				assert.Equal(t, int64(3), bValue.Value)
			},
		},
		{
			Src: `	
				var a = [[1, 42], [2, 3]]
				var b = a[0][1] 
            `,
			AssertFunc: func(evaluator *Evaluator) {
				b := assertFoundInScope(t, evaluator, "b", std.IntType)
				bValue := b.(*std.CometInt)
				assert.Equal(t, int64(42), bValue.Value)
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

func TestEvaluator_Eval_EvaluateStructDeclaration(t *testing.T) {
	tests := []struct {
		Src        string
		AssertFunc func(*Evaluator)
	}{
		{
			Src: `	
						struct a { }
            `,
			AssertFunc: func(evaluator *Evaluator) {
				s := assertFoundType(t, evaluator, "a")
				assert.Equal(t, "a", s.Name)
				assert.Equal(t, 0, len(s.Methods))
			},
		},
		{
			Src: `	
						struct a { 
							func init() { 
								var temp = 10
							}
						}
            `,
			AssertFunc: func(evaluator *Evaluator) {
				s := assertFoundType(t, evaluator, "a")
				assert.Equal(t, "a", s.Name)
				assert.Equal(t, 1, len(s.Methods))
			},
		},
		{
			Src: `	
						struct a { 
							func testa() { 
							}
							func testa(a) { 
							}
						}
            `,
			AssertFunc: func(evaluator *Evaluator) {
				assertNotFoundType(t, evaluator, "a")
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

func TestEvaluator_Eval_EvaluateInstanceCreation(t *testing.T) {
	tests := []struct {
		Src        string
		AssertFunc func(*Evaluator)
	}{
		{
			Src: `	
						struct A { 
						}
						var a = new A()
            `,
			AssertFunc: func(evaluator *Evaluator) {
				tp := assertFoundType(t, evaluator, "A")
				s := assertFoundInScope(t, evaluator, "a", std.ObjType)
				p, ok := s.(*std.CometInstance)
				assert.True(t, ok)
				assert.Equal(t, p.Struct, tp)
				assert.Equal(t, 0, len(p.Fields))
			},
		},
		{
			Src: `	
						struct A { 
						}
						var a = new A()
						var b = new A()
            `,
			AssertFunc: func(evaluator *Evaluator) {
				tp := assertFoundType(t, evaluator, "A")
				sa := assertFoundInScope(t, evaluator, "a", std.ObjType)
				sb := assertFoundInScope(t, evaluator, "b", std.ObjType)
				pa, oka := sa.(*std.CometInstance)
				pb, okb := sb.(*std.CometInstance)
				assert.True(t, oka)
				assert.True(t, okb)
				assert.Equal(t, pa.Struct, tp)
				assert.Equal(t, pb.Struct, tp)
				assert.Equal(t, 0, len(pa.Fields))
				assert.Equal(t, 0, len(pb.Fields))
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

func TestEvaluator_Eval_EvaluateMethodCall(t *testing.T) {
	tests := []struct {
		Src        string
		AssertFunc func(*Evaluator)
	}{
		{
			Src: `	
						struct A { 
							func init() { }
							func hello() {
								return 12
							}
						}
						var a = new A()
 						var res = a.hello()
            `,
			AssertFunc: func(evaluator *Evaluator) {
				tp := assertFoundType(t, evaluator, "A")
				s := assertFoundInScope(t, evaluator, "a", std.ObjType)
				p, ok := s.(*std.CometInstance)
				assert.True(t, ok)
				assert.Equal(t, p.Struct, tp)
				assert.Equal(t, 0, len(p.Fields))
				res := assertFoundInScope(t, evaluator, "res", std.IntType)
				val, ok := res.(*std.CometInt)
				assert.True(t, ok)
				assert.Equal(t, int64(12), val.Value)
			},
		},
		{
			Src: `	
						struct A { 
							func add(a, b) {
								return a + b
							}
						}
						var a = new A()
 						var res = a.add(10, 20)
            `,
			AssertFunc: func(evaluator *Evaluator) {
				tp := assertFoundType(t, evaluator, "A")
				s := assertFoundInScope(t, evaluator, "a", std.ObjType)
				p, ok := s.(*std.CometInstance)
				assert.True(t, ok)
				assert.Equal(t, p.Struct, tp)
				assert.Equal(t, 0, len(p.Fields))
				res := assertFoundInScope(t, evaluator, "res", std.IntType)
				val, ok := res.(*std.CometInt)
				assert.True(t, ok)
				assert.Equal(t, int64(30), val.Value)
			},
		},
		{
			Src: `	
						struct A { 
							func mul(a, b) {
								return a * b.get()	
							}
						}
						struct B { func get() { return 12 } } 
						var a = new A()
						var b = new B()
 						var res = a.mul(3, b)
            `,
			AssertFunc: func(evaluator *Evaluator) {
				tp := assertFoundType(t, evaluator, "A")
				s := assertFoundInScope(t, evaluator, "a", std.ObjType)
				p, ok := s.(*std.CometInstance)
				assert.True(t, ok)
				assert.Equal(t, p.Struct, tp)
				assert.Equal(t, 0, len(p.Fields))
				tp = assertFoundType(t, evaluator, "B")
				s = assertFoundInScope(t, evaluator, "b", std.ObjType)
				p, ok = s.(*std.CometInstance)
				assert.True(t, ok)
				assert.Equal(t, p.Struct, tp)
				assert.Equal(t, 0, len(p.Fields))
				res := assertFoundInScope(t, evaluator, "res", std.IntType)
				val, ok := res.(*std.CometInt)
				assert.True(t, ok)
				assert.Equal(t, int64(36), val.Value)
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

func assertStr(t *testing.T, v std.CometObject, expected string) {
	str, ok := v.(*std.CometStr)
	assert.True(t, ok)
	assert.Equal(t, expected, str.Value)
}

func assertFoundInScope(t *testing.T, ev *Evaluator, name string, expectedType std.CometType) std.CometObject {
	obj, found := ev.Scope.Lookup(name)
	assert.True(t, found)
	assert.True(t, expectedType == obj.Type())
	return obj
}

func assertFoundType(t *testing.T, ev *Evaluator, name string) *std.CometStruct {
	obj, found := ev.Types[name]
	assert.True(t, found)
	return obj
}

func assertNotFoundType(t *testing.T, ev *Evaluator, name string) {
	_, found := ev.Types[name]
	assert.False(t, found)
}

func parseOrDie(s string) parser.Node {
	return parser.New(s).Parse()
}
