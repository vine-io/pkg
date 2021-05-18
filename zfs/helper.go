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

package zfs

import (
	"bufio"
	"bytes"
	"io"
	"reflect"
	"strings"
)

func setValue(data []byte, into interface{}) {
	reader := bytes.NewReader(data)
	rd := bufio.NewReader(reader)
	for {
		line, err := rd.ReadString('\n')
		if err != nil && err != io.EOF {
			// 读取中出现错误直接退出
			break
		}

		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		set(into, line, "\t")
		if err == io.EOF {
			// 读取到尾部，返回
			return
		}
	}
}

// 利用反射动态设置 target
// @target: Target install
// @data: 需要处理的字符串
// @slim: 分隔符
func set(target interface{}, data string, slim string) {
	line := strings.Split(data, slim)
	getType := reflect.TypeOf(target).Elem()
	valueType := reflect.ValueOf(target).Elem()
	for i := 0; i < getType.NumField(); i++ {
		field := getType.Field(i)
		if field.Tag.Get("zfs") == line[1] {
			valueType.Field(i).SetString(line[2])
		}
	}
}
