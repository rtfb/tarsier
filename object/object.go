package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/rtfb/tarsier/ast"
)

// The constant values for Type.
const (
	ObjTypeInteger     = "INTEGER"
	ObjTypeString      = "STRING"
	ObjTypeBoolean     = "BOOLEAN"
	ObjTypeNull        = "NULL"
	ObjTypeReturnValue = "RETURN_VALUE"
	ObjTypeError       = "ERROR"
	ObjTypeFunction    = "FUNCTION"
	ObjTypeBuiltin     = "BUILTIN"
	ObjTypeArray       = "ARRAY"
	ObjTypeHash        = "HASH"
)

// HashKey contains a hash sum.
type HashKey struct {
	Type  Type
	Value uint64
}

// Hashable describes an interface for objects that can be hashed.
type Hashable interface {
	HashKey() HashKey
}

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

// HashKey returns a hashed value for an int.
func (i *Integer) HashKey() HashKey {
	return HashKey{
		Type:  i.Type(),
		Value: uint64(i.Value),
	}
}

// String is an implementation for a string Object type.
type String struct {
	Value string
}

// Type implements Object.
func (s *String) Type() Type {
	return ObjTypeString
}

// Inspect implements Object.
func (s *String) Inspect() string {
	return s.Value
}

// HashKey returns a hashed value for a string.
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{
		Type:  s.Type(),
		Value: h.Sum64(),
	}
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

// HashKey returns a hashed value for a bool.
func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{
		Type:  b.Type(),
		Value: value,
	}
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

// ReturnValue encapsulates the value after the return statement.
type ReturnValue struct {
	Value Object
}

// Type implements Object.
func (rv *ReturnValue) Type() Type {
	return ObjTypeReturnValue
}

// Inspect implements Object.
func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

// Error represents an execution error.
// TODO: extend this with stack trace and source code line:col (the latter will
// require extending the lexer to attach the line:col info to the tokens).
type Error struct {
	Message string
}

// Type implements Object.
func (e *Error) Type() Type {
	return ObjTypeError
}

// Inspect implements Object.
func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}

// Function represents a function object.
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Env
}

// Type implements Object.
func (f *Function) Type() Type {
	return ObjTypeFunction
}

// Inspect implements Object.
func (f *Function) Inspect() string {
	var out bytes.Buffer
	params := make([]string, len(f.Parameters))
	for i, p := range f.Parameters {
		params[i] = p.String()
	}
	out.WriteString("fn(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")
	return out.String()
}

// BuiltinFunction is a signature for implementation of any built-in function.
type BuiltinFunction func(args ...Object) Object

// Builtin represents a language-provided built-in function.
type Builtin struct {
	Fn BuiltinFunction
}

// Type implements Object.
func (b *Builtin) Type() Type {
	return ObjTypeBuiltin
}

// Inspect implements Object.
func (b *Builtin) Inspect() string {
	return "builtin function"
}

// Array is an object containing an ordered list of other objects.
type Array struct {
	Elements []Object
}

// Type implements Object.
func (a *Array) Type() Type {
	return ObjTypeArray
}

// Inspect implements Object.
func (a *Array) Inspect() string {
	var out bytes.Buffer
	elements := make([]string, len(a.Elements))
	for i, e := range a.Elements {
		elements[i] = e.Inspect()
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}

// HashPair represents a
type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

// Type implements Object.
func (h *Hash) Type() Type {
	return ObjTypeHash
}

// Inspect implements Object.
func (h *Hash) Inspect() string {
	var out bytes.Buffer
	pairs := make([]string, 0, len(h.Pairs))
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(),
			pair.Value.Inspect()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}
