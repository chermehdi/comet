package std

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/chermehdi/comet/parser"
)

// Type alias mapping some strings to types
type CometType string

const (
	IntType       = "INTEGER"
	BoolType      = "BOOLEAN"
	StrType       = "STR"
	ArrayType     = "ARRAY"
	FuncType      = "FUNCTION"
	ErrorType     = "ERROR"
	RangeType     = "RANGE"
	ObjType       = "OBJECT"
	ReturnWrapper = "ReturnWrapper"
	Nop           = "NOP"
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

type CometStr struct {
	Value string
	// Caching the size could prove beneficial, can't tell without benchmarks
	Size int
}

func (c *CometStr) Type() CometType {
	return StrType
}

func (c *CometStr) ToString() string {
	return fmt.Sprintf(`CometStr("%s")`, c.Value)
}

type CometArray struct {
	Length int
	Values []CometObject
}

func (c *CometArray) Type() CometType {
	return ArrayType
}

func (c *CometArray) ToString() string {
	var buf bytes.Buffer
	buf.WriteString("[")
	for i, obj := range c.Values {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(obj.ToString())
	}
	buf.WriteString("]")
	return buf.String()
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
	return ReturnWrapper
}

func (c *CometReturnWrapper) ToString() string {
	return fmt.Sprintf("CometWrapper(%s)", c.Value.ToString())
}

type CometFunc struct {
	Name   string
	Params []*parser.IdentifierExpression
	Body   *parser.BlockStatement
}

func (c *CometFunc) Type() CometType {
	return FuncType
}

func (c *CometFunc) ToString() string {
	// TODO(chermehdi): better ToString() representation for functions.
	return fmt.Sprintf("CometFunc")
}

type CometRange struct {
	From CometInt
	To   CometInt
}

func (c *CometRange) Type() CometType {
	return RangeType
}

func (c *CometRange) ToString() string {
	return fmt.Sprintf("CometRange(%d, %d)", c.From.Value, c.To.Value)
}

func CreateError(s string, params ...interface{}) CometObject {
	message := fmt.Sprintf(s, params...)
	return &CometError{
		message,
	}
}

// CometStruct represents a struct declaration in the comet language.
// Top level declaration
type CometStruct struct {
	Name    string
	Methods map[string]*CometFunc
}

// Add adds a method to a struct with performing sanity checks according to the
// language rules.
func (s *CometStruct) Add(fn *CometFunc) error {
	// As Comet is not typed it does not make sense to have method overloading
	// as A(1, 2) is equivalenet to A(1, 2, "") because no argument checking
	// is performed in the current implementation.
	_, found := s.Methods[fn.Name]
	if found {
		return errors.New(fmt.Sprintf("Method already exist with the name '%s' on '%s' struct", fn.Name, s.Name))
	}
	s.Methods[fn.Name] = fn
	return nil
}

// GetConstructor returns the CometFunc coresponding to the constructor of this
// type declaration.
// This should only be used at first time creating an instance of this type.
func (s *CometStruct) GetConstructor() (*CometFunc, bool) {
	fn, found := s.Methods["init"]
	return fn, found
}

// CometInstance is the object created from a given `Type`
type CometInstance struct {
	// Struct is the type definition for the given instance
	// every method on the struct declaration should take the current instance
	// as the this pointer reference.
	Struct *CometStruct
	// Fields represent the struct's state
	// Fields could be added at any point since this is a dynamic language
	Fields map[string]CometInstance
}

func (c *CometInstance) Type() CometType {
	return ObjType
}

func (c *CometInstance) ToString() string {
	// TODO: delegate to some special method on the Type declaration
	return fmt.Sprintf("CometInstance(Type=%s)", c.Struct.Name)
}

// NewInstance creates a new comet object
func NewInstance(typeDec *CometStruct) *CometInstance {
	return &CometInstance{
		Struct: typeDec,
		Fields: make(map[string]CometInstance, 0),
	}
}
