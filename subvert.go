package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	groups   = regexp.MustCompile(`(?:[^/\\]|\\.)*`)
	branches = regexp.MustCompile(`(?:[^,\\]|\\.)*`)
)

func parseExpr(expr string) ([][]string, error) {
	expr = strings.Replace(expr, `\,`, `\\,`, -1)
	expr = strings.Replace(expr, `\/`, `\\/`, -1)
	expr, err := strconv.Unquote(`"` + expr + `"`)
	if err != nil {
		return nil, err
	}

	gls := groups.FindAllString(expr, -1)
	tree := make([][]string, len(sls))
	for gi := 0; gi < len(gls); gi++ {
		bls := branches.FindAllString(gls[gi], -1)
		for bi := 0; bi < len(bls); bi++ {
			bls[bi] = strings.Replace(bls[bi], `\,`, `,`, -1)
			bls[bi] = strings.Replace(bls[bi], `\/`, `/`, -1)
			bls[bi] = regexp.QuoteMeta(bls[bi])
		}
		tree[gi] = bls
	}
	return tree, nil
}

func newMatcher(expr string) (m *regexp.Regexp, err error) {
	tree, err := parseExpr(expr)
	if err != nil {
		return nil, err
	}

	sls := make([]string, len(tree))
	for gi, bls := range tree {
		sls[gi] = "(" + strings.Join(bls, "|") + ")"
	}
	return regexp.Compile(strings.Join(sls, ""))
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
			return nil, fmt.Errorf("mismatch the number of group[%q]", gi)
		}

		r[gi] = make(map[string]string)
		for bi := 0; bi < len(from[gi]); bi++ {
			src, dst := from[gi][bi], to[gi][bi]
			if _, exist := r[gi][src]; exist {
				return nil, fmt.Errorf("group[%q] has duplicate items", gi)
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

func NewReplacer(from, to string) (r *Replacer, err error) {
	r = &Replacer{}

	r.matcher, err = newMatcher(from)
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
