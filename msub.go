package main

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	group     = regexp.MustCompile(`(?:[^/\\]|\\.)*`)
	branch    = regexp.MustCompile(`(?:[^,\\]|\\.)*`)
	backslash = regexp.MustCompile(`\\(.)`)
)

func parseExpr(expr string) (tree [][]string, err error) {
	gls := group.FindAllString(expr, -1)
	for gi := 0; gi < len(gls); gi++ {
		bls := branch.FindAllString(gls[gi], -1)
		for bi := 0; bi < len(bls); bi++ {
			bls[bi] = backslash.ReplaceAllString(bls[bi], "$1")
		}
		tree = append(tree, bls)
	}
	return tree, nil
}

func newMatcher(expr string, useBoundary bool) (m *regexp.Regexp, err error) {
	tree, err := parseExpr(expr)
	if err != nil {
		return nil, err
	}

	sls := make([]string, len(tree))
	for gi, bls := range tree {
		for bi := 0; bi < len(bls); bi++ {
			bls[bi] = regexp.QuoteMeta(bls[bi])
		}
		sls[gi] = "(" + strings.Join(bls, "|") + ")"
	}

	re := strings.Join(sls, "")
	if useBoundary {
		re = `\b` + re + `\b`
	}
	return regexp.Compile(re)
}

func newReplacement(exprFrom, exprTo string) ([]map[string]string, error) {
	from, err := parseExpr(exprFrom)
	if err != nil {
		return nil, err
	}
	to, err := parseExpr(exprTo)
	if err != nil {
		return nil, err
	}

	if len(from) != len(to) {
		return nil, fmt.Errorf("mismatch the number of group")
	}

	r := make([]map[string]string, len(from))
	for gi := 0; gi < len(from); gi++ {
		if len(from[gi]) != len(to[gi]) {
			return nil, fmt.Errorf("mismatch the number of branch at group[%d]", gi)
		}

		r[gi] = make(map[string]string)
		for bi := 0; bi < len(from[gi]); bi++ {
			src, dst := from[gi][bi], to[gi][bi]
			if _, exist := r[gi][src]; exist {
				return nil, fmt.Errorf("group[%d] has duplicate items", gi)
			}
			r[gi][src] = dst
		}
	}
	return r, nil
}

type Replacer struct {
	matcher     *regexp.Regexp
	replacement []map[string]string
}

func NewReplacer(from, to string, useBoundary bool) (r *Replacer, err error) {
	r = &Replacer{}

	r.matcher, err = newMatcher(from, useBoundary)
	if err != nil {
		return nil, err
	}
	r.replacement, err = newReplacement(from, to)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Replacer) ReplaceAll(s string) string {
	return r.matcher.ReplaceAllStringFunc(s, func(t string) string {
		m := r.matcher.FindStringSubmatch(t)[1:]

		a := make([]string, len(m))
		for i, from := range m {
			a[i] = r.replacement[i][from]
		}
		return strings.Join(a, "")
	})
}
