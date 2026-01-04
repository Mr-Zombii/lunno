package typechecker

import "fmt"

type Subst map[int]Type

type Checker struct {
	env     *Env
	nextVar int
	errors  []error
}

func (checker *Checker) freshVar() *TypeVar {
	tv := &TypeVar{ID: checker.nextVar}
	checker.nextVar++
	return tv
}

func apply(t Type, s Subst) Type {
	switch t := t.(type) {
	case *TypeVar:
		if r, ok := s[t.ID]; ok {
			return apply(r, s)
		}
		return t
	case *ListType:
		return &ListType{Element: apply(t.Element, s)}
	case *FunctionType:
		params := make([]Type, len(t.Parameters))
		for i, p := range t.Parameters {
			params[i] = apply(p, s)
		}
		return &FunctionType{
			Parameters: params,
			Return:     apply(t.Return, s),
		}
	default:
		return t
	}
}

func unify(a, b Type, s Subst) error {
	a = apply(a, s)
	b = apply(b, s)
	if av, ok := a.(*TypeVar); ok {
		s[av.ID] = b
		return nil
	}
	if bv, ok := b.(*TypeVar); ok {
		s[bv.ID] = a
		return nil
	}
	switch a := a.(type) {
	case *IntType, *FloatType, *BoolType,
		*StringType, *CharType, *UnitType:
		if a.String() != b.String() {
			return fmt.Errorf("type mismatch: %s vs %s", a, b)
		}
		return nil
	case *ListType:
		bt, ok := b.(*ListType)
		if !ok {
			return fmt.Errorf("expected list, got %s", b)
		}
		return unify(a.Element, bt.Element, s)
	case *FunctionType:
		bt, ok := b.(*FunctionType)
		if !ok {
			return fmt.Errorf("expected function, got %s", b)
		}
		if len(a.Parameters) != len(bt.Parameters) {
			return fmt.Errorf("arity mismatch")
		}
		for i := range a.Parameters {
			if err := unify(a.Parameters[i], bt.Parameters[i], s); err != nil {
				return err
			}
		}
		return unify(a.Return, bt.Return, s)
	}
	return fmt.Errorf("cannot unify %T and %T", a, b)
}
