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
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/masterzen/winrm"
)

var _ os.FileInfo = (*fileInfo)(nil)

// fileInfo is an artificial type designed to satisfy os.FileInfo.
type fileInfo struct {
	name  string
	size  int64
	mode  os.FileMode
	mtime time.Time
	sys   interface{}
}

// Name returns the base name of the file.
func (fi *fileInfo) Name() string { return fi.name }

// Size returns the length in bytes for regular files; system-dependent for others.
func (fi *fileInfo) Size() int64 { return fi.size }

// Mode returns file mode bits.
func (fi *fileInfo) Mode() os.FileMode { return fi.mode }

// ModTime returns the last modification time of the file.
func (fi *fileInfo) ModTime() time.Time { return fi.mtime }

// IsDir returns true if the file is a directory.
func (fi *fileInfo) IsDir() bool { return fi.Mode().IsDir() }

func (fi *fileInfo) Sys() interface{} { return fi.sys }

func fetchList(client *winrm.Client, remotePath string) ([]os.FileInfo, error) {
	script := fmt.Sprintf("Get-ChildItem %s", remotePath)
	stdout, _, _, err := client.RunWithString("powershell -Command \""+script+" | ConvertTo-Xml -NoTypeInformation -As String\"", "")
	if err != nil {
		return nil, fmt.Errorf("couldn't execute script %s: %v", script, err)
	}

	if stdout != "" {
		doc := pslist{}
		err := xml.Unmarshal([]byte(stdout), &doc)
		if err != nil {
			return nil, fmt.Errorf("couldn't parse results: %v", err)
		}

		return convertFileItems(doc.Objects), nil
	}

	return []os.FileInfo{}, nil
}

func convertFileItems(objects []psobject) []os.FileInfo {
	items := make([]os.FileInfo, 0)

	for _, object := range objects {
		stat := &fileInfo{}
		for _, property := range object.Properties {
			switch property.Name {
			case "Name":
				stat.name = property.Value
			case "Mode":
				if property.Value[0] == 'd' {
					stat.mode = os.ModeDir
				} else {
					stat.mode = os.ModeAppend
				}
			//case "FullName":
			//	items[i].Path = property.Value
			case "Length":
				if n, err := strconv.ParseInt(property.Value, 10, 64); err == nil {
					stat.size = n
				}
			case "LastWriteTime":
				stat.mtime, _ = time.Parse("2006/1/02 15:04:05", property.Value)
			}
		}
		stat.sys = struct{}{}
		items = append(items, stat)
	}

	return items
}
