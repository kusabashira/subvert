msub
====

Substitute multiple words at once
by FROM and TO patterns.

It's inspired by [tpope/vim-abolish](http://github.com/tpope/vim-abolish)

```
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
```

Usage
-----

```
$ msub [OPTION]... FROM TO [FILE]...
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
  msub true,false false,true ./file
  msub dog,cat/s cat,dog/s ~/Document/questionnaire
```

Installation
------------

### compiled binary

See [releases](https://github.com/nil-two/msub/releases)

### go get

```
go get github.com/nil-two/msub
```

Options
-------

### -h, --help

Display a help message.

### -v, --version

Output the version of msub.

### -b, --boundary

Replace only the string sandwiched word boundaries.

```
$ cat example
n, x = line[0], line[1]

$ cat example | msub n,x x,n
(matches /(n|x)/)
x, n = lixe[0], lixe[1]

$ cat example | msub --boundary n,x x,n
(matches /\b(n|x)\b/)
x, n = line[0], line[1]
```

Behavior
--------

Pattern can be separated by `,`.

Matched string will map to a string in the same index.

```
$ msub true,false false,true
true  -> false
false -> true

$ msub foo,bar,baz bar,baz,foo
foo -> bar
bar -> baz
baz -> foo
```

In addition, patterns can connect, separated by `/`.

Indexes are separately for each patterns.

```
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
```

### Inability

- loop patterns ( a,b -> a,b,a,b,a,b ... )
- three case variants ( box -> box, Box, BOX )

vim-abolish can both.

Syntax
------

Here is the syntax of pattern in extended BNF.

```
pattern = group , { "/" , group } ;
group   = branch , { "," , branch } ;
branch  = { [ "\" ] , ? unicode character ? - "/" - "," | "\/" | "\," } ;
```

`FROM` and `TO` are `pattern`.

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

nil2 <nil2@nil2.org>
