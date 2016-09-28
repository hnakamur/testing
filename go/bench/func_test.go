package bench

import "testing"

func funcfunc(v byte) func() byte {
	return func() byte {
		return v
	}
}

type obj struct {
	v byte
}

func (r *obj) f() byte {
	return r.v
}

func objfunc(v byte) func() byte {
	o := &obj{v}
	return o.f
}

func BenchmarkFuncFunc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		funcfunc('a')()
	}
}

func BenchmarkObjFunc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		objfunc('a')()
	}
}
