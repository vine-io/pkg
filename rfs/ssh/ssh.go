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

package ssh

import (
	"context"
	"io"
	"strings"

	"github.com/lack-io/pkg/rfs"
	"golang.org/x/crypto/ssh"
)

var _ rfs.Rfs = (*client)(nil)

type client struct {
	options Options

	cc *ssh.Client
}

func (c *client) Init() error {

	var err error
	c.cc, err = ssh.Dial(c.options.network, c.options.addr, &c.options.ClientConfig)
	return err
}

func (c *client) List(ctx context.Context) ([]*rfs.FileStat, error) {
	panic("implement me")
}

func (c *client) Write(ctx context.Context, toPath string, src io.Reader) error {
	panic("implement me")
}

func (c *client) Copy(ctx context.Context, fromPath, toPath string) error {
	panic("implement me")
}

func (c *client) Exec(ctx context.Context, cmd *rfs.Cmd) error {
	session, err := c.cc.NewSession()
	if err != nil {
		return err
	}

	shell := cmd.Name
	for _, arg := range cmd.Args {
		shell += " " + arg
	}

	for _, env := range cmd.Env {
		parts := strings.Split(env, "=")
		if len(parts) > 1 {
			if err = session.Setenv(parts[0], parts[1]); err != nil {
				return err
			}
		}
	}

	session.Stdin = cmd.Stdin
	session.Stdout = cmd.Stdout
	session.Stderr = cmd.Stderr

	ech := make(chan error, 1)
	done := make(chan struct{}, 1)
	go func() {
		ee := session.Run(shell)
		if ee != nil {
			ech <- ee
		} else {
			done <- struct{}{}
		}
	}()

	select {
	case <-ctx.Done():
		return rfs.ErrTimeout
	case err = <-ech:
		return err
	case <-done:
		return nil
	}
}

func New(opts ...Option) *client {
	var options Options

	for _, opt := range opts {
		opt(&options)
	}

	return &client{options: options}
}
