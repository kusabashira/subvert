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

	// multiple groups
	{`a/b`, regexp.MustCompile("(a)(b)")},
	{`a//b/c`, regexp.MustCompile(`(a)()(b)(c)`)},
	{`a/bc/def`, regexp.MustCompile("(a)(bc)(def)")},
	{`a,b/c`, regexp.MustCompile("(a|b)(c)")},
	{`/a`, regexp.MustCompile(`()(a)`)},
	{`a/`, regexp.MustCompile(`(a)()`)},
	{`/a/`, regexp.MustCompile(`()(a)()`)},

	// multiple groups with escape
	{`a/b\/c`, regexp.MustCompile("(a)(b/c)")},
	{`a/\/bc\//def`, regexp.MustCompile("(a)(/bc/)(def)")},
	{`a\,b,c/d,e\/f`, regexp.MustCompile("(a,b|c)(d|e/f)")},
}

func TestGenMatcher(t *testing.T) {
	for _, test := range genMatcherTests {
		expect := test.dst
		actual, err := newMatcher(test.src, false)
		if err != nil {
			t.Errorf("newMatcher(%q) returns %q, want nil",
				test.src, err)
		}
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("%q: got %q, want %q",
				test.src, actual, expect)
		}
	}
}

var genMatcherWithBoundaryTests = []struct {
	src string
	dst *regexp.Regexp
}{
	{`abc`, regexp.MustCompile(`\b(abc)\b`)},
	{`a,b`, regexp.MustCompile(`\b(a|b)\b`)},
	{`a\,b,c`, regexp.MustCompile(`\b(a,b|c)\b`)},
	{`a/b`, regexp.MustCompile(`\b(a)(b)\b`)},
	{`a/bc/def`, regexp.MustCompile(`\b(a)(bc)(def)\b`)},
	{`a,b/c`, regexp.MustCompile(`\b(a|b)(c)\b`)},
	{`a\,b,c/d,e\/f`, regexp.MustCompile(`\b(a,b|c)(d|e/f)\b`)},
}

func TestGenMatcherWithBoundary(t *testing.T) {
	for _, test := range genMatcherWithBoundaryTests {
		expect := test.dst
		actual, err := newMatcher(test.src, true)
		if err != nil {
			t.Errorf("newMatcher(%q) returns %q, want nil",
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

	// multiple groups
	{
		"a/b",
		"c/d",
		[]map[string]string{
			map[string]string{
				"a": "c",
			},
			map[string]string{
				"b": "d",
			},
		},
	},
	{
		"a//b/c",
		"d/e/f/g",
		[]map[string]string{
			map[string]string{
				"a": "d",
			},
			map[string]string{
				"": "e",
			},
			map[string]string{
				"b": "f",
			},
			map[string]string{
				"c": "g",
			},
		},
	},
	{
		"a,b/c",
		"d,e/f",
		[]map[string]string{
			map[string]string{
				"a": "d",
				"b": "e",
			},
			map[string]string{
				"c": "f",
			},
		},
	},
	{
		"/a",
		"a/",
		[]map[string]string{
			map[string]string{
				"": "a",
			},
			map[string]string{
				"a": "",
			},
		},
	},
	{
		"/a/",
		"b/c/d",
		[]map[string]string{
			map[string]string{
				"": "b",
			},
			map[string]string{
				"a": "c",
			},
			map[string]string{
				"": "d",
			},
		},
	},
}

func TestGenReplacement(t *testing.T) {
	for _, test := range genReplacementTests {
		expect := test.dst
		actual, err := newReplacement(test.srcFrom, test.srcTo)
		if err != nil {
			t.Errorf("newReplacement(%q, %q) returns %q, want nil",
				test.srcFrom, test.srcTo, err)
		}
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("%q, %q: got %q, want %q",
				test.srcFrom, test.srcTo, actual, expect)
		}
	}
}

var replaceTests = []struct {
	srcFrom string
	srcTo   string
	srcText string
	dst     string
}{
	// one branch
	{
		"abc",
		"def",
		"foo bar",
		"foo bar",
	},
	{
		"abc",
		"def",
		"abc def",
		"def def",
	},
	{
		"a",
		"b",
		"a b c a b c",
		"b b c b b c",
	},

	// multiple branches
	{
		"abc,def",
		"def,abc",
		"abc def",
		"def abc",
	},
	{
		"a,b,c,d",
		"e,f,g,h",
		"d c b a",
		"h g f e",
	},
	{
		"a, ",
		" ,a",
		"a a a",
		" a a ",
	},

	// multiple groups
	{
		"a/b",
		"c/d",
		"aa ab ac ad",
		"aa cd ac ad",
	},
	{
		"a//b/c",
		"d/e/f/g",
		"abc bca cab",
		"defg bca cab",
	},
	{
		"dog,cat/s",
		"cat,dog/s",
		"cats cats dogs dogs cats",
		"dogs dogs cats cats dogs",
	},
}

func TestReplace(t *testing.T) {
	for _, test := range replaceTests {
		r, err := NewReplacer(test.srcFrom, test.srcTo)
		if err != nil {
			t.Errorf("NewReplacer(%q, %q) returns %q, want nil",
				test.srcFrom, test.srcTo, err)
		}

		expect := test.dst
		actual := r.ReplaceAll(test.srcText)
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("Replacer{%q, %q}: %q: got %q, want %q",
				test.srcFrom, test.srcTo, test.srcText, actual, expect)
		}
	}
}
