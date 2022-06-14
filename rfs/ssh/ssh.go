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
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/sftp"
	"github.com/vine-io/pkg/rfs"
	"golang.org/x/crypto/ssh"
)

var _ rfs.Rfs = (*client)(nil)

type client struct {
	options Options
}

func (c *client) Init() error {

	var err error
	_, err = ssh.Dial(c.options.network, c.options.addr, &c.options.ClientConfig)
	if err != nil {
		return fmt.Errorf("%w: %v", rfs.ErrConnect, err)
	}
	return nil
}

func (c *client) dial() (*ssh.Client, error) {
	cc, err := ssh.Dial(c.options.network, c.options.addr, &c.options.ClientConfig)
	if err != nil {
		return nil, err
	}
	return cc, nil
}

type singleWriter struct {
	b  bytes.Buffer
	mu sync.Mutex
}

func (w *singleWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.b.Write(p)
}

func (c *client) Exec(ctx context.Context, cmd *rfs.Cmd, opts ...rfs.ExecOption) error {
	if cmd == nil {
		return rfs.ErrEmptyCmd
	}

	cc, err := c.dial()
	if err != nil {
		return err
	}
	defer cc.Close()

	session, err := cc.NewSession()
	if err != nil {
		return fmt.Errorf("%w: %v", rfs.ErrConnect, err)
	}
	defer session.Close()

	shell := cmd.Name
	for _, arg := range cmd.Args {
		shell += " " + arg
	}

	session.Stdin = cmd.Stdin
	session.Stdout = cmd.Stdout
	session.Stderr = cmd.Stderr

	ech := make(chan error, 1)
	done := make(chan struct{}, 1)
	go func() {
		var ee error
		for _, env := range cmd.Env {
			parts := strings.Split(env, "=")
			if len(parts) > 1 {
				if ee = session.Setenv(parts[0], parts[1]); ee != nil {
					ech <- ee
					return
				}
			}
		}

		if ee = session.Run(shell); ee != nil {
			ech <- ee
			return
		}
		done <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return rfs.ErrTimeout
	case err = <-ech:
		return fmt.Errorf("%w: %v", rfs.ErrRequest, err)
	case <-done:
		return nil
	}
}

func (c *client) List(ctx context.Context, remotePath string, opts ...rfs.ListOption) ([]os.FileInfo, error) {
	cc, err := c.dial()
	if err != nil {
		return nil, err
	}
	defer cc.Close()

	ftp, err := sftp.NewClient(cc)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", rfs.ErrConnect, err)
	}
	defer ftp.Close()

	var (
		ech   = make(chan error, 1)
		done  = make(chan struct{}, 1)
		items = make([]os.FileInfo, 0)
	)

	go func() {
		var ee error
		items, ee = ftp.ReadDir(remotePath)
		if ee != nil {
			ech <- ee
			return
		}

		done <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return nil, rfs.ErrTimeout
	case err = <-ech:
		return nil, fmt.Errorf("%w: %v", rfs.ErrRequest, err)
	case <-done:
	}

	return items, nil
}

func (c *client) Get(ctx context.Context, remotePath, localPath string, fn rfs.IOFn, opts ...rfs.GetOption) error {
	cc, err := c.dial()
	if err != nil {
		return err
	}
	defer cc.Close()

	ftp, err := sftp.NewClient(cc)
	if err != nil {
		return fmt.Errorf("%w: %v", rfs.ErrConnect, err)
	}
	defer ftp.Close()

	stat, err := ftp.Stat(remotePath)
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("%w: %v", rfs.ErrNotExists, err)
	}

	var (
		buf  = make([]byte, 32*1024)
		ech  = make(chan error, 1)
		done = make(chan struct{}, 1)
	)

	go func() {
		var ee error

		if stat.IsDir() {
			lstat, e1 := os.Stat(localPath)
			if e1 != nil {
				ech <- err
				return
			}
			if !lstat.IsDir() {
				ech <- fmt.Errorf("%w: %v", rfs.ErrAlreadyExists, localPath)
				return
			}

			walker := ftp.Walk(remotePath)
			for walker.Step() {
				if walker.Err() != nil || walker.Path() == remotePath {
					continue
				}

				sub := strings.TrimPrefix(walker.Path(), remotePath)
				if walker.Stat().IsDir() {
					_ = os.MkdirAll(filepath.Join(localPath, sub), os.ModePerm)
					continue
				}
				dst := filepath.Join(localPath, sub)
				if _, ee = c.get(ctx, ftp, walker.Path(), dst, buf, fn); ee != nil {
					break
				}
			}

		} else {
			lstat, _ := os.Stat(localPath)
			if lstat != nil && lstat.IsDir() {
				localPath = filepath.Join(localPath, filepath.Base(remotePath))
			}

			_, ee = c.get(ctx, ftp, remotePath, localPath, buf, fn)
		}

		if ee != nil {
			ech <- ee
			return
		}

		done <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return rfs.ErrTimeout
	case err = <-ech:
		return fmt.Errorf("%w: %v", rfs.ErrRequest, err)
	case <-done:
		return nil
	}
}

func (c *client) Put(ctx context.Context, localPath, remotePath string, fn rfs.IOFn, opts ...rfs.PutOption) error {
	cc, err := c.dial()
	if err != nil {
		return err
	}
	defer cc.Close()

	ftp, err := sftp.NewClient(cc)
	if err != nil {
		return fmt.Errorf("%w: %v", rfs.ErrConnect, err)
	}
	defer ftp.Close()

	stat, err := os.Stat(localPath)
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("%w: %v", rfs.ErrNotExists, err)
	}

	var (
		buf  = make([]byte, 32*1024)
		ech  = make(chan error, 1)
		done = make(chan struct{}, 1)
	)

	go func() {
		var ee error

		if stat.IsDir() {
			rstat, e1 := ftp.Stat(remotePath)
			if e1 != nil {
				ech <- err
				return
			}
			if !rstat.IsDir() {
				ech <- fmt.Errorf("%w: %v", rfs.ErrAlreadyExists, remotePath)
				return
			}

			ee = filepath.WalkDir(localPath, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if path == localPath {
					return nil
				}

				sub := strings.TrimPrefix(path, localPath)
				if d.IsDir() {
					return ftp.MkdirAll(filepath.Join(remotePath, sub))
				}
				dst := filepath.Join(remotePath, sub)
				_, e1 := c.put(ctx, ftp, path, dst, buf, fn)
				return e1
			})

		} else {
			rstat, _ := ftp.Stat(remotePath)
			if rstat != nil && rstat.IsDir() {
				remotePath = filepath.Join(remotePath, filepath.Base(localPath))
			}

			_, ee = c.put(ctx, ftp, localPath, remotePath, buf, fn)
		}

		if ee != nil {
			ech <- ee
			return
		}

		done <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return rfs.ErrTimeout
	case err = <-ech:
		return fmt.Errorf("%w: %v", rfs.ErrRequest, err)
	case <-done:
		return nil
	}
}

func New(opts ...Option) *client {
	var options Options

	for _, opt := range opts {
		opt(&options)
	}

	if options.network == "" {
		options.network = "tcp"
	}
	if options.HostKeyCallback == nil {
		options.HostKeyCallback = func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		}
	}

	return &client{options: options}
}
