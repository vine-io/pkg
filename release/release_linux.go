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
	"fmt"
	"io"
	"os"
	"regexp"
	"runtime"
	"strings"
)

func _get() (OSRelease, error) {
	r := OSRelease{}

	r.Kernel = runtime.GOOS
	r.Arch = runtime.GOARCH
	r.GoV = runtime.Version()

	target := "/etc/os-release"
	stat, _ := os.Stat(target)
	if stat != nil {
		fd, err := os.Open(target)
		if err != nil {
			return r, nil
		}

		rd := bufio.NewReader(fd)
		for {
			line, _, err := rd.ReadLine()
			if err == io.EOF {
				break
			}
			text := string(line)
			if strings.HasPrefix(text, "ID=") {
				parts := strings.Split(text, "=")
				r.Name = strings.Trim(strings.TrimSpace(parts[1]), `"`)
			}
			if strings.HasPrefix(text, "VERSION_ID=") {
				parts := strings.Split(text, "=")
				r.Version = strings.Trim(strings.TrimSpace(parts[1]), `"`)
			}
		}
	} else {
		target = "/etc/system-release"
		stat, _ = os.Stat(target)
		if stat == nil {
			return r, fmt.Errorf("unknown linux release")
		}

		release, err := os.ReadFile(target)
		if err != nil && err != io.EOF {
			return r, err
		}
		// CentOS and Red Hat
		re := regexp.MustCompile(`(CentOS|Red Hat) .* release (\d)\..* \(.*\)`)
		sm := re.FindStringSubmatch(string(release))
		if len(sm) > 2 {
			switch sm[1] {
			case "CentOS":
				r.Name = "centos"
			case "Red Hat":
				r.Name = "rhel"
			}

			r.Version = sm[2]
		}
	}

	return r, nil
}
