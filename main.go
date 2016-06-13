package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

var (
	cmdName = "msub"
	version = "0.3.0"

	flagset     = flag.NewFlagSet(cmdName, flag.ContinueOnError)
	useBoundary = flagset.Bool("boundary", false, "")
	isHelp      = flagset.Bool("help", false, "")
	isVersion   = flagset.Bool("version", false, "")
)

func init() {
	flagset.SetOutput(ioutil.Discard)
	flagset.BoolVar(useBoundary, "b", false, "")
	flagset.BoolVar(isHelp, "h", false, "")
	flagset.BoolVar(isVersion, "v", false, "")
}

func usage() {
	fmt.Fprintf(os.Stderr, `
Usage: %[1]s [OPTION]... FROM TO [FILE]...
Substitute multiple words at once
by FROM and TO patterns.

Options:
  -b, --boundary    use word boundary in matcher
  -h, --help        show this help message and exit
  -v, --version     output version information and exit

Syntax:
  pattern = group , { "/" , group } ;
  group   = branch , { "," , branch } ;
  branch  = { [ "\" ] , ? unicode character ? - "/" - "," | "\/" | "\," } ;

Examples:
  %[1]s true,false false,true ./file
  %[1]s dog,cat/s cat,dog/s ~/Document/questionnaire
`[1:], cmdName)
}

func printVersion() {
	fmt.Fprintln(os.Stderr, version)
}

func printErr(err interface{}) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", cmdName, err)
}

func do(rep *Replacer, r io.Reader) error {
	b := bufio.NewScanner(r)
	for b.Scan() {
		fmt.Println(rep.Replace(b.Text()))
	}
	return b.Err()
}

func _main() int {
	if err := flagset.Parse(os.Args[1:]); err != nil {
		printErr(err)
		return 2
	}
	if *isHelp {
		usage()
		return 0
	}
	if *isVersion {
		printVersion()
		return 0
	}

	if flagset.NArg() < 1 {
		printErr("no specify FROM and TO")
		return 2
	}
	if flagset.NArg() < 2 {
		printErr("no specify TO")
		return 2
	}
	from, to := flagset.Arg(0), flagset.Arg(1)

	rep, err := NewReplacer(from, to, *useBoundary)
	if err != nil {
		printErr(err)
		return 2
	}

	var r io.Reader
	if flagset.NArg() < 3 {
		r = os.Stdin
	} else {
		var a []io.Reader
		for _, file := range flagset.Args()[2:] {
			f, err := os.Open(file)
			if err != nil {
				printErr(err)
				return 1
			}
			defer f.Close()
			a = append(a, f)
		}
		r = io.MultiReader(a...)
	}

	if err = do(rep, r); err != nil {
		printErr(err)
		return 1
	}
	return 0
}

func main() {
	e := _main()
	os.Exit(e)
}
