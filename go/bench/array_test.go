package bench

import "testing"

func BenchmarkString(b *testing.B) {
	str := "foobar"
	var slice string
	for i := 0; i < b.N; i++ {
		slice = str[0:3]
	}
	b.Logf("%v", slice)
}

func BenchmarkArrayString(b *testing.B) {
	array := []string{"foo", "bar"}
	var slice string
	for i := 0; i < b.N; i++ {
		slice = array[0]
	}
	b.Logf("%v", slice)
}
