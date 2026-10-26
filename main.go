package main

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/jerr3n/thyme/cmd"
)

func RosemaryVersion(rosemaryExec string, options cmd.RosemaryOptions) (string, error) {
	var out bytes.Buffer
	if err := cmd.RunRosemary(context.Background(), rosemaryExec, options, []string{"version"}, cmd.RunIO{Stdout: &out}); err != nil {
		return "", fmt.Errorf("rosemary version: %w", err)
	}
	return strings.TrimSpace(out.String()), nil
}
func GetRosemaryPath() (string, error) {
	path, err := exec.LookPath("rosemary")
	if err != nil {
		return "", err
	}
	return path, nil
}
func Run(in string, ctx context.Context, options cmd.RosemaryOptions, rosemaryExec *string, io cmd.RunIO) error {
	strs := strings.Split(in, " ")
	if strs[0] == "rosemary" || strs[0][0] == '/' {
		strs = strs[1:]
		if strs[0] == "run" {
			strs = strs[1:]
		}
	}
	var err error
	executable := *rosemaryExec
	if rosemaryExec == nil { // not always false do not listen to your ide
		executable, err = GetRosemaryPath()
		if err != nil {
			return err
		}
	}
	err = cmd.RunRosemary(ctx, executable, options, strs, io)
	if err != nil {
		return err
	}
	return nil
}
