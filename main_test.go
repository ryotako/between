package main

import (
	"bytes"
	"strconv"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}

	tests := []struct {
		stdin, input, want string
		err                bool
	}{
		{stdin: seq(5), input: "2 4", want: "2\n3\n4\n", err: false},
	}

	for _, test := range tests {
		cli.inStream = bytes.NewBufferString(test.stdin)
		outStream.Reset()
		errStream.Reset()

		status := cli.Run(strings.Fields(test.input))

		if test.err {
			if status == 0 {
				t.Errorf("%q >> status code should be non-zero", test.input)
			}
			if len(errStream.String()) == 0 {
				t.Errorf("%q >> status code %d, but no error messages", test.input, status)
			}
		} else {
			if status != 0 {
				t.Errorf("%q >> status code %d should be zero", test.input, status)
			}
			got := outStream.String()
			if test.want != got {
				t.Errorf("%q >> want %q, but %q", test.input, test.want, got)
			}
		}
	}
}

func seq(n int) string {
	s := ""
	for i := 1; i <= n; i++ {
		s += strconv.Itoa(i) + "\n"
	}
	return s
}
