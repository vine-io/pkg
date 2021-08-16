package winrm

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/vine-io/pkg/rfs"
)

func Test_client_Exec(t *testing.T) {
	cc := New(
		Host("192.168.1.152", 5985),
		Auth("Administrator", "Cisco123"),
		Insecure(true),
		Timeout(time.Second*3),
	)

	cc.Init()

	ctx := context.TODO()
	b := bytes.NewBuffer([]byte(""))
	c := rfs.NewCmd("ipconfig /all", nil, b)
	if err := cc.Exec(ctx, c); err != nil {
		t.Fatal(err)
	}
	t.Log(string(b.Bytes()))
}

func Test_client_List(t *testing.T) {
	cc := New(
		Host("192.168.1.152", 5985),
		Auth("Administrator", "Cisco123"),
		Insecure(true),
		Timeout(time.Second*3),
	)

	ctx := context.TODO()
	outs, err := cc.List(ctx, "C:\\howlink\\aaa")
	if err != nil {
		t.Fatal(err)
	}

	for _, stat := range outs {
		t.Logf("%v | %v | %v | %v | %v", stat.Name(), stat.Size(), stat.Mode(), stat.ModTime(), stat.IsDir())
	}
}

func Test_client_Get(t *testing.T) {

	cc := New(
		Host("192.168.1.152", 5985),
		Auth("Administrator", "Cisco123"),
		Insecure(true),
		Timeout(time.Second*3),
	)

	if err := cc.Init(); err != nil {
		t.Fatal(err)
	}

	ctx := context.TODO()
	src := "C:\\howlink\\aaa\\hello"
	dst := "/tmp/aaa/hello"
	err := cc.Get(ctx, src, dst, func(metric *rfs.IOMetric) {
		t.Log(metric)
	})
	if err != nil {
		t.Fatal(err)
	}
}

func Test_client_Put(t *testing.T) {

	cc := New(
		Host("192.168.1.152", 5985),
		Auth("Administrator", "Cisco123"),
		Insecure(true),
		Timeout(time.Second*3),
	)

	if err := cc.Init(); err != nil {
		t.Fatal(err)
	}

	ctx := context.TODO()
	src := "/tmp/sfz.exe"
	dst := "C:\\howlink\\aaa\\sfz.exe"
	err := cc.Put(ctx, src, dst, func(metric *rfs.IOMetric) {
		t.Log(metric)
	})
	if err != nil {
		t.Fatal(err)
	}
}
