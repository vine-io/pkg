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

package winrm

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/lack-io/pkg/rfs"
	"github.com/masterzen/winrm"
	"github.com/nu7hatch/gouuid"
)

func (c *client) Copy(fromPath, toPath string, fn rfs.IOFn) error {
	f, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("couldn't read file %s: %v", fromPath, err)
	}

	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return fmt.Errorf("couldn't stat file %s: %v", fromPath, err)
	}

	if !fi.IsDir() {
		return c.Write(toPath, f, fn)
	} else {
		fw := fileWalker{
			client:  c.cc,
			config:  c.opts,
			toDir:   toPath,
			fromDir: fromPath,
			fn:      fn,
		}
		return filepath.Walk(fromPath, fw.copyFile)
	}
}

func (c *client) Write(toPath string, src *os.File, fn rfs.IOFn) error {
	return doCopy(c.cc, c.opts, src, winPath(toPath), fn)
}

type fileWalker struct {
	client  *winrm.Client
	config  Options
	toDir   string
	fromDir string
	fn      rfs.IOFn
}

func (fw *fileWalker) copyFile(fromPath string, fi os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if !shouldUploadFile(fi) {
		return nil
	}

	hostPath, _ := filepath.Abs(fromPath)
	fromDir, _ := filepath.Abs(fw.fromDir)
	relPath, _ := filepath.Rel(fromDir, hostPath)
	toPath := filepath.Join(fw.toDir, relPath)

	f, err := os.Open(hostPath)
	if err != nil {
		return fmt.Errorf("couldn't read file %s: %v", fromPath, err)
	}

	return doCopy(fw.client, fw.config, f, winPath(toPath), fw.fn)
}

func shouldUploadFile(fi os.FileInfo) bool {
	// Ignore dir entries and OS X special hidden file
	return !fi.IsDir() && ".DS_Store" != fi.Name()
}

func doCopy(client *winrm.Client, config Options, in *os.File, toPath string, fn rfs.IOFn) error {
	tempFile, err := tempFileName()
	if err != nil {
		return fmt.Errorf("error generating unique filename: %v", err)
	}
	tempPath := "$env:TEMP\\" + tempFile

	defer func() {
		_ = cleanupContent(client, tempPath)
	}()

	var metric *rfs.IOMetric
	if fn != nil {
		metric = &rfs.IOMetric{
			Name: filepath.Base(in.Name()),
			From: in.Name(),
			To:   toPath,
		}
		stat, _ := in.Stat()
		if stat != nil {
			metric.Total = stat.Size()
		}
	}
	err = uploadContent(client, 32*1024, "%TEMP%\\"+tempFile, in, metric, fn)
	if err != nil {
		return fmt.Errorf("error uploading file to %s: %v", tempPath, err)
	}

	err = restoreContent(client, tempPath, toPath)
	if err != nil {
		return fmt.Errorf("error restoring file from %s to %s: %v", tempPath, toPath, err)
	}

	return nil
}

func uploadContent(client *winrm.Client, maxChunks int, filePath string, reader io.Reader, metric *rfs.IOMetric, fn rfs.IOFn) error {
	var err error
	done := false
	for !done {
		done, err = uploadChunks(client, filePath, maxChunks, reader, metric, fn)
		if err != nil {
			return err
		}
	}

	return nil
}

func uploadChunks(client *winrm.Client, filePath string, maxChunks int, reader io.Reader, metric *rfs.IOMetric, fn rfs.IOFn) (bool, error) {
	shell, err := client.CreateShell()
	if err != nil {
		return false, fmt.Errorf("couldn't create shell: %v", err)
	}
	defer shell.Close()

	// Upload the file in chunks to get around the Windows command line size limit.
	// Base64 encodes each set of three bytes into four bytes. In addition the output
	// is padded to always be a multiple of four.
	//
	//   ceil(n / 3) * 4 = m1 - m2
	//
	//   where:
	//     n  = bytes
	//     m1 = max (8192 character command limit.)
	//     m2 = len(filePath)

	chunkSize := ((8000 - len(filePath)) / 4) * 3
	chunk := make([]byte, chunkSize)

	if maxChunks == 0 {
		maxChunks = 1
	}

	last := time.Now()
	for i := 0; i < maxChunks; i++ {
		n, err := reader.Read(chunk)

		if err != nil && err != io.EOF {
			return false, err
		}

		if fn != nil {
			now := time.Now()
			metric.Block += int64(n)
			metric.Speed = int64(float64(n) / (now.Sub(last).Seconds()))
			last = now
			fn(metric)
		}

		if n == 0 {
			return true, nil
		}

		content := base64.StdEncoding.EncodeToString(chunk[:n])
		if err = appendContent(shell, filePath, content); err != nil {
			return false, err
		}
	}

	return false, nil
}

func restoreContent(client *winrm.Client, fromPath, toPath string) error {
	shell, err := client.CreateShell()
	if err != nil {
		return err
	}

	defer shell.Close()
	script := fmt.Sprintf(`
		$tmp_file_path = [System.IO.Path]::GetFullPath("%s")
		$dest_file_path = [System.IO.Path]::GetFullPath("%s".Trim("'"))
		if (Test-Path $dest_file_path) {
			if (Test-Path -Path $dest_file_path -PathType container) {
				Exit 1
			} else {
				rm $dest_file_path
			}
		}
		else {
			$dest_dir = ([System.IO.Path]::GetDirectoryName($dest_file_path))
			New-Item -ItemType directory -Force -ErrorAction SilentlyContinue -Path $dest_dir | Out-Null
		}

		if (Test-Path $tmp_file_path) {
			$reader = [System.IO.File]::OpenText($tmp_file_path)
			$writer = [System.IO.File]::OpenWrite($dest_file_path)
			try {
				for(;;) {
					$base64_line = $reader.ReadLine()
					if ($base64_line -eq $null) { break }
					$bytes = [System.Convert]::FromBase64String($base64_line)
					$writer.write($bytes, 0, $bytes.Length)
				}
			}
			finally {
				$reader.Close()
				$writer.Close()
			}
		} else {
			echo $null > $dest_file_path
		}
	`, fromPath, toPath)

	cmd, err := shell.Execute(winrm.Powershell(script))
	if err != nil {
		return err
	}
	defer cmd.Close()

	var wg sync.WaitGroup
	copyFunc := func(w io.Writer, r io.Reader) {
		defer wg.Done()
		io.Copy(w, r)
	}

	wg.Add(2)
	go copyFunc(os.Stdout, cmd.Stdout)
	go copyFunc(os.Stderr, cmd.Stderr)

	cmd.Wait()
	wg.Wait()

	if cmd.ExitCode() != 0 {
		return fmt.Errorf("restore operation returned code=%d", cmd.ExitCode())
	}
	return nil
}

func cleanupContent(client *winrm.Client, filePath string) error {
	shell, err := client.CreateShell()
	if err != nil {
		return err
	}

	defer shell.Close()
	script := fmt.Sprintf(`
		$tmp_file_path = [System.IO.Path]::GetFullPath("%s")
		if (Test-Path $tmp_file_path) {
			Remove-Item $tmp_file_path -ErrorAction SilentlyContinue
		}
	`, filePath)

	cmd, err := shell.Execute(winrm.Powershell(script))
	if err != nil {
		return err
	}
	defer cmd.Close()

	var wg sync.WaitGroup
	copyFunc := func(w io.Writer, r io.Reader) {
		defer wg.Done()
		io.Copy(w, r)
	}

	wg.Add(2)
	go copyFunc(os.Stdout, cmd.Stdout)
	go copyFunc(os.Stderr, cmd.Stderr)

	cmd.Wait()
	wg.Wait()

	if cmd.ExitCode() != 0 {
		return fmt.Errorf("cleanup operation returned code=%d", cmd.ExitCode())
	}
	return nil
}

func appendContent(shell *winrm.Shell, filePath, content string) error {
	cmd, err := shell.Execute(fmt.Sprintf(`echo %s >> %s`, content, filePath))

	if err != nil {
		return err
	}

	defer cmd.Close()
	var wg sync.WaitGroup
	copyFunc := func(w io.Writer, r io.Reader) {
		defer wg.Done()
		io.Copy(w, r)
	}

	wg.Add(2)
	go copyFunc(os.Stdout, cmd.Stdout)
	go copyFunc(os.Stderr, cmd.Stderr)

	cmd.Wait()
	wg.Wait()

	if cmd.ExitCode() != 0 {
		return fmt.Errorf("upload operation returned code=%d", cmd.ExitCode())
	}

	return nil
}

func tempFileName() (string, error) {
	uniquePart, err := uuid.NewV4()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("winrmcp-%s.tmp", uniquePart), nil
}

func (c *client) get(ctx context.Context, remotePath, localPath string, fn rfs.IOFn) error {
	writer, err := os.Create(localPath)
	if err != nil {
		return err
	}
	defer writer.Close()

	var metric *rfs.IOMetric
	if fn != nil {
		metric = &rfs.IOMetric{
			Name: filepath.Base(remotePath),
			From: remotePath,
			To:   localPath,
		}

		info, _ := fetchList(c.cc, remotePath)
		if len(info) > 0 {
			metric.Total = info[0].Size()
		}
	}

	return readContent(ctx, c.cc, remotePath, writer, metric, fn)
}

func readContent(ctx context.Context, client *winrm.Client, remotePath string, writer *os.File, metric *rfs.IOMetric, fn rfs.IOFn) error {
	shell, err := client.CreateShell()
	if err != nil {
		return err
	}

	defer shell.Close()
	script := fmt.Sprintf(`
		$dest_file_path = [System.IO.Path]::GetFullPath("%s".Trim("'"))
		if (Test-Path $dest_file_path) {
			if (Test-Path -Path $dest_file_path -PathType container) {
				Exit 1
			}
		}

		$reader = [System.IO.File]::OpenText($dest_file_path)
		try {
			for(;;) {
				$line = $reader.ReadLine()
				if ($line -eq $null) { break }
				echo $line
			}
		}
		finally {
			$reader.Close()
		}
	`, remotePath)

	cmd, err := shell.Execute(winrm.Powershell(script))
	if err != nil {
		return err
	}
	defer cmd.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	buf := []byte("")
	go func() {
		defer wg.Done()
		cmd.Stderr.Read(buf)
	}()
	go func() {
		defer wg.Done()
		io.Copy(writer, cmd.Stdout)
	}()

	cmd.Wait()
	wg.Wait()

	if cmd.ExitCode() != 0 {
		return fmt.Errorf("read file operation returned code=%d: %v", cmd.ExitCode(), string(buf))
	}
	return nil
}
