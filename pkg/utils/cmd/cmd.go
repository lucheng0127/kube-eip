package cmd

import (
	"errors"
	"os/exec"
	"strings"
)

type CmdMgr struct {
	Path string
}

func NewCmdMgr(name string) (*CmdMgr, error) {
	binPath, err := exec.LookPath(name)

	if err != nil {
		return nil, err
	}

	return &CmdMgr{binPath}, nil
}

func (c *CmdMgr) Execute(args ...string) (string, error) {
	var stderr strings.Builder
	var stdout strings.Builder
	cmd := exec.Cmd{
		Path:   c.Path,
		Args:   append([]string{c.Path}, args...),
		Stderr: &stderr,
		Stdout: &stdout,
	}

	if err := cmd.Run(); err != nil {
		return "", errors.New(stderr.String())
	}

	return stdout.String(), nil
}
