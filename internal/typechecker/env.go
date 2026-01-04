package typechecker

type Env struct {
	parent *Env
	values map[string]*Scheme
}

func newEnv(parent *Env) *Env {
	return &Env{
		parent: parent,
		values: map[string]*Scheme{},
	}
}

func (env *Env) get(name string) (*Scheme, bool) {
	if s, ok := env.values[name]; ok {
		return s, true
	}
	if env.parent != nil {
		return env.parent.get(name)
	}
	return nil, false
}

func (env *Env) set(name string, s *Scheme) {
	env.values[name] = s
}
