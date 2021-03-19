// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package zfs

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

func SendPipe(ctx context.Context, source, last string, fn func(context.Context, io.ReadCloser) error) error {
	var (
		shell string
		done  = make(chan struct{})
		errCh = make(chan error, 1)
	)

	if last != "" {
		shell = fmt.Sprintf("%s send -i %s %s", "zfs", last, source)
	} else {
		shell = fmt.Sprintf("%s send %s", "zfs", source)
	}
	log.Println(shell)
	sender := exec.CommandContext(ctx, "/bin/sh", "-c", shell)
	sender.Stderr = os.Stderr
	rd, err := sender.StdoutPipe()
	if err != nil {
		return err
	}

	go func() {
		if err := fn(ctx, rd); err != nil && err != io.EOF {
			errCh <- fmt.Errorf("%w: send failed")
		} else {
			close(done)
		}
	}()

	go func() {
		if err := sender.Run(); err != nil {
			errCh <- fmt.Errorf("sender %v", err)
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err = <-errCh:
		return err
	case <-done:
		return nil
	}
}

func RecvPipe(ctx context.Context, target string, fn func(context.Context, io.WriteCloser) error) error {
	var (
		done  = make(chan struct{})
		errCh = make(chan error, 1)
	)

	shell := fmt.Sprintf("zfs receive -Fu %s", target)
	log.Println(shell)
	recver := exec.CommandContext(ctx, "/bin/sh", "-c", shell)
	recver.Stderr = os.Stderr
	wd, err := recver.StdinPipe()
	if err != nil {
		return err
	}

	go func() {
		if err := fn(ctx, wd); err != nil && err != io.EOF {
			errCh <- fmt.Errorf("%w: receive failed", err)
		} else {
			close(done)
		}
	}()

	go func() {
		if err := recver.Run(); err != nil {
			errCh <- fmt.Errorf("recer %v", err)
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err = <-errCh:
		return err
	case <-done:
		return nil
	}
}
