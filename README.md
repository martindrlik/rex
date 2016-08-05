# rex
I believe rex is fast and simple file search utility that uses regular expression to search in files.

## Options
Following command line arguments should help you specify your search criteria for fastest match.
* **dot** search dotfiles
* **exp** search regular expression exp in files
* **ner** print no error
* **pre** search pre prefixed file names
* **suf** search suf suffixed file names

## Installation
You might need to be familiar with [Go Programming Language](http://github.com/golang/go). At least I would advise to read [How to Write Go Code](https://golang.org/doc/code.html).
Finally you can use `go get` to install rex by following command:
```
$ go get github.com/martindrlik/rex
(no output means success)
```

## Example
Finally following shows how to use rex.
```
$ cd ~/src/go
$ rex -pre g -suf go -exp hello
doc/progs/go1.go
src/cmd/go/go_test.go
src/compress/gzip/gunzip_test.go
src/compress/gzip/gzip_test.go
src/encoding/gob/gobencdec_test.go
src/cmd/compile/internal/gc/global_test.go
```
