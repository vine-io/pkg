package winrm

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/vine-io/pkg/rfs"
	"github.com/masterzen/winrm"
)

var _ rfs.Rfs = (*client)(nil)

type client struct {
	opts Options

	cc *winrm.Client
}

func (c *client) Init() error {

	var err error
	c.cc, err = winrm.NewClient(&c.opts.Endpoint, c.opts.username, c.opts.password)
	if err != nil {
		return fmt.Errorf("%w: %v", rfs.ErrConnect, err)
	}

	return nil
}

func (c *client) Exec(ctx context.Context, cmd *rfs.Cmd, opts ...rfs.ExecOption) error {
	if cmd == nil {
		return rfs.ErrEmptyCmd
	}
	if cmd.Stdout == nil || cmd.Stderr == nil {
		return rfs.ErrMissingCmd
	}

	var (
		err   error
		shell = cmd.Name
		ech   = make(chan error, 1)
		done  = make(chan struct{}, 1)
	)
	for _, arg := range cmd.Args {
		shell += " " + arg
	}
	go func() {
		var ee error
		var bash *winrm.Shell
		bash, ee = c.cc.CreateShell()
		if ee != nil {
			ech <- ee
			return
		}
		defer bash.Close()

		var command *winrm.Command
		command, ee = bash.Execute(cmd.Name, cmd.Args...)
		if ee != nil {
			ech <- ee
			return
		}

		var wg sync.WaitGroup
		wg.Add(3)

		go func() {
			if cmd.Stdin == nil {
				wg.Done()
				return
			}

			defer func() {
				command.Stdin.Close()
				wg.Done()
			}()
			io.Copy(command.Stdin, cmd.Stdin)
		}()
		go func() {
			defer wg.Done()
			io.Copy(cmd.Stdout, command.Stdout)
		}()
		go func() {
			defer wg.Done()
			io.Copy(cmd.Stderr, command.Stderr)
		}()

		command.Wait()
		wg.Wait()
		command.Close()

		if command.ExitCode() != 0 {
			stderr := []byte("")
			stdout := []byte("")
			cmd.Stdout.Write(stdout)
			cmd.Stderr.Write(stderr)
			ech <- fmt.Errorf("%v: %v", string(stderr), string(stderr))
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
	cc, err := winrm.NewClient(&c.opts.Endpoint, c.opts.username, c.opts.password)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", rfs.ErrConnect, err)
	}

	var (
		ech   = make(chan error, 1)
		done  = make(chan struct{}, 1)
		items = make([]os.FileInfo, 0)
	)

	go func() {
		var ee error
		items, ee = fetchList(cc, remotePath)
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
	var (
		err  error
		ech  = make(chan error, 1)
		done = make(chan struct{}, 1)
	)

	go func() {
		info, ee := fetchList(c.cc, remotePath)
		if ee != nil {
			ech <- fmt.Errorf("%w: remote path %v", rfs.ErrNotExists, remotePath)
			return
		}
		if len(info) == 0 {
			ech <- fmt.Errorf("%w: empty directory", rfs.ErrNotExists)
			return
		}

		if len(info) == 1 {
			ee = c.get(ctx, remotePath, localPath, fn)
		} else {
			_ = os.MkdirAll(localPath, 0755)
			ee = c.walker(ctx, info, remotePath, localPath, fn)
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
	var (
		err  error
		ech  = make(chan error, 1)
		done = make(chan struct{}, 1)
	)

	go func() {
		var ee error

		ee = c.Copy(localPath, remotePath, fn)
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

	return &client{opts: options}
}
