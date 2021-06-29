// MIT License
//
// Copyright (c) 2021 Lack
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package rfs

import (
	"context"
	"errors"
	"io"
	"os"
	"sync"
)

var (
	ErrConnect       = errors.New("connect error")
	ErrEmptyCmd      = errors.New("empty cmd")
	ErrInvalidWrite  = errors.New("invalid write result")
	ErrTimeout       = errors.New("request timeout")
	ErrMissingCmd    = errors.New("missing stdout or stdout")
	ErrRequest       = errors.New("request exception")
	ErrNotExists     = errors.New("file does not exist")
	ErrAlreadyExists = errors.New("file already exists")
)

type Cmd struct {
	Name string
	Args []string
	Env  []string

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

type singleWriter struct {
	io.Writer
	mu sync.Mutex
}

func (w *singleWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.Writer.Write(p)
}

func NewCmd(name string, args []string, writer io.Writer) *Cmd {
	w := singleWriter{Writer: writer}
	return &Cmd{
		Name:   name,
		Args:   args,
		Env:    []string{},
		Stdin:  os.Stdin,
		Stdout: &w,
		Stderr: &w,
	}
}

type IOMetric struct {
	Name  string
	From  string
	To    string
	Total int64
	Block int64
	Speed int64
}

type IOFn func(*IOMetric)

type Rfs interface {
	Init() error
	Exec(ctx context.Context, cmd *Cmd, opts ...ExecOption) error
	List(ctx context.Context, remotePath string, opts ...ListOption) ([]os.FileInfo, error)
	Get(ctx context.Context, remotePath, localPath string, fn IOFn, opts ...GetOption) error
	Put(ctx context.Context, localPath, remotePath string, fn IOFn, opts ...PutOption) error
}
