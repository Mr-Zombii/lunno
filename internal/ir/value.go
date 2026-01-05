package ir

import "fmt"

type Value interface {
	String() string
}

type LiteralValue struct {
	Kind  string
	Value interface{}
}

func (l *LiteralValue) String() string {
	return fmt.Sprintf("%v", l.Value)
}

type Variable struct {
	Name string
}

func (v *Variable) String() string {
	return v.Name
}

type Temporary struct {
	ID int
}

func (t *Temporary) String() string {
	return fmt.Sprintf("t%d", t.ID)
}
