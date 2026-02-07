// Package sharing manages session sharing between accounts.
package sharing

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/delhombre/cxa/pkg/codex"
)

// Mode represents the sharing mode.
type Mode string

const (
	ModeDisabled Mode = "disabled"
	ModeGlobal   Mode = "global"
	ModeGroup    Mode = "group"
)

// Config holds the sharing configuration.
type Config struct {
	Mode            Mode              `json:"mode"`
	IncludeSettings bool              `json:"include_settings"`
	Groups          map[string]string `json:"groups"` // account -> group mapping
}

// Manager handles session sharing between accounts.
type Manager struct {
	paths  *codex.Paths
	config *Config
}

// NewManager creates a new sharing manager.
func NewManager() *Manager {
	return &Manager{
		paths:  codex.NewPaths(),
		config: &Config{Mode: ModeDisabled},
	}
}

// LoadConfig loads the sharing configuration from disk.
func (m *Manager) LoadConfig() error {
	data, err := os.ReadFile(m.paths.SharingConfigFile())
	if err != nil {
		if os.IsNotExist(err) {
			m.config = &Config{Mode: ModeDisabled}
			return nil
		}
		return err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return err
	}

	m.config = &config
	return nil
}

// SaveConfig writes the sharing configuration to disk.
func (m *Manager) SaveConfig() error {
	if err := m.paths.EnsureDirs(); err != nil {
		return err
	}

	data, err := json.MarshalIndent(m.config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.paths.SharingConfigFile(), data, 0644)
}

// IsEnabled returns true if sharing is enabled.
func (m *Manager) IsEnabled() bool {
	return m.config.Mode == ModeGlobal || m.config.Mode == ModeGroup
}

// GetMode returns the current sharing mode.
func (m *Manager) GetMode() Mode {
	return m.config.Mode
}

// IncludesSettings returns whether settings are shared.
func (m *Manager) IncludesSettings() bool {
	return m.config.IncludeSettings
}

// Enable enables global sharing.
func (m *Manager) Enable(includeSettings bool) error {
	m.config.Mode = ModeGlobal
	m.config.IncludeSettings = includeSettings

	// Create shared directory
	if err := os.MkdirAll(m.paths.SharedDir, 0755); err != nil {
		return err
	}

	// Setup symlinks
	if err := m.SetupSymlinks(); err != nil {
		return err
	}

	return m.SaveConfig()
}

// Disable disables sharing and copies data locally.
func (m *Manager) Disable() error {
	// First, copy shared data back to local
	if err := m.RemoveSymlinks(); err != nil {
		return err
	}

	m.config.Mode = ModeDisabled
	m.config.IncludeSettings = false

	return m.SaveConfig()
}

// SetupSymlinks creates symlinks from ~/.codex to the shared location.
func (m *Manager) SetupSymlinks() error {
	if !m.IsEnabled() {
		return nil
	}

	targetDir := m.getShareTarget("")
	if targetDir == "" {
		return nil
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return err
	}

	// Setup symlinks for shareable items
	for _, item := range codex.ShareableItems {
		if err := m.setupSymlink(item, targetDir); err != nil {
			return fmt.Errorf("failed to setup symlink for %s: %w", item, err)
		}
	}

	// Optionally setup symlinks for settings
	if m.config.IncludeSettings {
		for _, item := range codex.OptionalShareableItems {
			if err := m.setupSymlink(item, targetDir); err != nil {
				return fmt.Errorf("failed to setup symlink for %s: %w", item, err)
			}
		}
	}

	return nil
}

func (m *Manager) setupSymlink(item, targetDir string) error {
	src := filepath.Join(m.paths.Home, item)
	dest := filepath.Join(targetDir, item)

	// Check if already a symlink to the correct location
	if link, err := os.Readlink(src); err == nil {
		if link == dest {
			return nil // Already correct
		}
		// Wrong symlink, remove it
		os.Remove(src)
	}

	// If source exists and is not a symlink, migrate it
	if info, err := os.Lstat(src); err == nil && info.Mode()&os.ModeSymlink == 0 {
		if _, err := os.Stat(dest); os.IsNotExist(err) {
			// Migrate to shared location
			if err := os.Rename(src, dest); err != nil {
				return err
			}
		} else {
			// Both exist, remove local copy
			os.RemoveAll(src)
		}
	}

	// Ensure target exists
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		// Create empty target
		if filepath.Ext(item) != "" {
			// File
			if err := os.WriteFile(dest, []byte{}, 0644); err != nil {
				return err
			}
		} else {
			// Directory
			if err := os.MkdirAll(dest, 0755); err != nil {
				return err
			}
		}
	}

	// Create symlink
	return os.Symlink(dest, src)
}

// RemoveSymlinks replaces symlinks with copies of the shared data.
func (m *Manager) RemoveSymlinks() error {
	allItems := append(codex.ShareableItems, codex.OptionalShareableItems...)

	for _, item := range allItems {
		src := filepath.Join(m.paths.Home, item)

		// Check if it's a symlink
		link, err := os.Readlink(src)
		if err != nil {
			continue // Not a symlink
		}

		// Remove the symlink
		os.Remove(src)

		// Copy the target data back
		if _, err := os.Stat(link); err == nil {
			if err := copyPath(link, src); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *Manager) getShareTarget(account string) string {
	switch m.config.Mode {
	case ModeGlobal:
		return m.paths.SharedDir
	case ModeGroup:
		if group, ok := m.config.Groups[account]; ok {
			return filepath.Join(m.paths.GroupsDir, group)
		}
		return ""
	default:
		return ""
	}
}

// Status returns the current sharing status.
func (m *Manager) Status() (mode Mode, sharedDir string, symlinks map[string]string) {
	mode = m.config.Mode
	symlinks = make(map[string]string)

	if mode == ModeGlobal {
		sharedDir = m.paths.SharedDir
	}

	allItems := append(codex.ShareableItems, codex.OptionalShareableItems...)
	for _, item := range allItems {
		src := filepath.Join(m.paths.Home, item)
		if link, err := os.Readlink(src); err == nil {
			symlinks[item] = link
		} else if _, err := os.Stat(src); err == nil {
			symlinks[item] = "(local)"
		} else {
			symlinks[item] = "(missing)"
		}
	}

	return
}

// copyPath copies a file or directory.
func copyPath(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return copyDir(src, dst)
	}
	return copyFile(src, dst)
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(src, path)
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		return copyFile(path, dstPath)
	})
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	info, _ := os.Stat(src)
	return os.WriteFile(dst, data, info.Mode())
}
