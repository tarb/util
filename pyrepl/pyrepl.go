package pyrepl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"sync"
)

//
type PythonREPL struct {
	cmd    *exec.Cmd
	stdout *bufio.Reader
	stdin  io.WriteCloser

	callQueue chan call

	closer sync.Once
}

type call struct {
	in  string
	out chan []byte
}

//
func NewPythonREPL(script string) (*PythonREPL, error) {
	var err error
	var pr = &PythonREPL{}

	pr.cmd = exec.Command("python", "-i", script)

	var out io.ReadCloser
	if out, err = pr.cmd.StdoutPipe(); err != nil {
		return nil, err
	}

	pr.stdout = bufio.NewReader(out)

	if pr.stdin, err = pr.cmd.StdinPipe(); err != nil {
		return nil, err
	}

	pr.callQueue = make(chan call, 50)
	return pr, nil
}

//
func (repl *PythonREPL) Call(code string) <-chan []byte {
	var c = call{
		in:  code,
		out: make(chan []byte, 1),
	}

	repl.callQueue <- c

	return c.out
}

//
func (repl *PythonREPL) Close() {
	repl.closer.Do(func() {
		repl.stdin.Close()
		close(repl.callQueue)
	})
}

//
func (repl *PythonREPL) Start() error {
	go repl.listenCalls()
	return repl.cmd.Start()
}

//
func (repl *PythonREPL) listenCalls() {
	for c := range repl.callQueue {
		var err error

		if _, err = fmt.Fprintln(repl.stdin, c.in, `; print('\0')`); err != nil {
			fmt.Fprintln(repl.stdin, `print('\0')`)
			continue
		}

		var out []byte
		if out, err = repl.stdout.ReadBytes('\x00'); err != nil {
			repl.Close()
			continue
		}

		c.out <- bytes.TrimSpace(out[:len(out)-1])
		close(c.out)
	}
}
