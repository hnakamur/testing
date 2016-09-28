package bench

type scanner struct {
	s string
	i int
}

func (s *scanner) next() (c byte) {
	if s.i < len(s.s) {
		c = s.s[s.i]
		s.i++
	}
	return
}

func (s *scanner) skipIf(f func(byte) bool) {
	var c byte
	for f(c) {
		c = s.next()
	}
}

func isE(c byte) bool {
	return c == 'e' || c == 'E'
}
