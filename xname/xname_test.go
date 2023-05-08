package xname

import "testing"

func TestGen(t *testing.T) {
	s1 := Gen()
	s2 := Gen()

	if len(s1) != len(s2) {
		t.Fatalf("Gen() failed")
	}

	t.Log(s1, s2)

	s1 = Gen(C(20), Lowercase(), Digit())
	t.Log(s1)
}

func TestGen6(t *testing.T) {
	s1 := Gen6()

	t.Log(s1)
}

func BenchmarkGen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Gen(C(10))
	}
}
