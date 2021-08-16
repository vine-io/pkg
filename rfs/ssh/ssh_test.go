package ssh

import (
	"bytes"
	"context"
	"testing"

	"github.com/vine-io/pkg/rfs"
)

func Test_client_Exec(t *testing.T) {

	cc := New(Host("192.168.3.111:22"), Auth("root", "123456"))
	if err := cc.Init(); err != nil {
		t.Fatal(err)
	}

	ctx := context.TODO()
	b := bytes.NewBuffer([]byte(""))
	c := rfs.NewCmd("ifconfig", nil, b)
	if err := cc.Exec(ctx, c); err != nil {
		t.Fatal(err)
	}
	t.Log(string(b.Bytes()))
}

func Test_client_List(t *testing.T) {

	cc := New(Host("192.168.3.111:22"), Auth("root", "123456"))
	if err := cc.Init(); err != nil {
		t.Fatal(err)
	}

	ctx := context.TODO()
	stats, err := cc.List(ctx, "/")
	if err != nil {
		t.Fatal(err)
	}

	for _, stat := range stats {
		t.Logf("%v | %v | %v | %v | %v", stat.Name(), stat.Size(), stat.Mode(), stat.ModTime(), stat.IsDir())
	}
}

func Test_client_Get(t *testing.T) {

	cc := New(Host("192.168.3.111:22"), Auth("root", "123456"))
	if err := cc.Init(); err != nil {
		t.Fatal(err)
	}

	ctx := context.TODO()
	src := "/root/EMC-ScaleIO-gateway-2.6-11000.113.x86_64.rpm"
	dst := "/tmp/"
	err := cc.Get(ctx, src, dst, func(metric *rfs.IOMetric) {
		t.Log(metric)
	})
	if err != nil {
		t.Fatal(err)
	}
}

func Test_client_Put(t *testing.T) {

	cc := New(Host("192.168.3.111:22"), Auth("root", "123456"))
	if err := cc.Init(); err != nil {
		t.Fatal(err)
	}

	ctx := context.TODO()
	src := "ssh.go"
	dst := "/root/aaa/ssh.go"
	err := cc.Put(ctx, src, dst, func(metric *rfs.IOMetric) {
		t.Log(metric)
	})
	if err != nil {
		t.Fatal(err)
	}
}
