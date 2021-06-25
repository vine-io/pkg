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
	"os"
	"path/filepath"
	"time"

	"github.com/lack-io/pkg/rfs"
	"github.com/pkg/sftp"
)

func (c *client) get(
	ctx context.Context,
	ftp *sftp.Client,
	src string,
	dst string,
	buf []byte,
	fn rfs.IOFn,
) (written int64, err error) {

	reader, err := ftp.Open(src)
	if err != nil {
		return 0, err
	}
	defer reader.Close()

	writer, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer writer.Close()

	var metric *rfs.IOMetric
	if fn != nil {
		metric = &rfs.IOMetric{
			Name: filepath.Base(reader.Name()),
			From: src,
			To:   dst,
		}
		if stat, _ := reader.Stat(); stat != nil {
			metric.Total = stat.Size()
		}
	}

	return fcopy(ctx, reader, writer, metric, buf, fn)
}

func (c *client) put(
	ctx context.Context,
	ftp *sftp.Client,
	src string,
	dst string,
	buf []byte,
	fn rfs.IOFn,
) (written int64, err error) {

	reader, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer reader.Close()

	writer, err := ftp.Create(dst)
	if err != nil {
		return 0, err
	}
	defer writer.Close()

	var metric *rfs.IOMetric
	if fn != nil {
		metric = &rfs.IOMetric{
			Name: reader.Name(),
			From: src,
			To:   dst,
		}
		if stat, _ := reader.Stat(); stat != nil {
			metric.Total = stat.Size()
		}
	}

	return fcopy(ctx, reader, writer, metric, buf, fn)
}

func fcopy(
	ctx context.Context,
	reader io.Reader,
	writer io.Writer,
	metric *rfs.IOMetric,
	buf []byte,
	fn rfs.IOFn,
) (written int64, err error) {
	if buf == nil {
		buf = make([]byte, 32*1024)
	}

	last := time.Now()
	sub := int64(0)
	for {
		select {
		case <-ctx.Done():
			return 0, rfs.ErrTimeout
		default:
		}

		nr, er := reader.Read(buf)
		if nr > 0 {
			nw, ew := writer.Write(buf[0:nr])
			if nw < 0 || nr < nw {
				nw = 0
				if ew == nil {
					ew = errInvalidWrite
				}
			}
			written += int64(nw)
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if fn != nil {
			now := time.Now()
			metric.Block = written
			metric.Speed = int64(float64(written-sub) / (now.Sub(last).Seconds()))
			last = now
			sub = written
			fn(metric)
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err
}
