package eval

import "fmt"

// Type alias mapping some strings to types
type CometType string

const (
	IntType   = "INTEGER"
	BoolType  = "BOOLEAN"
	ErrorType = "ERROR"
	Nop       = "NOP"
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

type CometError struct {
	Message string
}

func (c *CometError) Type() CometType {
	return ErrorType
}

func (c *CometError) ToString() string {
	return fmt.Sprintf("Comet error: \n\n\t%s", c.Message)
}

type NopObject struct{}

func (n *NopObject) Type() CometType {
	return Nop
}

func (n *NopObject) ToString() string {
	return "CometNop"
}

type CometReturnWrapper struct {
	Value CometObject
}

func (c *CometReturnWrapper) Type() CometType {
	return "CometReturnType"
}

func (c *CometReturnWrapper) ToString() string {
	return fmt.Sprintf("CometWrapper(%s)", c.Value.ToString())
}
