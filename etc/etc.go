package etc

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/jerr3n/thyme/cmd"
)

func GetRosemaryPath() (string, error) {
	path, err := exec.LookPath("rosemary")
	if err != nil {
		return "", err
	}
	return path, nil
}

func RosemaryVersion(ctx context.Context, rosemaryExec string, options cmd.RosemaryOptions) (string, error) {
	var out bytes.Buffer
	if err := cmd.RunRosemary(ctx, rosemaryExec, options, []string{"--version"}, cmd.RunIO{Stdout: &out}); err != nil {
		return "", fmt.Errorf("rosemary --version: %w", err)
	}
	return strings.TrimSpace(out.String()), nil
}
