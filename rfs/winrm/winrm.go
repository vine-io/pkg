package winrm

import (
	"context"
	"fmt"
	"io"

	"github.com/lack-io/pkg/rfs"
	"github.com/masterzen/winrm"
)

var _ rfs.Rfs = (*client)(nil)

type client struct {
	opts Options

	cc *winrm.Client
}

func New(opts ...Option) *client {
	var options Options

	for _, opt := range opts {
		opt(&options)
	}

	return &client{opts: options}
}

func (c *client) Init() error {
	return nil
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
	if cmd == nil {
		return rfs.ErrEmptyCmd
	}

	cc, err := winrm.NewClient(&c.opts.Endpoint, c.opts.username, c.opts.password)
	if err != nil {
		return fmt.Errorf("connect to server: %w", err)
	}

	var (
		shell = cmd.Name
		ech   = make(chan error, 1)
		done  = make(chan struct{}, 1)
	)
	for _, arg := range cmd.Args {
		shell += " " + arg
	}
	go func() {
		_, err = cc.RunWithInput(shell, cmd.Stdout, cmd.Stderr, cmd.Stdin)
		if err != nil {
			ech <- err
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
