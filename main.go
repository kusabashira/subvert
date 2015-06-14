package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func shortUsage() {
	os.Stderr.WriteString(`
Usage: msub [OPTION]... FROM TO [FILE]...
Try 'msub --help' for more information.
`[1:])
}

func usage() {
	os.Stderr.WriteString(`
Usage: msub [OPTION]... FROM TO [FILE]...
Substitute multiple words at once
by FROM and TO patterns.

Options:
  -b, --boundary    use word boundary in matcher
  -h, --help        show this help message

Syntax:
  pattern = group {"/" group}
  group   = branch {"," branch}
  branch  = {letter | "\/" | "\,"}

Examples:
  msub true,false false,true ./file
  msub dog,cat/s cat,dog/s ~/Document/questionnaire
`[1:])
}

func printError(err error) {
	fmt.Fprintln(os.Stderr, "msub:", err)
}

func do(r *Replacer, src io.Reader) error {
	b := bufio.NewScanner(src)
	for b.Scan() {
		fmt.Println(r.ReplaceAll(b.Text()))
	}
	return b.Err()
}

func _main() int {
	var useBoundary bool
	flag.BoolVar(&useBoundary, "b", false, "")
	flag.BoolVar(&useBoundary, "boundary", false, "")

	var isHelp bool
	flag.BoolVar(&isHelp, "h", false, "")
	flag.BoolVar(&isHelp, "help", false, "")
	flag.Usage = shortUsage
	flag.Parse()
	switch {
	case isHelp:
		usage()
		return 0
	case flag.NArg() < 1:
		printError(fmt.Errorf("no specify FROM and TO"))
		return 2
	case flag.NArg() < 2:
		printError(fmt.Errorf("no specify TO"))
		return 2
	}
	from, to := flag.Arg(0), flag.Arg(1)

	r, err := NewReplacer(from, to, useBoundary)
	if err != nil {
		printError(err)
		return 2
	}

	if flag.NArg() < 3 {
		if err = do(r, os.Stdin); err != nil {
			printError(err)
			return 1
		}
		return 0
	}

	var srcls []io.Reader
	for _, name := range flag.Args()[2:] {
		f, err := os.Open(name)
		if err != nil {
			printError(err)
			return 1
		}
		defer f.Close()
		srcls = append(srcls, f)
	}
	if err = do(r, io.MultiReader(srcls...)); err != nil {
		printError(err)
		return 1
	}
	return 0
}

func main() {
	e := _main()
	os.Exit(e)
}
