//go:build unix

package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"syscall"

	"github.com/jerr3n/thyme/cmd"
	"github.com/jerr3n/thyme/prefix"
)

type Server struct {
	Prefix prefix.Prefix
}

func (s Server) env() ([]string, error) {
	dir, err := s.Prefix.Path()
	if err != nil {
		return nil, err
	}
	return []string{"ROSEMARY_PREFIX=" + dir}, nil
}

func (s Server) Start(ctx context.Context, rosemaryExec string) error {
	env, err := s.env()
	if err != nil {
		return err
	}
	if err := cmd.Exec(ctx, rosemaryExec, []string{"server"}, cmd.RunIO{Env: env}); err != nil {
		return fmt.Errorf("rosemary server: %w", err)
	}
	return nil
}

func (s Server) Foreground(ctx context.Context, rosemaryExec string, runIO cmd.RunIO) error {
	env, err := s.env()
	if err != nil {
		return err
	}
	runIO.Env = append(runIO.Env, env...)
	if err := cmd.Exec(ctx, rosemaryExec, []string{"server", "--foreground"}, runIO); err != nil {
		return fmt.Errorf("rosemary server --foreground: %w", err)
	}
	return nil
}

func (s Server) Stop(ctx context.Context, rosemaryExec string, runIO cmd.RunIO) error {
	env, err := s.env()
	if err != nil {
		return err
	}
	runIO.Env = append(runIO.Env, env...)
	if err := cmd.Exec(ctx, rosemaryExec, []string{"server", "--stop"}, runIO); err != nil {
		return fmt.Errorf("rosemary server --stop: %w", err)
	}
	return nil
}

func (s Server) SocketPath() (string, error) {
	dir, err := s.Prefix.Path()
	if err != nil {
		return "", err
	}
	leaf := "server-host"
	if info, err := os.Stat(dir); err == nil {
		if st, ok := info.Sys().(*syscall.Stat_t); ok {
			leaf = fmt.Sprintf("server-%x-%x", uint64(st.Dev), uint64(st.Ino))
		}
	}
	return filepath.Join(runtimeDir(), leaf, "socket"), nil
}

func (s Server) Running() bool {
	socket, err := s.SocketPath()
	if err != nil {
		return false
	}
	conn, err := net.Dial("unix", socket)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func runtimeDir() string {
	if runtime := os.Getenv("XDG_RUNTIME_DIR"); runtime != "" {
		return filepath.Join(runtime, "rosemary")
	}
	return fmt.Sprintf("/tmp/rosemary-%d", os.Getuid())
}
