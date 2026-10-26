package thyme

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jerr3n/thyme/cmd"
	"github.com/jerr3n/thyme/util"
)

func RosemaryVersion(rosemaryExec string, options cmd.RosemaryOptions) (string, error) {
	var out bytes.Buffer
	if err := cmd.RunRosemary(context.Background(), rosemaryExec, options, []string{"--version"}, cmd.RunIO{Stdout: &out}); err != nil {
		return "", fmt.Errorf("rosemary --version: %w", err)
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
	return RunArgs(strings.Fields(in), ctx, options, rosemaryExec, io)
}

func RunArgs(args []string, ctx context.Context, options cmd.RosemaryOptions, rosemaryExec *string, io cmd.RunIO) error {
	strs := args
	if len(strs) > 0 && filepath.Base(strs[0]) == "rosemary" {
		strs = strs[1:]
		if len(strs) > 0 && strs[0] == "run" {
			strs = strs[1:]
		}
	}
	if len(strs) == 0 {
		return errors.New("nothing to run")
	}
	expanded := make([]string, len(strs))
	for i, s := range strs {
		e, err := util.ExpandHome(s)
		if err != nil {
			return err
		}
		expanded[i] = e
	}
	strs = expanded
	var err error
	var executable string
	if rosemaryExec != nil {
		executable = *rosemaryExec
	} else {
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
