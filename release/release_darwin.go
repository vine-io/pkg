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

package release

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func _get() (OSRelease, error) {
	r := OSRelease{}

	r.Kernel = runtime.GOOS
	r.Arch = runtime.GOARCH
	r.GoV = runtime.Version()

	cmd := exec.Command("sw_vers")
	cmd.Env = os.Environ()
	out, err := cmd.CombinedOutput()
	if err != nil {
		return r, err
	}

	rd := bufio.NewReader(bytes.NewReader(out))
	for {
		line, _, err := rd.ReadLine()
		if err == io.EOF {
			break
		}
		text := string(line)
		if strings.HasPrefix(text, "ProductName:") {
			t := strings.TrimSpace(strings.TrimPrefix(text, "ProductName:"))
			r.Name = t
		}
		if strings.HasPrefix(text, "ProductVersion:") {
			t := strings.TrimSpace(strings.TrimPrefix(text, "ProductVersion:"))
			r.Version = t
		}
	}

	return r, nil
}