package object

// Env holds the execution environment.
type Env struct {
	store map[string]Object
}

// NewEnv creates an Env.
func NewEnv() *Env {
	s := make(map[string]Object)
	return &Env{store: s}
}

// Get looks up a value by the name it's bound to.
func (e *Env) Get(name string) (Object, bool) {
	val, ok := e.store[name]
	return val, ok
}

// Set associates (binds) a given object with a given name.
func (e *Env) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
