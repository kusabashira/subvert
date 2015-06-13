package main

import (
	"os"
)

func shortUsage() {
	os.Stderr.WriteString(`
Usage: subvert [OPTION]... FROM TO [FILE]...
Try 'subvert --help' for more information.
`[1:])
}

func usage() {
	os.Stderr.WriteString(`
Usage: subvert [OPTION]... FROM TO [FILE]...
Substitute multiple words at once.

Options:
  -h, --help                show this help message
`[1:])
}
