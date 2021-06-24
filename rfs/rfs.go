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
)

var (
	ErrEmptyCmd = errors.New("empty cmd")
	ErrTimeout  = errors.New("request timeout")
)

type FileStat struct {
	Name    string
	Size    string
	Mod     uint32
	ModTime int64
	IsDir   string
}

type Cmd struct {
	Name string
	Args []string
	Env  []string

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

type Rfs interface {
	Init() error
	List(ctx context.Context) ([]*FileStat, error)
	Write(ctx context.Context, toPath string, src io.Reader) error
	Copy(ctx context.Context, fromPath, toPath string) error
	Exec(ctx context.Context, cmd *Cmd) error
}
