package cmd

import (
	"context"
	"io"
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
	Stdout io.Writer
	Stderr io.Writer
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

func RunRosemary(ctx context.Context, rosemaryExec string, options RosemaryOptions, args []string, runIO RunIO) error {
	flags := RosemaryFlags(options)
	cmd := exec.CommandContext(ctx, rosemaryExec, slices.Concat(flags, args)...)
	cmd.Stdout = runIO.Stdout
	cmd.Stderr = runIO.Stdout
	return cmd.Run()
}
