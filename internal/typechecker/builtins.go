package typechecker

func registerBuiltins(env *Env) {
	tv := &TypeVar{ID: 0}
	s := &Scheme{
		TypeVars: []int{tv.ID},
		Type: &FunctionType{
			Parameters: []Type{&StringType{}},
			Return:     &UnitType{},
		},
	}
	env.set("builtin_print", s)
}
