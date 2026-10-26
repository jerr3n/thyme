package prefix

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/jerr3n/thyme/cmd"
	"github.com/jerr3n/thyme/util"
	"go.yaml.in/yaml/v4"
)

type StateYAML struct {
	Version     string `yaml:"version"`
	CreatedWith struct {
		RuntimeVersion string `yaml:"runtime_version"`
	} `yaml:"created_with"`
	Base struct {
		Project      string   `yaml:"project"`
		Arch         string   `yaml:"arch"`
		Sequence     string   `yaml:"sequence"`
		Tag          string   `yaml:"tag"`
		Tree         string   `yaml:"tree"`
		ManifestHash string   `yaml:"manifest_sha256"`
		ArtifactHash string   `yaml:"artifact_sha256"`
		Trust        string   `yaml:"trust"`
		Policy       string   `yaml:"policy"`
		Materialised []string `yaml:"materialised"`
	} `yaml:"base"`
	AcceptedSequence *string `yaml:"accepted_sequence"`
}

type Prefix struct {
	Dir   string
	State StateYAML
}

type SystemFiles string

const (
	SystemFilesRepair    SystemFiles = "repair"
	SystemFilesWarn      SystemFiles = "warn"
	SystemFilesImmutable SystemFiles = "immutable"
)

type CreateOptions struct {
	Arch        string
	Spec        string
	Bundle      string
	SystemFiles SystemFiles
	NoBase      bool
}

func DefaultDir() (string, error) {
	if data := os.Getenv("XDG_DATA_HOME"); data != "" {
		return filepath.Join(data, "rosemary/prefixes/default"), nil
	}
	home := os.Getenv("HOME")
	if home == "" {
		return "", errors.New("no $HOME, so there is no default prefix")
	}
	return filepath.Join(home, ".local/share/rosemary/prefixes/default"), nil
}

func (p Prefix) Path() (string, error) {
	if p.Dir == "" {
		return DefaultDir()
	}
	return util.ExpandHome(p.Dir)
}

func (p Prefix) statePath() (string, error) {
	dir, err := p.Path()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, ".rosemary", "state.yaml"), nil
}

func (p Prefix) Exists() bool {
	dest, err := p.statePath()
	if err != nil {
		return false
	}
	_, err = os.Stat(dest)
	return err == nil
}

func (p *Prefix) ResolveYAML() error {
	if !p.Exists() {
		return fs.ErrNotExist
	}
	dest, err := p.statePath()
	if err != nil {
		return err
	}
	f, err := os.ReadFile(dest)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(f, &p.State)
}

func (p Prefix) run(ctx context.Context, rosemaryExec, action string, args []string, runIO cmd.RunIO) error {
	dir, err := p.Path()
	if err != nil {
		return err
	}
	all := append([]string{"prefix", action, dir}, args...)
	if err := cmd.Exec(ctx, rosemaryExec, all, runIO); err != nil {
		return fmt.Errorf("rosemary prefix %s: %w", action, err)
	}
	return nil
}

func (p Prefix) Create(ctx context.Context, rosemaryExec string, options CreateOptions, runIO cmd.RunIO) error {
	var args []string
	if options.Arch != "" {
		args = append(args, "--arch", options.Arch)
	}
	if options.Spec != "" {
		args = append(args, "--spec", options.Spec)
	}
	if options.Bundle != "" {
		args = append(args, "--bundle", options.Bundle)
	}
	if options.SystemFiles != "" {
		args = append(args, "--system-files", string(options.SystemFiles))
	}
	if options.NoBase {
		args = append(args, "--no-base")
	}
	return p.run(ctx, rosemaryExec, "create", args, runIO)
}

func (p Prefix) Info(ctx context.Context, rosemaryExec string) (string, error) {
	var out, errOut bytes.Buffer
	if err := p.run(ctx, rosemaryExec, "info", nil, cmd.RunIO{Stdout: &out, Stderr: &errOut}); err != nil {
		if msg := strings.TrimSpace(errOut.String()); msg != "" {
			return "", fmt.Errorf("%w: %s", err, msg)
		}
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func (p Prefix) Verify(ctx context.Context, rosemaryExec string, runIO cmd.RunIO) error {
	return p.run(ctx, rosemaryExec, "verify", nil, runIO)
}

func (p Prefix) Repair(ctx context.Context, rosemaryExec string, dryRun bool, runIO cmd.RunIO) error {
	return p.run(ctx, rosemaryExec, "repair", dryRunArgs(dryRun), runIO)
}

func (p Prefix) ResetSystem(ctx context.Context, rosemaryExec string, dryRun bool, runIO cmd.RunIO) error {
	return p.run(ctx, rosemaryExec, "reset-system", dryRunArgs(dryRun), runIO)
}

func (p Prefix) Update(ctx context.Context, rosemaryExec string, allowDowngrade bool, runIO cmd.RunIO) error {
	var args []string
	if allowDowngrade {
		args = append(args, "--allow-downgrade")
	}
	return p.run(ctx, rosemaryExec, "update", args, runIO)
}

func dryRunArgs(dryRun bool) []string {
	if dryRun {
		return []string{"--dry-run"}
	}
	return nil
}

func CacheList(ctx context.Context, rosemaryExec string, runIO cmd.RunIO) error {
	return cmd.Exec(ctx, rosemaryExec, []string{"prefix", "cache", "list"}, runIO)
}

func CacheGC(ctx context.Context, rosemaryExec string, runIO cmd.RunIO) error {
	return cmd.Exec(ctx, rosemaryExec, []string{"prefix", "cache", "gc"}, runIO)
}
