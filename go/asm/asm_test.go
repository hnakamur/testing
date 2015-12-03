package asm_test

import (
	"testing"

	"github.com/harukasan/testing/go/asm"
)

func TestRet(t *testing.T) {
	tests := []struct {
		x uint64
	}{
		{0},
		{1},
		{0xFFFFFFFFFFFFFFFF},
	}

	for _, test := range tests {
		if got := asm.Ret(test.x); got != test.x {
			t.Errorf("%x: %x", test.x, got)
		}
	}
}

func TestSum(t *testing.T) {
	tests := []struct {
		x uint64
		y uint64
		a uint64
	}{
		{1, 1, 2},
		{1, 2, 3},
		{2, 2, 4},
	}

	for _, test := range tests {
		if got := asm.Sum(test.x, test.y); got != test.a {
			t.Errorf("%d + %d: got %x, expect %x", test.x, test.y, got, test.a)
		}
	}
}
