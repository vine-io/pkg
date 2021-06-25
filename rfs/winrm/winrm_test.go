package winrm

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/lack-io/pkg/rfs"
)

func Test_client_Exec(t *testing.T) {
	cc := New(
		Host("192.168.1.152", 5985),
		Auth("Administrator", "Cisco123"),
		Insecure(true),
		Timeout(time.Second*3),
	)

	ctx := context.TODO()
	cmd := rfs.Cmd{
		Name:   "ipconfig /all",
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	if err := cc.Exec(ctx, &cmd); err != nil {
		t.Fatal(err)
	}
}

func Test_client_List(t *testing.T) {
	cc := New(
		Host("192.168.1.152", 5985),
		Auth("Administrator", "Cisco123"),
		Insecure(true),
		Timeout(time.Second*3),
	)

	ctx := context.TODO()
	outs, err := cc.List(ctx, "C:\\howlink")
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
	src := "C:\\howlink\\aaa\\winrm.go"
	dst := "/tmp/winrm.go"
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
	src := "winrm.go"
	dst := "C:\\howlink\\aaa\\winrm.go"
	err := cc.Put(ctx, src, dst, func(metric *rfs.IOMetric) {
		t.Log(metric)
	})
	if err != nil {
		t.Fatal(err)
	}
}
