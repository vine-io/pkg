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

type TestWithName struct {
	S *Sub `inject:"sub"`
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

	t2 := &Test{}
	g.Resolve(t2)

	t.Log(t2)
}

func TestContainer_ResolveByName(t *testing.T) {
	sub := &Sub{A: "aa"}
	if err := g.Provide(&Object{Value: sub, Name: "sub"}); err != nil {
		t.Fatal(err)
	}

	tn := &TestWithName{}
	if err := g.Provide(&Object{Value: tn}); err != nil {
		t.Fatal(err)
	}

	if err := g.Populate(); err != nil {
		t.Fatal(err)
	}

	t.Log(tn.S)
}
