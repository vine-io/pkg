package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/lack-io/vine"
	log "github.com/lack-io/vine/service/logger"
	zfs "github.com/lack-io/vine/testdata/zfs/proto"
)

type storage struct {
}

func (s storage) Recv(bctx context.Context, stream zfs.Storage_RecvStream) error {

	// 首次的数据中保存快照信息，确定传输是否正常
	b, err := stream.Recv()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(bctx)
	defer cancel()
	// 根据快照信息构建命令行
	shell := fmt.Sprintf("zfs receive -Fu %s", b.Target)
	log.Infof(shell)
	// 启动 zfs 接收命令, 并通过 pipe 暴露输出流
	recv := exec.CommandContext(ctx, "/bin/sh", "-c", shell)
	recv.Stderr = os.Stderr
	wd, err := recv.StdinPipe()
	if err != nil {
		return err
	}

	// sum 来统计传输数据
	var (
		sum   = 0
		errCh = make(chan error, 1)
		done  = make(chan struct{})
	)
	// 启动一个 goroutine 来接收数据
	go func() {
		for {
			b, _ = stream.Recv()
			//if err != nil && err != io.EOF {
			//	log.Error(err)
			//	return err
			//} else {
			//
			//}

			if b != nil {
				n, e1 := wd.Write(b.Chunk[0:b.Length])
				if e1 != nil && e1 != io.EOF {
					log.Fatal(err)
					return
				}
				if b.Done {
					log.Infof("Done!!")
					break
				}
				sum += n
				//log.Infof("write chunk %d", n)
			}
		}
		log.Infof("recv %dMB", sum/1024/1024)
	}()

	// 启动命令行子进程, 等待接收完成
	if err := recv.Run(); err != nil {
		log.Fatal(err)
	}

	return err
}

var _ zfs.StorageHandler = (*storage)(nil)

func main() {
	srv := vine.NewService(
		vine.Name("go.vine.zfs"),
		vine.Address("127.0.0.1:2333"),
	)

	srv.Init()

	zfs.RegisterStorageHandler(srv.Server(), new(storage))

	srv.Run()
}
