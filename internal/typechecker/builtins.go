package typechecker

func registerBuiltins(env *Env) {
	tv := &TypeVar{ID: 0}
	s := &Scheme{
		TypeVars: []int{tv.ID},
		Type: &FunctionType{
			Parameters: []Type{tv},
			Return:     &UnitType{},
		},
	}
	env.set("_builtin_print", s)
}
