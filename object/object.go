package object

import "fmt"

// The constant values for Type.
const (
	ObjTypeInteger = "INTEGER"
	ObjTypeBoolean = "BOOLEAN"
	ObjTypeNull    = "NULL"
)

// Type is an identifier for a type of an object.
type Type string

// Object represents an object in the program.
type Object interface {
	Type() Type
	Inspect() string
}

// Integer is an implementation for an integer Object type.
type Integer struct {
	Value int64
}

// Type implements Object.
func (i *Integer) Type() Type {
	return ObjTypeInteger
}

// Inspect implements Object.
func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

// Boolean is an implementation for a boolean Object type.
type Boolean struct {
	Value bool
}

// Type implements Object.
func (b *Boolean) Type() Type {
	return ObjTypeBoolean
}

// Inspect implements Object.
func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

// Null is an implementation of a null Object.
type Null struct{}

// Type implements Object.
func (n *Null) Type() Type {
	return ObjTypeNull
}

// Inspect implements Object.
func (n *Null) Inspect() string {
	return "null"
}
