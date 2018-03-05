package main

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	groupRegexp             = regexp.MustCompile(`(?:[^/\\]|\\.)*`)
	branchRegexp            = regexp.MustCompile(`(?:[^,\\]|\\.)*`)
	escapedCharacterRegexp  = regexp.MustCompile(`\\(.)`)
	trailingBackslashRegexp = regexp.MustCompile(`\\+$`)
)

func parsePattern(pattern string) (tree [][]string, err error) {
	pattern = trailingBackslashRegexp.ReplaceAllStringFunc(pattern, func(s string) string {
		return strings.Repeat(`\\`, len(s)/2)
	})

	groups := groupRegexp.FindAllString(pattern, -1)
	for gi := 0; gi < len(groups); gi++ {
		branches := branchRegexp.FindAllString(groups[gi], -1)
		for bi := 0; bi < len(branches); bi++ {
			branches[bi] = escapedCharacterRegexp.ReplaceAllString(branches[bi], "$1")
		}
		tree = append(tree, branches)
	}
	return tree, nil
}

func toPatternRegexp(srcPattern string, useBoundary bool) (re *regexp.Regexp, err error) {
	srcTree, err := parsePattern(srcPattern)
	if err != nil {
		return nil, err
	}

	quotedGroups := make([]string, len(srcTree))
	for gi, branches := range srcTree {
		quotedBranches := make([]string, len(branches))
		for bi := 0; bi < len(branches); bi++ {
			quotedBranches[bi] = regexp.QuoteMeta(branches[bi])
		}
		quotedGroups[gi] = "(" + strings.Join(quotedBranches, "|") + ")"
	}

	rawRegexp := strings.Join(quotedGroups, "")
	if useBoundary {
		rawRegexp = `\b` + rawRegexp + `\b`
	}
	return regexp.Compile(rawRegexp)
}

func toReplaceTable(srcPattern, dstPattern string) (table []map[string]string, err error) {
	srcTree, err := parsePattern(srcPattern)
	if err != nil {
		return nil, err
	}
	dstTree, err := parsePattern(dstPattern)
	if err != nil {
		return nil, err
	}

	if len(srcTree) != len(dstTree) {
		return nil, fmt.Errorf("mismatch the number of groups")
	}

	table = make([]map[string]string, len(srcTree))
	for gi := 0; gi < len(srcTree); gi++ {
		if len(srcTree[gi]) != len(srcTree[gi]) {
			return nil, fmt.Errorf("mismatch the number of branches at group[%d]", gi)
		}

		table[gi] = make(map[string]string)
		for bi := 0; bi < len(srcTree[gi]); bi++ {
			src, dst := srcTree[gi][bi], dstTree[gi][bi]
			if _, exist := table[gi][src]; exist {
				return nil, fmt.Errorf("group[%d] has duplicated items", gi)
			}
			table[gi][src] = dst
		}
	}
	return table, nil
}

type Replacer struct {
	patternRegexp *regexp.Regexp
	replaceTable  []map[string]string
}

func NewReplacer(srcPattern string, dstPattern string, useBoundary bool) (r *Replacer, err error) {
	r = &Replacer{}

	r.patternRegexp, err = toPatternRegexp(srcPattern, useBoundary)
	if err != nil {
		return nil, err
	}
	r.replaceTable, err = toReplaceTable(srcPattern, dstPattern)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Replacer) Replace(s string) string {
	return r.patternRegexp.ReplaceAllStringFunc(s, func(t string) string {
		// Submatch 0 is the match of the entire expression.
		//
		// See: https://golang.org/pkg/regexp/
		//
		groups := r.patternRegexp.FindStringSubmatch(t)[1:]
		for gi := 0; gi < len(groups); gi++ {
			groups[gi] = r.replaceTable[gi][groups[gi]]
		}
		return strings.Join(groups, "")
	})
}
