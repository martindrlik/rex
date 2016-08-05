package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
)

var (
	dot = flag.Bool("dot", false, "search dotfiles")
	exp = flag.String("exp", "", "")
	pre = flag.String("pre", "", "search pre prefixed file names")
	suf = flag.String("suf", "", "search suf suffixed file names")
)

var (
	rex *regexp.Regexp
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
	check(readDir("."))
}

func readDir(dirName string) error {
	dir, err := os.Open(dirName)
	if err != nil {
		return err
	}
	defer dir.Close()
	fi, err := dir.Readdir(0)
	if err != nil {
		return err
	}
	for _, fi := range fi {
		name := fi.Name()
		if !*dot && name[0] == '.' {
			continue
		}
		if fi.IsDir() {
			full := path.Join(dirName, name)
			check(readDir(full))
			continue
		}
		if *pre != "" && !strings.HasPrefix(name, *pre) {
			continue
		}
		if *suf != "" && !strings.HasSuffix(name, *suf) {
			continue
		}
		full := path.Join(dirName, name)
		if *exp == "" {
			fmt.Println(full)
			continue
		}
		m, err := match(full)
		if check(err) {
			continue
		}
		if m {
			fmt.Println(full)
		}
	}
	return nil
}

func match(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	return rex.MatchReader(r), nil
}

func check(err error) bool {
	if err == nil {
		return false
	}
	fmt.Fprintf(os.Stderr, "rex: %v\n", err)
	return true
}
