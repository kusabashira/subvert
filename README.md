msub
====

[![Build Status](https://travis-ci.org/kusabashira/msub.svg?branch=master)](https://travis-ci.org/kusabashira/msub)

Substitute multiple words at once
by FROM and TO patterns.

It's inspired by [tpope/vim-abolish](http://github.com/tpope/vim-abolish)

	$ cat questionnaire
	1 true
	2 true
	3 false
	4 false
	5 true

	$ cat questionnaire | msub true,false false,true
	1 false
	2 false
	3 true
	4 true
	5 false

Usage
-----

	$ msub [OPTION]... FROM TO [FILE]...
	Substitute multiple words at once
	by FROM and TO patterns.

	Options:
	  -b, --boundary    use word boundary in matcher
	  -h, --help        show this help message

Installation
------------

###go get

	go get github.com/kusabashira/msub

Behavior
--------

Pattern can be separated by a `,`.

Matched string will map to a string in the same index.

	$ msub true,false false,true
	true  -> false
	false -> true

	$ msub foo,bar,baz bar,baz,foo
	foo -> bar
	bar -> baz
	baz -> foo

In addition, patterns can connect, separated by a `/`.

Indexes are separately for each patterns.

	$ msub cat,dog/,s dog,cat/,s
	cat  -> dog
	cats -> dogs
	dog  -> cat
	dogs -> cats

	$ msub 'V,v/im/ ,/s,S/cript' 'V,V/im/ , /s,s/cript'
	Vim script -> Vim script
	Vim Script -> Vim script
	Vimscript  -> Vim script
	VimScript  -> Vim script
	vim script -> Vim script
	vim Script -> Vim script
	vimscript  -> Vim script
	vimScript  -> Vim script

###Inability

- loop patterns ( a,b -> a,b,a,b,a,b ... )
- three case variants (box -> box, Box, BOX)

vim-abolish can both.

Syntax
------

Here is the syntax of msub in extended BNF. 

	pattern = group {"/" group}
	group   = branch {"," branch}
	branch  = {letter | "\/" | "\,"}

- FROM and TO are `pattern`.
- letter is a unicode character (ignore `/` and `,`).

Correspondence of vim-abolish is as follows:

| msub                 | vim-abolish        |
|----------------------|--------------------|
| foo                  | foo                |
| true,false           | {true,false}       |
| dog,cat/s            | {dog,cat}s         |
| ,f,s/print/,f,ln     | {,f,s}print{,f,ln} |
| a,b,a,b/a,a,a,a      | {a,b}{a}           |
| V,V/im/ , /s,s/cript | {V}im{ }{s}cript   |

License
-------

MIT License

Author
------

kusabashira <kusabashira227@gmail.com>
