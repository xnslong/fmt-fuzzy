package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

var sepFlag string
var ignoreMismatch string
var splitter = "[](){},;"

func main() {
	Format(os.Stdin, os.Stdout)
}

func Format(input io.Reader, stdout io.Writer) {
	scanner := bufio.NewScanner(input)
	scanner.Split(Splitter)
	s := &stack{}
	o := &output{
		cleanLine: true,
		writer:    stdout,
		indentStr: strings.Repeat(" ", 4),
		indent:    s.Depth,
	}

	for scanner.Scan() {
		data := scanner.Text()
		data = strings.TrimSpace(data)
		if len(data) == 0 {
			continue
		}

		sp, spType := IsSplitter(data)

		switch spType {
		case SplitterTypeOpen:
			o.Print(data)
			o.NewLine()
			s.Push(sp)
		case SplitterTypeSplit:
			o.Print(data)
			o.NewLine()
		case SplitterTypeNone:
			o.Print(data)
		case SplitterTypeClose:
			Close(s, sp)
			if !o.IsClean() {
				o.NewLine()
			}
			o.Print(data)
		}
	}
}

func Close(s *stack, closeSp byte) bool {
	for {
		pop, ok := s.Pop()
		if !ok {
			return false
		}
		if Closes(closeSp, pop) {
			return true
		}
	}
}

var quoteMap = map[byte]byte{
	'[': ']',
	'{': '}',
	'(': ')',
}

func Closes(closeSp byte, openSp byte) bool {
	return quoteMap[openSp] == closeSp
}

type output struct {
	cleanLine bool
	indentStr string
	indent    func() int
	writer    io.Writer
}

func (o *output) Print(v string) {
	if o.cleanLine {
		i := o.indent()
		fmt.Fprint(o.writer, strings.Repeat(o.indentStr, i))
		o.cleanLine = false
	}

	fmt.Fprint(o.writer, v)
}

func (o *output) NewLine() {
	fmt.Fprintln(o.writer)
	o.cleanLine = true
}

func (o *output) IsClean() bool {
	return o.cleanLine
}

type SplitterType int

const (
	SplitterTypeOpen = SplitterType(iota)
	SplitterTypeClose
	SplitterTypeSplit
	SplitterTypeNone
)

func IsSplitter(data string) (q byte, ok SplitterType) {
	if len(data) > 1 {
		return 0, SplitterTypeNone
	}

	v := data[0]
	switch v {
	case '{', '[', '(':
		return v, SplitterTypeOpen
	case '}', ']', ')':
		return v, SplitterTypeClose
	case ',', ';':
		return v, SplitterTypeSplit
	default:
		return 0, SplitterTypeNone
	}

}

type stack struct {
	elements []byte
}

func (s *stack) Push(v byte) {
	s.elements = append(s.elements, v)
}

func (s *stack) Pop() (byte, bool) {
	l := s.Depth()
	if l == 0 {
		return 0, false
	}

	e := s.elements[l-1]
	s.elements = s.elements[:l-1]
	return e, true
}

func (s *stack) Peek() (byte, bool) {
	l := s.Depth()
	if l == 0 {
		return 0, false
	}

	return s.elements[l-1], true
}

func (s *stack) IsEmpty() bool {
	return s.Depth() == 0
}

func (s *stack) Depth() int {
	return len(s.elements)
}

func Splitter(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	index := bytes.IndexAny(data, splitter)
	if index == 0 {
		return 1, data[:1], nil
	}
	if index > 0 {
		return index, data[:index], nil
	}
	if atEOF {
		return len(data), data, nil
	}

	// require more data
	return 0, nil, nil
}
