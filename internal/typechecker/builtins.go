package typechecker

func registerBuiltins(env *Env) {
	env.set("_builtin_print", &FunctionType{
		Parameters: []Type{
			&TypeVar{ID: -1},
		},
		Return: &UnitType{},
	})
}
