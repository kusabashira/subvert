package main

import (
	"reflect"
	"regexp"
	"strings"
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
			continue
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
			continue
		}
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("%q: got %q, want %q",
				test.src, actual, expect)
		}
	}
}

var genReplacementTests = []struct {
	from        string
	to          string
	replacement []map[string]string
}{
	// one branch
	{
		from: "abc",
		to:   "def",
		replacement: []map[string]string{
			map[string]string{
				"abc": "def",
			},
		},
	},
	{
		from: "abcdef",
		to:   "ghijkl",
		replacement: []map[string]string{
			map[string]string{
				"abcdef": "ghijkl",
			},
		},
	},

	// multiple branches
	{
		from: "a,b",
		to:   "b,a",
		replacement: []map[string]string{
			map[string]string{
				"a": "b",
				"b": "a",
			},
		},
	},
	{
		from: "a,,b,c",
		to:   "d,e,f,g",
		replacement: []map[string]string{
			map[string]string{
				"a": "d",
				"":  "e",
				"b": "f",
				"c": "g",
			},
		},
	},
	{
		from: ",a",
		to:   "a,",
		replacement: []map[string]string{
			map[string]string{"": "a", "a": ""},
		},
	},
	{
		from: "a,b,c",
		to:   ",d,",
		replacement: []map[string]string{
			map[string]string{
				"a": "",
				"b": "d",
				"c": "",
			},
		},
	},

	// multiple groups
	{
		from: "a/b",
		to:   "c/d",
		replacement: []map[string]string{
			map[string]string{
				"a": "c",
			},
			map[string]string{
				"b": "d",
			},
		},
	},
	{
		from: "a//b/c",
		to:   "d/e/f/g",
		replacement: []map[string]string{
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
		from: "a,b/c",
		to:   "d,e/f",
		replacement: []map[string]string{
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
		from: "/a",
		to:   "a/",
		replacement: []map[string]string{
			map[string]string{
				"": "a",
			},
			map[string]string{
				"a": "",
			},
		},
	},
	{
		from: "/a/",
		to:   "b/c/d",
		replacement: []map[string]string{
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

	// special chars
	{
		from: "( , )",
		to:   "(,)",
		replacement: []map[string]string{
			map[string]string{
				"( ": "(",
				" )": ")",
			},
		},
	},
	{
		from: "^*/|$",
		to:   "[+/?]",
		replacement: []map[string]string{
			map[string]string{
				"^*": "[+",
			},
			map[string]string{
				"|$": "?]",
			},
		},
	},
}

func TestGenReplacement(t *testing.T) {
	for _, test := range genReplacementTests {
		expect := test.replacement
		actual, err := newReplacement(test.from, test.to)
		if err != nil {
			t.Errorf("newReplacement(%q, %q) returns %q, want nil",
				test.from, test.to, err)
			continue
		}
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("%q, %q: got %q, want %q",
				test.from, test.to, actual, expect)
		}
	}
}

var replaceTests = []struct {
	from string
	to   string
	src  string
	dst  string
}{
	// one branch
	{
		from: "abc",
		to:   "def",
		src:  "foo bar",
		dst:  "foo bar",
	},
	{
		from: "abc",
		to:   "def",
		src:  "abc def",
		dst:  "def def",
	},
	{
		from: "a",
		to:   "b",
		src:  "a b c a b c",
		dst:  "b b c b b c",
	},

	// multiple branches
	{
		from: "abc,def",
		to:   "def,abc",
		src:  "abc def",
		dst:  "def abc",
	},
	{
		from: "a,b,c,d",
		to:   "e,f,g,h",
		src:  "d c b a",
		dst:  "h g f e",
	},
	{
		from: "a, ",
		to:   " ,a",
		src:  "a a a",
		dst:  " a a ",
	},

	// multiple groups
	{
		from: "a/b",
		to:   "c/d",
		src:  "aa ab ac ad",
		dst:  "aa cd ac ad",
	},
	{
		from: "a//b/c",
		to:   "d/e/f/g",
		src:  "abc bca cab",
		dst:  "defg bca cab",
	},
	{
		from: "dog,cat/s",
		to:   "cat,dog/s",
		src:  "cats cats dogs dogs cats",
		dst:  "dogs dogs cats cats dogs",
	},
}

func TestReplace(t *testing.T) {
	for _, test := range replaceTests {
		r, err := NewReplacer(test.from, test.to, false)
		if err != nil {
			t.Errorf("NewReplacer(%q, %q) returns %q, want nil",
				test.from, test.to, err)
			continue
		}

		expect := test.dst
		actual := r.ReplaceAll(test.src)
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("Replacer{%q, %q}: %q: got %q, want %q",
				test.from, test.to, test.src, actual, expect)
		}
	}
}

func BenchmarkStringsReplace(b *testing.B) {
	src := strings.Repeat("aaa bbb\n", 1000)
	rep := strings.NewReplacer("aaa", "bbb", "bbb", "aaa")
	for i := 0; i < b.N; i++ {
		rep.Replace(src)
	}
}

func BenchmarkReplacerReplace(b *testing.B) {
	src := strings.Repeat("aaa bbb\n", 1000)
	from, to := "aaa,bbb", "bbb,aaa"
	r, err := NewReplacer(from, to, false)
	if err != nil {
		b.Fatalf("NewReplacer(%q, %q, false) returns %q, want nil",
			from, to, err)
	}
	for i := 0; i < b.N; i++ {
		r.ReplaceAll(src)
	}
}
