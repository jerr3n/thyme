package cmd

import (
	"context"
	"io"
	"os"
	"os/exec"
	"reflect"
	"slices"
)

type RosemaryOptions struct {
	Slow        bool `rosemary:"slow"`
	Native      bool `rosemary:"native"`
	Emulated    bool `rosemary:"emulated"`
	ObjcDbg     bool `rosemary:"objc-debug"`
	NewInstance bool `rosemary:"new"`
}

type RunIO struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	Env    []string
}

func RosemaryFlags(options RosemaryOptions) []string {
	flags := []string{}

	v := reflect.ValueOf(options)
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		name, ok := t.Field(i).Tag.Lookup("rosemary")
		if !ok || name == "" {
			continue
		}
		if v.Field(i).Kind() == reflect.Bool && v.Field(i).Bool() {
			flags = append(flags, "--"+name)
		}
	}

	return flags
}

func Exec(ctx context.Context, rosemaryExec string, args []string, runIO RunIO) error {
	cmd := exec.CommandContext(ctx, rosemaryExec, args...)
	cmd.Stdin = runIO.Stdin
	cmd.Stdout = runIO.Stdout
	cmd.Stderr = runIO.Stderr
	if len(runIO.Env) > 0 {
		cmd.Env = append(os.Environ(), runIO.Env...)
	}
	return cmd.Run()
}

func RunRosemary(ctx context.Context, rosemaryExec string, options RosemaryOptions, args []string, runIO RunIO) error {
	return Exec(ctx, rosemaryExec, slices.Concat(RosemaryFlags(options), args), runIO)
}
