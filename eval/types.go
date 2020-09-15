package eval

import "fmt"

// Type alias mapping some strings to types
type CometType string

const (
	IntType  = "INTEGER"
	BoolType = "BOOLEAN"
)

// Every object (or primitive) in the comet programming language will be representated
// As an instance of this interface.
type CometObject interface {
	// Returns the type of this instance, see CometType for details about available/possible types
	Type() CometType

	// A string representation For Debugging / REPL purposes
	ToString() string
}

type CometInt struct {
	Value int64
}

func (i *CometInt) Type() CometType {
	return IntType
}

func (i *CometInt) ToString() string {
	return fmt.Sprintf("CometInt(%d)", i.Value)
}

type CometBool struct {
	Value bool
}

func (b *CometBool) Type() CometType {
	return BoolType
}

func (b *CometBool) ToString() string {
	return fmt.Sprintf("CometBool(%v)", b.Value)
}
