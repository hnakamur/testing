package bench

import (
	"strconv"
	"testing"
)

func BenchmarkSwitch(b *testing.B) {
	var a byte
	for i := 0; i < b.N; i++ {
		switch i % 10 {
		case 0:
			a = 'a'
		case 1:
			a = 'b'
		case 2:
			a = 'c'
		case 3:
			a = 'd'
		case 4:
			a = 'e'
		case 5:
			a = 'f'
		case 6:
			a = 'g'
		case 7:
			a = 'h'
		case 8:
			a = 'i'
		case 9:
			a = 'j'
		}
	}
	_ = a
}

func BenchmarkSwitchString(b *testing.B) {
	var a byte
	for i := 0; i < b.N; i++ {
		switch strconv.Itoa(i % 10) {
		case "0":
			a = 'a'
		case "1":
			a = 'b'
		case "2":
			a = 'c'
		case "3":
			a = 'd'
		case "4":
			a = 'e'
		case "5":
			a = 'f'
		case "6":
			a = 'g'
		case "7":
			a = 'h'
		case "8":
			a = 'i'
		case "9":
			a = 'j'
		}
	}
	_ = a
}

func BenchmarkLookupTable(b *testing.B) {
	v := [10]byte{
		0: 'a',
		1: 'b',
		2: 'c',
		3: 'd',
		4: 'e',
		5: 'f',
		6: 'g',
		7: 'h',
		8: 'i',
		9: 'j',
	}
	var a byte
	for i := 0; i < b.N; i++ {
		a = v[i%10]
	}
	_ = a
}

// 5x slower
func BenchmarkLookupFunc(b *testing.B) {
	v := [10]func() (v byte){
		0: func() (v byte) { return 'a' },
		1: func() (v byte) { return 'b' },
		2: func() (v byte) { return 'c' },
		3: func() (v byte) { return 'd' },
		4: func() (v byte) { return 'e' },
		5: func() (v byte) { return 'f' },
		6: func() (v byte) { return 'g' },
		7: func() (v byte) { return 'h' },
		8: func() (v byte) { return 'i' },
		9: func() (v byte) { return 'j' },
	}
	var a byte
	for i := 0; i < b.N; i++ {
		a = v[i%10]()
	}
	_ = a
}

func BenchmarkMapFunc(b *testing.B) {
	v := map[string]func() (v byte){
		"0": func() (v byte) { return 'a' },
		"1": func() (v byte) { return 'b' },
		"2": func() (v byte) { return 'c' },
		"3": func() (v byte) { return 'd' },
		"4": func() (v byte) { return 'e' },
		"5": func() (v byte) { return 'f' },
		"6": func() (v byte) { return 'g' },
		"7": func() (v byte) { return 'h' },
		"8": func() (v byte) { return 'i' },
		"9": func() (v byte) { return 'j' },
	}
	var a byte
	for i := 0; i < b.N; i++ {
		a = v[strconv.Itoa(i%10)]()
	}
	_ = a
}

