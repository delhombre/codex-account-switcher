// Package codex provides utilities for working with Codex CLI paths.
package codex

import (
	"os"
	"path/filepath"
)

// Paths contains all relevant Codex paths.
type Paths struct {
	Home      string // ~/.codex
	DataDir   string // ~/codex-data (account storage)
	StateDir  string // ~/.codex-switch (state tracking)
	SharedDir string // ~/codex-data/shared
	GroupsDir string // ~/codex-data/groups
}

// ShareableItems are the items that can be shared between accounts.
var ShareableItems = []string{
	"sessions",
	"sqlite",
	"history.jsonl",
	".codex-global-state.json",
}

// AccountSpecificItems are items that remain per-account.
var AccountSpecificItems = []string{
	"auth.json",
	"license.secret",
}

// OptionalShareableItems can optionally be shared.
var OptionalShareableItems = []string{
	"config.toml",
	"settings.json",
}

// NewPaths creates a new Paths instance with default locations.
func NewPaths() *Paths {
	home, _ := os.UserHomeDir()
	return &Paths{
		Home:      filepath.Join(home, ".codex"),
		DataDir:   filepath.Join(home, "codex-data"),
		StateDir:  filepath.Join(home, ".codex-switch"),
		SharedDir: filepath.Join(home, "codex-data", "shared"),
		GroupsDir: filepath.Join(home, "codex-data", "groups"),
	}
}

// AccountsDir returns the path to the accounts directory.
func (p *Paths) AccountsDir() string {
	return filepath.Join(p.DataDir, "accounts")
}

// AccountPath returns the path for a specific account.
func (p *Paths) AccountPath(name string) string {
	return filepath.Join(p.AccountsDir(), name)
}

// StateFile returns the path to the state file.
func (p *Paths) StateFile() string {
	return filepath.Join(p.StateDir, "state.json")
}

// SharingConfigFile returns the path to the sharing config.
func (p *Paths) SharingConfigFile() string {
	return filepath.Join(p.StateDir, "sharing.json")
}

// EnsureDirs creates all necessary directories.
func (p *Paths) EnsureDirs() error {
	dirs := []string{
		p.DataDir,
		p.StateDir,
		p.AccountsDir(),
	}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}

// CodexExists checks if ~/.codex exists.
func (p *Paths) CodexExists() bool {
	_, err := os.Stat(p.Home)
	return err == nil
}
