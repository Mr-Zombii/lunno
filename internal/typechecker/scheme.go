package typechecker

type Scheme struct {
	TypeVars []int
	Type     Type
}

func contains(slice []int, val int) bool {
	for _, x := range slice {
		if x == val {
			return true
		}
	}
	return false
}

func generalize(env *Env, typ Type) *Scheme {
	free := freeTypeVars(typ)
	envFree := envFreeTypeVars(env)
	var quantified []int
	for _, v := range free {
		if !contains(envFree, v) {
			quantified = append(quantified, v)
		}
	}
	return &Scheme{
		TypeVars: quantified,
		Type:     typ,
	}
}

func instantiate(s *Scheme, checker *Checker) Type {
	subst := Subst{}
	for _, id := range s.TypeVars {
		subst[id] = checker.freshVar()
	}
	return apply(s.Type, subst)
}

func freeTypeVars(t Type) []int {
	set := map[int]struct{}{}
	var collect func(Type)
	collect = func(tt Type) {
		switch ty := tt.(type) {
		case *TypeVar:
			set[ty.ID] = struct{}{}
		case *ListType:
			collect(ty.Element)
		case *FunctionType:
			for _, p := range ty.Parameters {
				collect(p)
			}
			collect(ty.Return)
		}
	}
	collect(t)
	var res []int
	for id := range set {
		res = append(res, id)
	}
	return res
}

func envFreeTypeVars(env *Env) []int {
	set := map[int]struct{}{}
	for _, t := range env.values {
		for _, id := range freeTypeVarsScheme(t) {
			set[id] = struct{}{}
		}
	}
	if env.parent != nil {
		for _, id := range envFreeTypeVars(env.parent) {
			set[id] = struct{}{}
		}
	}

	var res []int
	for id := range set {
		res = append(res, id)
	}
	return res
}

func freeTypeVarsScheme(scheme interface{}) []int {
	switch s := scheme.(type) {
	case *Scheme:
		ft := freeTypeVars(s.Type)
		quantified := map[int]struct{}{}
		for _, id := range s.TypeVars {
			quantified[id] = struct{}{}
		}
		var res []int
		for _, id := range ft {
			if _, ok := quantified[id]; !ok {
				res = append(res, id)
			}
		}
		return res
	default:
		return freeTypeVars(s.(Type))
	}
}
