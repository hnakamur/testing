package bench

import "testing"

func BenchmarkScanFor(b *testing.B) {
	s := &scanner{s: "EEEEEEEEEEEEEEEEND"}
	for i := 0; i < b.N; i++ {
		var c byte
		for isE(c) {
			c = s.next()
		}
	}
}

// 7x slow
func BenchmarkScanFunc(b *testing.B) {
	s := &scanner{s: "EEEEEEEEEEEEEEEEND"}
	for i := 0; i < b.N; i++ {
		s.skipIf(isE)
	}
}
