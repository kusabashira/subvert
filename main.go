package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

const (
	cmdName    = "msub"
	cmdVersion = "0.3.0"
)

type CLI struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer

	useBoundary bool
	isHelp      bool
	isVersion   bool
}

func NewCLI(stdin io.Reader, stdout io.Writer, stderr io.Writer) *CLI {
	return &CLI{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}
}

func (c *CLI) parseOptions(args []string) (leftArgs []string, err error) {
	f := flag.NewFlagSet(cmdName, flag.ContinueOnError)
	f.SetOutput(ioutil.Discard)

	f.BoolVar(&c.useBoundary, "b", false, "")
	f.BoolVar(&c.useBoundary, "boundary", false, "")
	f.BoolVar(&c.isHelp, "h", false, "")
	f.BoolVar(&c.isHelp, "help", false, "")
	f.BoolVar(&c.isVersion, "v", false, "")
	f.BoolVar(&c.isVersion, "version", false, "")

	if err = f.Parse(args); err != nil {
		return nil, err
	}
	return f.Args(), nil
}

func (c *CLI) printUsage() {
	fmt.Fprintf(c.stderr, `
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

func (c *CLI) printVersion() {
	fmt.Fprintf(c.stderr, "%s\n", cmdVersion)
}

func (c *CLI) printErr(err interface{}) {
	fmt.Fprintf(c.stderr, "%s: %s\n", cmdName, err)
}

func (c *CLI) do(rep *Replacer, r io.Reader) error {
	bs := bufio.NewScanner(r)
	for bs.Scan() {
		fmt.Fprintln(c.stdout, rep.Replace(bs.Text()))
	}
	return bs.Err()
}

func (c *CLI) Run(args []string) int {
	leftArgs, err := c.parseOptions(args)
	if err != nil {
		c.printErr(err)
		return 2
	}

	if c.isHelp {
		c.printUsage()
		return 0
	}
	if c.isVersion {
		c.printVersion()
		return 0
	}

	if len(leftArgs) < 1 {
		c.printErr("no specify FROM and TO")
		return 2
	}
	if len(leftArgs) < 2 {
		c.printErr("no specify TO")
		return 2
	}

	srcPattern := leftArgs[0]
	dstPattern := leftArgs[1]
	filePathes := leftArgs[2:]

	rep, err := NewReplacer(srcPattern, dstPattern, c.useBoundary)
	if err != nil {
		c.printErr(err)
		return 2
	}

	exitCode := 0
	if len(filePathes) == 0 {
		if err = c.do(rep, c.stdin); err != nil {
			c.printErr(err)
			exitCode = 1
		}
	} else {
		for _, filePath := range filePathes {
			f, err := os.Open(filePath)
			if err != nil {
				c.printErr(err)
				exitCode = 1
				continue
			}
			defer f.Close()

			if err = c.do(rep, f); err != nil {
				c.printErr(err)
				exitCode = 1
			}
		}
	}
	return exitCode
}

func main() {
	c := NewCLI(os.Stdin, os.Stdout, os.Stderr)
	e := c.Run(os.Args[1:])
	os.Exit(e)
}
