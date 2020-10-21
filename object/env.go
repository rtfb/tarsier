package object

// Env holds the execution environment.
type Env struct {
	store map[string]Object
	outer *Env
}

// NewEnv creates an Env.
func NewEnv() *Env {
	s := make(map[string]Object)
	return &Env{store: s}
}

// NewEnclosedEnv creates a nested environment inside a given outer one.
func NewEnclosedEnv(outer *Env) *Env {
	env := NewEnv()
	env.outer = outer
	return env
}

// Get looks up a value by the name it's bound to.
func (e *Env) Get(name string) (Object, bool) {
	val, ok := e.store[name]
	if !ok && e.outer != nil {
		val, ok = e.outer.Get(name)
	}
	return val, ok
}

// Set associates (binds) a given object with a given name.
func (e *Env) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
