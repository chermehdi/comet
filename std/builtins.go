package std

import (
	"fmt"
	"strconv"
)

type Callback func(args ...CometObject) CometObject

type Builtin struct {
	Name string
	Func Callback
}

// Global builtin singletons
var (
	TrueObject  = &CometBool{true}
	FalseObject = &CometBool{false}
	NopInstance = &NopObject{}
)

var Builtins = []*Builtin{
	{
		Name: "printf",
		Func: func(args ...CometObject) CometObject {
			if len(args) == 0 {
				// Just an empty line call
				return CreateError("Expected 1 or more arguments, got none.")
			}
			if args[0].Type() != StrType {
				return CreateError("First argument expected to be CometString got '%s' instead", args[0].Type())
			}
			transArgs := make([]interface{}, 0)
			for i := 1; i < len(args); i++ {
				transArgs = append(transArgs, extractPrimitive(args[i]))
			}
			format := args[0].(*CometStr)
			fmt.Printf(format.Value, transArgs)
			return NopInstance
		},
	},
	{
		Name: "println",
		Func: func(args ...CometObject) CometObject {
			if len(args) == 0 {
				fmt.Println()
				return NopInstance
			}
			if len(args) != 1 {
				return CreateError("Expected 0 or 1 arguments, got %d.", len(args))
			}
			// This works if args[0] is a string, int or boolean
			// Maybe we should only allow this for the defined types, but for the time being it's not required.
			fmt.Println(extractPrimitive(args[0]))
			return NopInstance
		},
	},
}

// Standard library to convert any object type to a string value.
// Newly added types should add their string conversion implementation as well.
func ToString(object CometObject) *CometStr {
	switch n := object.(type) {
	case *CometStr:
		return n
	case *CometBool:
		return &CometStr{Value: strconv.FormatBool(n.Value), Size: 4}
	case *CometInt:
		value := strconv.FormatInt(n.Value, 10)
		return &CometStr{Value: value, Size: len(value)}
	case *CometFunc:
		value := n.ToString()
		return &CometStr{Value: value, Size: len(value)}
	case *CometError:
		value := n.Message
		return &CometStr{Value: value, Size: len(value)}
	default:
		panic("All types should have been exhausted!!")
	}
}

func extractPrimitive(object CometObject) interface{} {
	switch n := object.(type) {
	case *CometStr:
		return n.Value
	case *CometBool:
		return n.Value
	case *CometInt:
		return n.Value
	default:
		return object
	}
}
