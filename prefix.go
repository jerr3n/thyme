package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

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

func (p Prefix) Exists() bool {
	dest, err := util.ExpandHome(p.Dir)
	if err != nil {
		return false
	}
	_, err = os.Stat(dest)
	if errors.Is(err, fs.ErrNotExist) {
		return false
	} else if err != nil {
		return false
	}
	dest = fmt.Sprintf("%s/state.yaml")
	_, err = os.Stat(dest)
	if errors.Is(err, fs.ErrNotExist) {
		return false
	} else if err != nil {
		return false
	}
	return true
}

func (p *Prefix) ResolveYAML() error {
	if !p.Exists() {
		return fs.ErrNotExist
	}
	home, err := util.ExpandHome(p.Dir)
	if err != nil {
		return err
	}
	dest := fmt.Sprintf("%s/state.yaml", home)
	f, err := os.ReadFile(dest)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(f, p.State)
	if err != nil {
		return err
	}
	return nil
}
