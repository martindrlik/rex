// rex is regular expression using file search utility.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
)

var (
	dot = flag.Bool("dot", false, "search dotfiles")
	ner = flag.Bool("ner", false, "print no error")

	exp = flag.String("exp", "", "search regular expression exp in files")
	pre = flag.String("pre", "", "search pre prefixed file names")
	suf = flag.String("suf", "", "search suf suffixed file names")
)

var (
	rex *regexp.Regexp

	me sync.Mutex
	wg sync.WaitGroup
)

func main() {
	flag.Parse()
	var err error
	if *exp != "" {
		rex, err = regexp.Compile(*exp)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "rex: -exp: %v\n", err)
		os.Exit(1)
	}
	readDir(".")
	wg.Wait()
}

func readDir(name string) {
	dir, err := os.Open(name)
	if err != nil {
		pe(err)
		return
	}
	defer dir.Close()
	fi, err := dir.Readdir(0)
	if err != nil {
		pe(err)
		return
	}
	wg.Add(len(fi))
	for _, fi := range fi {
		go func(fi os.FileInfo) {
			readFile(name, fi)
			wg.Done()
		}(fi)
	}
}

func readFile(dir string, fi os.FileInfo) {
	name := fi.Name()
	if !*dot && name[0] == '.' {
		return
	}
	if fi.IsDir() {
		readDir(path.Join(dir, name))
		return
	}
	if *pre != "" && !strings.HasPrefix(name, *pre) {
		return
	}
	if *suf != "" && !strings.HasSuffix(name, *suf) {
		return
	}
	ful := path.Join(dir, name)
	if *exp == "" {
		po(ful)
		return
	}
	f, err := os.Open(ful)
	if err != nil {
		pe(err)
		return
	}
	r := bufio.NewReader(f)
	m := rex.MatchReader(r)
	f.Close()
	if !m {
		return
	}
	po(ful)
}

func pe(err error) {
	if *ner {
		return
	}
	me.Lock()
	fmt.Fprintf(os.Stderr, "rex: %v\n", err)
	me.Unlock()
}

func po(ln string) {
	me.Lock()
	fmt.Println(ln)
	me.Unlock()
}
