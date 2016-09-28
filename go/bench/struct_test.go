package bench

import "testing"

type str1 struct {
	v string
}

func BenchmarkStruct(b *testing.B) {
	array := []string{"foo", "bar"}
	var str str1
	for i := 0; i < b.N; i++ {
		str = str1{array[0]}
	}
	b.Logf("%v", str)
}

func BenchmarkStructPointer(b *testing.B) {
	array := []string{"foo", "bar"}
	var str *str1
	for i := 0; i < b.N; i++ {
		str = &str1{array[0]}
	}
	b.Logf("%v", str)
}
