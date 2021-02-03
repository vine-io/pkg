package scst

import (
	"reflect"
	"testing"
)

func Test_readDirs(t *testing.T) {
	type args struct {
		f       string
		ignores []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "test-readDirs-1", args: args{f: "../../"}, want: []string{"iscsi", "testdata"}},
		{name: "test-readDirs-1", args: args{f: "../../", ignores: []string{"iscsi"}}, want: []string{"testdata"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := readDirs(tt.args.f, tt.args.ignores...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readDirs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readFiles(t *testing.T) {
	type args struct {
		f       string
		ignores []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "test-readFiles-1", args: args{f: "../../testdata", ignores: []string{"walk"}}, want: []string{"main.go"}},
		{name: "test-readFiles-2", args: args{f: "../../testdata", ignores: []string{""}}, want: []string{"walk", "main.go"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := readFiles(tt.args.f, tt.args.ignores...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}
