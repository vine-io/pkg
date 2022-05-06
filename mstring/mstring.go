// MIT License
//
// Copyright (c) 2022 Lack
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

package mstring

import (
	"bytes"
	"math/rand"
	"time"
)

const (
	lowercase   = "abcdefghijklmnopqrstuvwxyz"
	uppercase   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digit       = "0123456"
	punctuation = "!\"#$%&\\'()*+,-./:;<=>?@[\\\\]^_`{|}~"
)

type Options struct {
	c          int
	l, u, d, p bool
}

func newOptions(opts ...Option) Options {
	options := Options{}

	for _, opt := range opts {
		opt(&options)
	}

	if options.c == 0 {
		options.c = 4
	}

	if !options.l && !options.u && !options.d && !options.p {
		options.l = true
	}

	return options
}

type Option func(*Options)

func C(c int) Option {
	return func(o *Options) {
		o.c = c
	}
}

func Lowercase() Option {
	return func(o *Options) {
		o.l = true
	}
}

func Uppercase() Option {
	return func(o *Options) {
		o.u = true
	}
}

func Digit() Option {
	return func(o *Options) {
		o.d = true
	}
}

func Punctuation() Option {
	return func(o *Options) {
		o.p = true
	}
}

func Gen(opts ...Option) string {
	options := newOptions(opts...)

	b := bytes.NewBuffer([]byte(""))
	if options.l {
		b.WriteString(lowercase)
	}
	if options.u {
		b.WriteString(uppercase)
	}
	if options.d {
		b.WriteString(digit)
	}
	if options.p {
		b.WriteString(punctuation)
	}
	target := b.String()
	length := len(target)

	rand.Seed(time.Now().UnixNano())
	out := bytes.NewBuffer([]byte(""))
	for i := 0; i < options.c; i++ {
		n := rand.Intn(length)
		out.WriteByte(target[n])
	}
	return out.String()
}
