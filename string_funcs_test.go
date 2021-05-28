// -----------------------------------------------------------------------------
// github.com/balacode/udpt                              /[string_funcs_test.go]
// (c) balarabe@protonmail.com                                      License: MIT
// -----------------------------------------------------------------------------

package udpt

import (
	"testing"
)

// to run all tests in this file:
// go test -v -run Test_string_*

// -----------------------------------------------------------------------------

// getPart(s, prefix, suffix string) string
//
// go test -run Test_string_getPart_
//
func Test_string_getPart_(t *testing.T) {
	for _, it := range []struct {
		s      string
		prefix string
		suffix string
		//
		expect string
	}{
		{"", "", "", ""},
		//
		// both prefix and suffix are blank: return 's' as it is
		{"name:cat;", "", "", "name:cat;"},
		//
		// prefix is blank, suffix given: return everything before the suffix
		{"name:cat;", "", ";", "name:cat"},
		//
		// prefix given, suffix is blank: return everything after the prefix
		{"name:cat;", "name:", "", "cat;"},
		//
		// both prefix and suffix specified: return the substring between them
		{"name:cat;", "name:", ";", "cat"},
		//
		// non-existent prefix: return a blank string
		{"name:cat;", "age:", "", ""},
		//
		// non-existent suffix: return a blank string
		{"name:cat;", "name:", ".", ""},
		//
		// additional test
		{"xyz class::sum; 123", "class::", ";", "sum"},
	} {
		got := getPart(it.s, it.prefix, it.suffix)
		if got != it.expect {
			t.Errorf("0xE2C73F"+" getPart(%#v, %#v, %#v)"+
				"\n expect: %#v"+
				"\n    got: %#v",
				it.s, it.prefix, it.suffix, it.expect, got)
		}
	}
}

// joinArgs(tag string, a ...interface{}) string
//
// go test -run Test_string_joinArgs_
//
func Test_string_joinArgs_(t *testing.T) {
	for _, it := range []struct {
		tag string
		a   []interface{}
		//
		expect string
	}{
		{"", nil, ""},
		{" ", nil, ""},
		{" tag ", nil, "tag"},
		{"tag", nil, "tag"},
		{"tag ", nil, "tag"},
		{"tag", []interface{}{"a", 1, "b", 2, 3, "c"}, "tag a 1 b 2 3 c"},
		{"", []interface{}{"abc", 123, "de", 34}, "abc 123 de 34"},
	} {
		got := joinArgs(it.tag, it.a...)
		if got != it.expect {
			t.Errorf("0xE7B77D"+" joinArgs(%#v, %#v)"+
				"\n expect: %#v"+
				"\n    got: %#v",
				it.tag, it.a, it.expect, got)
		}
	}
}

// padf(minLength int, format string, a ...interface{}) string
//
// go test -run Test_string_padf_
//
func Test_string_padf_(t *testing.T) {
	for _, it := range []struct {
		minLength int
		format    string
		a         []interface{}
		//
		expect string
	}{
		{0, "", nil, ""},
		{0, "%s", []interface{}{"abc"}, "abc"},
		{6, "%s", []interface{}{"abc"}, "abc   "},
	} {
		got := padf(it.minLength, it.format, it.a...)
		if got != it.expect {
			t.Errorf("0xEA7E4F"+" padf(%#v, %#v, %#v)"+
				"\n expect: %#v"+
				"\n    got: %#v",
				it.minLength, it.format, it.a, it.expect, got)
		}
	}
}

// end
