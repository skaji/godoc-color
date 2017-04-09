package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"

	"golang.org/x/crypto/ssh/terminal"
)

var (
	regexpCode               = regexp.MustCompile(`^\t`)
	regexpSection            = regexp.MustCompile(`^[A-Z ]+$`)
	regexpSectionCode        = regexp.MustCompile(`^(?:const|func|var|type) .*$`)
	regexpSectionCodeComment = regexp.MustCompile(`^\s*//`)
)

var (
	colorCode               = "36"
	colorSection            = "4;1"
	colorSectionCode        = "32"
	colorSectionCodeComment = "1;30"
)

var (
	escape1 = regexp.MustCompile(`^\}`)
	escape2 = regexp.MustCompile(`^\)`)
)

type Trans struct {
	Out    io.Writer
	escape *regexp.Regexp
}

func (t *Trans) write(color, line string) {
	fmt.Fprintf(t.Out, "\x1b[%sm%s\x1b[m\n", color, line)
}

func (t *Trans) Render(line string) {
	if t.escape != nil {
		color := colorSectionCode
		if regexpSectionCodeComment.MatchString(line) {
			color = colorSectionCodeComment
		}
		t.write(color, line)
		if t.escape.MatchString(line) {
			t.escape = nil
		}
	} else if regexpCode.MatchString(line) {
		t.write(colorCode, line)
	} else if regexpSection.MatchString(line) {
		t.write(colorSection, line)
	} else if regexpSectionCode.MatchString(line) {
		if last := line[len(line)-1:]; last == "{" {
			t.escape = escape1
		} else if last == "(" {
			t.escape = escape2
		}
		t.write(colorSectionCode, line)
	} else {
		fmt.Fprintln(t.Out, line)
	}
}

func main() {
	if terminal.IsTerminal(0) {
		fmt.Fprintf(os.Stderr, "Usage: godoc ARGS | %s | less -R\n", os.Args[0])
		os.Exit(1)
	}
	trans := &Trans{Out: os.Stdout}
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		trans.Render(line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
