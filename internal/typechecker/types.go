package typechecker

import (
	"fmt"
	"lunno/internal/parser"
)

type Type interface {
	String() string
	isType()
}

type (
	IntType    struct{}
	FloatType  struct{}
	BoolType   struct{}
	StringType struct{}
	CharType   struct{}
	UnitType   struct{}

	ListType struct {
		Element Type
	}

	FunctionType struct {
		Parameters []Type
		Return     Type
	}

	TypeVar struct {
		ID int
	}
)

func (*IntType) isType()      {}
func (*FloatType) isType()    {}
func (*BoolType) isType()     {}
func (*StringType) isType()   {}
func (*CharType) isType()     {}
func (*UnitType) isType()     {}
func (*ListType) isType()     {}
func (*FunctionType) isType() {}
func (*TypeVar) isType()      {}

func (*IntType) String() string {
	return "int"
}

func (*FloatType) String() string {
	return "float"
}

func (*BoolType) String() string {
	return "bool"
}

func (*StringType) String() string {
	return "string"
}

func (*CharType) String() string {
	return "char"
}

func (*UnitType) String() string {
	return "unit"
}

func (t *ListType) String() string {
	return "list(" + t.Element.String() + ")"
}

func (t *FunctionType) String() string {
	s := "fn("
	for i, p := range t.Parameters {
		if i > 0 {
			s += ", "
		}
		s += p.String()
	}
	s += ") -> " + t.Return.String()
	return s
}

func (t *TypeVar) String() string {
	return fmt.Sprintf("T%d", t.ID)
}

func (checker *Checker) resolveType(typ parser.TypeNode) Type {
	switch t := typ.(type) {
	case *parser.SimpleType:
		switch t.Name {
		case "int":
			return &IntType{}
		case "float":
			return &FloatType{}
		case "bool":
			return &BoolType{}
		case "string":
			return &StringType{}
		case "char":
			return &CharType{}
		case "unit":
			return &UnitType{}
		default:
			return checker.freshVar()
		}
	case *parser.ListType:
		return &ListType{
			Element: checker.resolveType(t.Element),
		}
	case *parser.FunctionType:
		params := make([]Type, len(t.Parameters))
		for i, p := range t.Parameters {
			params[i] = checker.resolveType(p)
		}
		var ret Type = &UnitType{}
		if t.Return != nil {
			ret = checker.resolveType(t.Return)
		}
		return &FunctionType{
			Parameters: params,
			Return:     ret,
		}
	}
	return nil
}
