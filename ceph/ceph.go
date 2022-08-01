package ceph

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

type CephCmd struct {
	ceph string
	rbd  string
}

func NewCeph() (*CephCmd, error) {
	cmd := &CephCmd{}

	return cmd, nil
}

func shellExec(ctx context.Context, shell string, envs map[string]string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "/bin/bash", shell)
	for k, v := range envs {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	return bytes.TrimSuffix(out, []byte("\n")), nil
}
