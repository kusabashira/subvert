msub
====

Substitute multiple words at once
by FROM and TO patterns.

Its based on [tpope/vim-abolish](http://github.com/tpope/vim-abolish)

	$ cat questionnaire
	1 cat
	2 cat
	3 dog
	4 dog
	5 cat

	$ cat questionnaire | msub cat,dog dog,cat
	1 dog
	2 dog
	3 cat
	4 cat
	5 dog

Usage
-----

	$ msub [OPTION]... FROM TO [FILE]...
	Substitute multiple words at once
	by FROM and TO patterns.

	Options:
	  -h, --help        show this help message

Installation
------------

###go get

	go get github.com/kusabashira/msub

Syntax
------

Here is the syntax of msub in extended BNF. 

	pattern = group {"/" group}
	group   = branch {"," branch}
	branch  = letter | "\/" | "\,"

- FROM and TO is a `pattern`.
- letter is a unicode character (ignore "/" and ",")

In this way, It is unlike vim-abolish.
Correspondence is as follows.

| msub             | vim-abolish        |
|------------------|--------------------|
| foo              | foo                |
| true,false       | {true,false}       |
| dog,cat/s        | {dog,cat}s         |
| ,f,s/print/,f,ln | {,f,s}print{,f,ln} |

License
-------

MIT License

Author
------

wara <kusabashira227@gmail.com>
