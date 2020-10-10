package std

import (
	"fmt"
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
				return CreateError("Expected 0 or 1 arguments, got %s.", len(args))
			}
			// This works if args[0] is a string, int or boolean
			// Maybe we should only allow this for the defined types, but for the time being it's not required.
			fmt.Println(extractPrimitive(args[0]))
			return NopInstance
		},
	},
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
