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

func RunRosemary(ctx context.Context, options RosemaryOptions, args []string, stdout *io.Writer, stderr *io.Writer) error {
	flags := RosemaryFlags(options)
	cmd := exec.CommandContext(ctx, "rosemary", slices.Concat(flags, args)...)
	cmd.Stdout = *stdout
	cmd.Stderr = *stderr
	return cmd.Run()
}
