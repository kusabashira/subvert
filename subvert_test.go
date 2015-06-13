package main

import (
	"reflect"
	"regexp"
	"testing"
)

var genMatcherTests = []struct {
	src string
	dst *regexp.Regexp
}{
	// one branch
	{`abc`, regexp.MustCompile(`(abc)`)},
	{`abcdef`, regexp.MustCompile(`(abcdef)`)},

	// multiple branches
	{`a,b`, regexp.MustCompile(`(a|b)`)},
	{`a,,b,c`, regexp.MustCompile(`(a||b|c)`)},
	{`a,bc,def`, regexp.MustCompile(`(a|bc|def)`)},
	{`,a`, regexp.MustCompile(`(|a)`)},
	{`a,`, regexp.MustCompile(`(a|)`)},
	{`,a,`, regexp.MustCompile(`(|a|)`)},

	// use escape
	{`a\,b`, regexp.MustCompile(`(a,b)`)},
	{`a\,bc\,def`, regexp.MustCompile(`(a,bc,def)`)},

	// multiple branches with escape
	{`a\,b,c`, regexp.MustCompile(`(a,b|c)`)},
	{`a,bc\,def`, regexp.MustCompile(`(a|bc,def)`)},

	// regexp quote
	{`a+b`, regexp.MustCompile(`(a\+b)`)},
	{`(a|bc)*def`, regexp.MustCompile(`(\(a\|bc\)\*def)`)},

	// unquote special values
	{`a\\bc`, regexp.MustCompile("(a\\\\bc)")},
	{`a\tb\,c`, regexp.MustCompile("(a\tb,c)")},
	{`a\tbc\n\ndef`, regexp.MustCompile("(a\tbc\n\ndef)")},

	// multiple sequenses
	{`a/b`, regexp.MustCompile("(a)(b)")},
	{`a//b/c`, regexp.MustCompile(`(a)()(b)(c)`)},
	{`a/bc/def`, regexp.MustCompile("(a)(bc)(def)")},
	{`a,b/c`, regexp.MustCompile("(a|b)(c)")},
	{`/a`, regexp.MustCompile(`()(a)`)},
	{`a/`, regexp.MustCompile(`(a)()`)},
	{`/a/`, regexp.MustCompile(`()(a)()`)},

	// multiple sequenses with escape
	{`a/b\/c`, regexp.MustCompile("(a)(b/c)")},
	{`a/\/bc\//def`, regexp.MustCompile("(a)(/bc/)(def)")},
	{`a\,b,c/d,e\/f`, regexp.MustCompile("(a,b|c)(d|e/f)")},
}

func TestGenMatcher(t *testing.T) {
	for _, test := range genMatcherTests {
		expect := test.dst
		actual, err := newMatcher(test.src)
		if err != nil {
			t.Errorf("NewSubvert(%q) returns %q, want nil",
				test.src, err)
		}
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("%q: got %q, want %q",
				test.src, actual, expect)
		}
	}
}

var genReplacementTests = []struct {
	srcFrom string
	srcTo   string
	dst     []map[string]string
}{
	// one branch
	{
		"abc",
		"def",
		[]map[string]string{
			map[string]string{
				"abc": "def",
			},
		},
	},
	{
		"abcdef",
		"ghijkl",
		[]map[string]string{
			map[string]string{
				"abcdef": "ghijkl",
			},
		},
	},

	// multiple branches
	{
		"a,b",
		"b,a",
		[]map[string]string{
			map[string]string{
				"a": "b",
				"b": "a",
			},
		},
	},
	{
		"a,,b,c",
		"d,e,f,g",
		[]map[string]string{
			map[string]string{
				"a": "d",
				"":  "e",
				"b": "f",
				"c": "g",
			},
		},
	},
	{
		",a",
		"a,",
		[]map[string]string{
			map[string]string{"": "a", "a": ""},
		},
	},
	{
		"a,b,c",
		",d,",
		[]map[string]string{
			map[string]string{
				"a": "",
				"b": "d",
				"c": "",
			},
		},
	},
}

func TestGenReplacement(t *testing.T) {
	for _, test := range genReplacementTests {
		expect := test.dst
		actual, err := newReplacement(test.srcFrom, test.srcTo)
		if err != nil {
			t.Errorf("NewSubvert(%q, %q) returns %q, want nil",
				test.srcFrom, test.srcTo, err)
		}
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("%q, %q: got %q, want %q",
				test.srcFrom, test.srcTo, actual, expect)
		}
	}
}
