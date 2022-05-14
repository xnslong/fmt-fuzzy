package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type A struct {
	S  string
	I  int
	B  bool
	F  float32
	A  *A
	BV B
}

type B struct {
	S string
	I int
	B bool
	F float32
}

func Make() *A {
	return &A{
		S: "hello world",
		I: 1,
		B: true,
		F: 0.1,
		A: &A{},
		BV: B{
			S: "foo bar",
			I: 2,
			B: false,
			F: 0.2,
		},
	}
}

func TestGo(t *testing.T) {
	v := Make()
	fmt.Printf("%+v", v)
}

func TestFormat(t *testing.T) {
	type Case struct {
		Name   string
		Input  string
		Output string
	}

	cases := []Case{
		{
			Name:  "fuzzy",
			Input: `10.186.36.13(somebody@hostname:dir):#{"a":0,"b":{"c":"1013","d":"1049","e":"1453",},"f":"success"}`,
			Output: `10.186.36.13(
    somebody@hostname:dir
):#{
    "a":0,
    "b":{
        "c":"1013",
        "d":"1049",
        "e":"1453",
    },
    "f":"success"
}`,
		},
		{
			Name: "formatted already",
			Input: `10.186.36.13(
    somebody@hostname:dir
):#{
    "a":0,
    "b":{
        "c":"1013",
        "d":"1049",
        "e":"1453",
    },
    "f":"success"
}`,
			Output: `10.186.36.13(
    somebody@hostname:dir
):#{
    "a":0,
    "b":{
        "c":"1013",
        "d":"1049",
        "e":"1453",
    },
    "f":"success"
}`,
		},
		{
			Name:  "json",
			Input: `{"a":0,"b":{"c":"1013","d":"1049","e":"1453"},"f":"success"}`,
			Output: `{
    "a":0,
    "b":{
        "c":"1013",
        "d":"1049",
        "e":"1453"
    },
    "f":"success"
}`,
		},
		{
			Name:  "gostruct",
			Input: "&{S:hello world, I:1, B:true, F:0.1, A:0xc0000b4500, BV:{S:foo bar, I:2, B:false, F:0.2,}}",
			Output: `&{
    S:hello world,
    I:1,
    B:true,
    F:0.1,
    A:0xc0000b4500,
    BV:{
        S:foo bar,
        I:2,
        B:false,
        F:0.2,
    }
}`,
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {

			rdr := strings.NewReader(c.Input)
			wtr := &bytes.Buffer{}

			Format(rdr, wtr)

			assert.Equal(t, c.Output, wtr.String())
		})
	}
}
