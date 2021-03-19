package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/lack-io/vine/service/client"
	"github.com/lack-io/vine/service/client/grpc"
	log "github.com/lack-io/vine/service/logger"
	zfs "github.com/lack-io/vine/testdata/zfs/proto"
)

var (
	source string
	last   string
	target string
)

func main() {
	flag.StringVar(&source, "source", "", "zfs src snapshot")
	flag.StringVar(&last, "last", "", "zfs increment snapshot")
	flag.StringVar(&target, "target", "", "zfs dest snapshot")
	flag.Parse()

	if len(source) == 0 || len(target) == 0 {
		log.Fatal("must be `source` and `target`")
	}

	shell := ""
	ctx := context.TODO()
	if last != "" {
		shell = fmt.Sprintf("%s send -i %s %s", "zfs", last, source)
	} else {
		shell = fmt.Sprintf("%s send %s", "zfs", source)
	}
	log.Infof(shell)
	send := exec.CommandContext(ctx, "/bin/sh", "-c", shell)
	send.Stderr = os.Stderr
	rd, err := send.StdoutPipe()
	if err != nil {
		log.Fatalf("get stdout pipe: %v", err)
	}

	cli := grpc.NewClient()

	cc := zfs.NewStorageService("go.vine.client", cli)

	stream, err := cc.Recv(ctx, client.WithAddress("127.0.0.1:2333"))
	if err != nil {
		log.Fatal(err)
	}

	in := zfs.RecvRequest{Target: target}
	err = stream.Send(&in)
	if err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 1024*32)
	go func() {
		for {
			n, err := rd.Read(buf)
			if err != nil && err != io.EOF {
				log.Fatal(err)
				return
			} else {

			}

			log.Infof("send chunk %d", n)
			in := zfs.RecvRequest{Target: target, Length: int64(n), Chunk: buf[0:n]}
			if err == io.EOF {
				in.Done = true
			}
			stream.Send(&in)

			if err == io.EOF {
				break
			}
		}
	}()

	if err := send.Run(); err != nil {
		log.Fatal(err)
	}
}
