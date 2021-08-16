package inject

import (
	"testing"
)

var g Container

type Inline struct {
	S *Sub `inject:""`
}

type Test struct {
	*Inline `inject:"private"`

	Name string
}

type Sub struct {
	A string
}

func TestContainer_PopulateTarget(t *testing.T) {
	g.Provide(&Object{Value: &Sub{A: "a"}})

	tt := &Test{}
	if err := g.PopulateTarget(tt); err != nil {
		t.Fatal(err)
	}

	t.Log(tt.S)
}

func TestContainer_Resolve(t *testing.T) {
	t1 := &Test{Name: "a"}
	g.Provide(&Object{Value: t1})

	t2 :=&Test{}
	g.Resolve(t2)

	t.Log(t2)
}