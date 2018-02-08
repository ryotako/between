package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
)

var (
	flagA bool
	flagH bool
	flagV bool
)

const (
	VERSION = "0.0"
)

type CLI struct {
	inStream             io.Reader
	outStream, errStream io.Writer
}

func main() {

	f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	f.BoolVar(&flagA, "all", false, "display all matched lines")
	f.BoolVar(&flagA, "a", false, "display all matched lines")
	f.BoolVar(&flagH, "help", false, "show help")
	f.BoolVar(&flagH, "h", false, "show help")
	f.BoolVar(&flagV, "version", false, "print the version")
	f.BoolVar(&flagV, "v", false, "print the version")

	f.Parse(os.Args[1:])

	c := &CLI{
		inStream:  os.Stdin,
		outStream: os.Stdout,
		errStream: os.Stderr,
	}

	os.Exit(c.Run(f.Args()))
}

func (c *CLI) Run(args []string) int {
	if flagH {
		c.Usage()
		return 0
	}

	if flagV {
		fmt.Fprintf(c.outStream, "between version %s", VERSION)
		return 0
	}

	ss := []io.Reader{}

	switch len(args) {
	case 0, 1:
		fmt.Fprintln(c.errStream, "between: too few arguments")
		return 1
	case 2:
		ss = append(ss, c.inStream)
	default:
		for _, a := range args[2:] {
			f, err := os.Open(a)
			if err != nil {
				fmt.Fprintf(c.errStream, "between: cannot open `%s`\n", a)
				continue
			}
			defer f.Close()

			ss = append(ss, f)
		}
	}

	r1, err := regexp.Compile(args[0])
	if err != nil {
		fmt.Fprintf(c.errStream, "between: invalid regexp `%s`\n", args[0])
		return 1
	}
	r2, err := regexp.Compile(args[1])
	if err != nil {
		fmt.Fprintf(c.errStream, "between: invalid regexp `%s`\n", args[1])
		return 1
	}

	for _, s := range ss {
		s := bufio.NewScanner(s)
		ok := false
		cnt := 0
		for s.Scan() {
			text := s.Text()

			if !ok && r1.MatchString(text) {
				ok = true
				cnt++
			}

			if ok {
				if cnt == 1 || flagA {
					fmt.Fprintln(c.outStream, text)
				}
			}

			if ok && r2.MatchString(text) {
				ok = false
			}
		}

	}

	return 0
}

func (c *CLI) Usage() {
	fmt.Fprintf(c.outStream, `NAME:
   between - display from the matched line to the matched line

USAGE:
   between [options] cond1 cond2 [files...]

VERSION:
   %s

OPTIONS:
   -a, --all      display all matched lines
   -h, --help     show help
   -v, --version  print the version
`, VERSION)
}
