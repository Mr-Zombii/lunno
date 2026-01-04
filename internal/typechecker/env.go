package typechecker

type Env struct {
	parent *Env
	values map[string]Type
}

func newEnv(parent *Env) *Env {
	return &Env{
		parent: parent,
		values: map[string]Type{},
	}
}

func (env *Env) get(name string) (Type, bool) {
	if t, ok := env.values[name]; ok {
		return t, true
	}
	if env.parent != nil {
		return env.parent.get(name)
	}
	return nil, false
}

func (env *Env) set(name string, t Type) {
	env.values[name] = t
}
