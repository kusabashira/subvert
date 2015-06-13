package main

import (
	"os"
)

func usage() {
	os.Stderr.WriteString(`
Usage: subvert [OPTION]... FROM TO [FILE]...
Substitute multiple words at once.

Options:
  -h, --help                show this help message
`[1:])
}
