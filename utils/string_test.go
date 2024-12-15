// Package utils string unit tests.
package utils_test

import (
	"testing"

	"github.com/callerobertsson/resty/utils"
)

func TestContainsString(t *testing.T) {
	cases := []struct {
		s  string
		ss []string
		e  bool
	}{
		{"", []string{}, false},
		{"nop", []string{}, false},
		{"abc", []string{"abc"}, true},
		{"abc", []string{"cba"}, false},
		{"abc", []string{"cba", "abc"}, true},
		{"def", []string{"cba", "abc"}, false},
	}

	for i, c := range cases {
		got := utils.ContainsString(c.ss, c.s)

		if c.e != got {
			t.Errorf("Case %d failed: ContainsString(%v, %q) == %v, expected %v", i, c.ss, c.s, got, c.e)
		}
	}
}
