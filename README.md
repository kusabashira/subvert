msub
====

Substitute multiple words at once
by FROM and TO patterns.

It's inspired by [tpope/vim-abolish](http://github.com/tpope/vim-abolish)

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
	  -b, --boundary    use word boundary in matcher
	  -h, --help        show this help message

Installation
------------

###go get

	go get github.com/kusabashira/msub

Behavior
--------

Pattern can be separated by a ",".

Matched string will map to a string in the same index.

	$ msub true,false false,true
	true  -> false
	false -> true

	$ msub foo,bar,baz bar,baz,foo
	foo -> bar
	bar -> baz
	baz -> foo

In addition, pattern can connect, separated by a "/".

Indexes are separately for each pattern.

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

###Inability now

- loop pattern ( a,b -> a,b,a,b,a,b ... )
- three case variants (box -> box, Box, BOX)

vim-abolish can both.

Syntax
------

Here is the syntax of msub in extended BNF. 

	pattern = group {"/" group}
	group   = branch {"," branch}
	branch  = {letter | "\/" | "\,"}

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

If the replacement should be identical to the pattern

License
-------

MIT License

Author
------

wara <kusabashira227@gmail.com>
