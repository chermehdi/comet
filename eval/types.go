package eval

import "fmt"

// Type alias mapping some strings to types
type CometType string

const (
	IntegerType = "INTEGER"
	BooleanType = "BOOLEAN"
)

// Every object (or primitive) in the comet programming language will be representated
// As an instance of this interface.
type CometObject interface {
	// Returns the type of this instance, see CometType for details about available/possible types
	Type() CometType

	// A string representation For Debugging / REPL purposes
	ToString() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() CometType {
	return IntegerType
}

func (i *Integer) ToString() string {
	return fmt.Sprintf("Integer(%d)", i.Value)
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() CometType {
	return BooleanType
}

func (b *Boolean) ToString() string {
	return fmt.Sprintf("Boolean(%v)", b.Value)
}

