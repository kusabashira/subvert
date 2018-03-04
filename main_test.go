package main

import (
	"bytes"
	"regexp"
	"strings"
	"testing"
)

func TestHelp(t *testing.T) {
	stdin := bytes.NewReader(nil)
	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)
	args := []string{"--help"}

	c := NewCLI(stdin, stdout, stderr)

	exitCode := c.Run(args)
	if exitCode != 0 {
		t.Fatalf("Run(%q) should return 0, but got %d", args, exitCode)
	}

	output := stderr.String()
	if !strings.HasPrefix(output, "Usage") {
		t.Fatalf("Run(%q) should output usage to stderr", args)
	}
}

func TestVersion(t *testing.T) {
	stdin := bytes.NewReader(nil)
	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)
	args := []string{"--version"}

	c := NewCLI(stdin, stdout, stderr)

	exitCode := c.Run(args)
	if exitCode != 0 {
		t.Fatalf("Run(%q) should return 0, but got %d", args, exitCode)
	}

	output := stderr.String()
	if !regexp.MustCompile(`^\d+.\d+.\d`).MatchString(output) {
		t.Fatalf("Run(%q) should output version to stderr", args)
	}
}

type runningTest struct {
	Description string
	Args        []string
	Src         string
	Dst         string
}

var runningTests = []runningTest{
	{
		Description: "AAA -> BBB",
		Args:        []string{"AAA", "BBB"},
		Src: `
AAA AAA
BBB CCC
`[1:],
		Dst: `
BBB BBB
BBB CCC
`[1:],
	},
	{
		Description: "AAA -> BBB, BBB -> AAA",
		Args:        []string{"AAA,BBB", "BBB,AAA"},
		Src: `
AAA AAA
BBB CCC
`[1:],
		Dst: `
BBB BBB
AAA CCC
`[1:],
	},
	{
		Description: "AACC -> DDFF, BBCC -> EEFF",
		Args:        []string{"AA,BB/CC", "DD,EE/FF"},
		Src: `
AACC BBCC AABB
DDFF EEFF
`[1:],
		Dst: `
DDFF EEFF AABB
DDFF EEFF
`[1:],
	},
	{
		Description: "Fix Vim script",
		Args:        []string{"V,v/im/ ,/s,S/cript", "V,V/im/ , /s,s/cript"},
		Src: `
Vim script Vim Script Vimscript VimScript
vim script vim Script vimscript vimScript
`[1:],
		Dst: `
Vim script Vim script Vim script Vim script
Vim script Vim script Vim script Vim script
`[1:],
	},
}

func TestRun(t *testing.T) {
	for _, test := range runningTests {
		stdin := bytes.NewReader([]byte(test.Src))
		stdout := bytes.NewBuffer([]byte{})
		stderr := bytes.NewBuffer([]byte{})

		c := NewCLI(stdin, stdout, stderr)

		exitCode := c.Run(test.Args)
		if exitCode != 0 {
			t.Errorf("Run(%q) should return 0, but got %d", test.Args, exitCode)
			continue
		}

		expectOutput := test.Dst
		actualOutput := stdout.String()
		if actualOutput != expectOutput {
			t.Errorf("Run(%q): %s\ninput:\n%s\nactual:\n%s\nexpect:\n%s",
				test.Args, test.Description, test.Src, actualOutput, expectOutput)
		}
	}
}
