// Package storage provides account storage implementations.
package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/delhombre/cxa/internal/account"
	"github.com/delhombre/cxa/internal/sharing"
	"github.com/delhombre/cxa/pkg/codex"
)

// DirectoryRepository implements account.Repository using directories.
// This is much faster than zip-based storage.
type DirectoryRepository struct {
	paths *codex.Paths
}

// NewDirectoryRepository creates a new directory-based repository.
func NewDirectoryRepository() *DirectoryRepository {
	return &DirectoryRepository{
		paths: codex.NewPaths(),
	}
}

// List returns all saved accounts.
func (r *DirectoryRepository) List() ([]*account.Account, error) {
	accountsDir := r.paths.AccountsDir()
	if err := r.paths.EnsureDirs(); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(accountsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []*account.Account{}, nil
		}
		return nil, err
	}

	var accounts []*account.Account
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		acc, err := r.Get(entry.Name())
		if err != nil {
			continue // Skip invalid accounts
		}
		accounts = append(accounts, acc)
	}

	return accounts, nil
}

// Get retrieves an account by name.
func (r *DirectoryRepository) Get(name string) (*account.Account, error) {
	accountPath := r.paths.AccountPath(name)
	metaPath := filepath.Join(accountPath, ".account.json")

	data, err := os.ReadFile(metaPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Account exists but no metadata, create basic account
			info, statErr := os.Stat(accountPath)
			if statErr != nil {
				return nil, fmt.Errorf("account '%s' not found", name)
			}
			return &account.Account{
				Name:      name,
				CreatedAt: info.ModTime(),
				UpdatedAt: info.ModTime(),
			}, nil
		}
		return nil, err
	}

	var acc account.Account
	if err := json.Unmarshal(data, &acc); err != nil {
		return nil, err
	}

	return &acc, nil
}

// Save stores the current ~/.codex as the given account.
func (r *DirectoryRepository) Save(name string) (*account.Account, error) {
	if !r.paths.CodexExists() {
		return nil, errors.New("~/.codex not found - please login first with 'codex login'")
	}

	if err := r.paths.EnsureDirs(); err != nil {
		return nil, err
	}

	accountPath := r.paths.AccountPath(name)

	// Remove existing account data if exists
	_ = os.RemoveAll(accountPath)

	// Copy ~/.codex to account directory
	if err := copyDir(r.paths.Home, accountPath); err != nil {
		return nil, fmt.Errorf("failed to save account: %w", err)
	}

	// Create account metadata
	acc := &account.Account{
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Try to read email from auth.json
	authPath := filepath.Join(accountPath, "auth.json")
	if authData, err := os.ReadFile(authPath); err == nil {
		var authInfo struct {
			Tokens struct {
				IDToken string `json:"id_token"`
			} `json:"tokens"`
		}
		if json.Unmarshal(authData, &authInfo) == nil {
			// Could parse JWT to get email, simplified for now
		}
	}

	// Save metadata
	metaPath := filepath.Join(accountPath, ".account.json")
	metaData, _ := json.MarshalIndent(acc, "", "  ")
	if err := os.WriteFile(metaPath, metaData, 0644); err != nil {
		return nil, err
	}

	// Update current account state
	if err := r.saveState(name); err != nil {
		return nil, err
	}

	return acc, nil
}

// Delete removes an account.
func (r *DirectoryRepository) Delete(name string) error {
	accountPath := r.paths.AccountPath(name)
	if _, err := os.Stat(accountPath); os.IsNotExist(err) {
		return fmt.Errorf("account '%s' not found", name)
	}
	return os.RemoveAll(accountPath)
}

// Activate switches to the given account.
func (r *DirectoryRepository) Activate(name string) error {
	accountPath := r.paths.AccountPath(name)
	if _, err := os.Stat(accountPath); os.IsNotExist(err) {
		return fmt.Errorf("account '%s' not found", name)
	}

	// Get current account to save it first
	current, _ := r.Current()
	if current != "" && current != name {
		// Save current state before switching
		if r.paths.CodexExists() {
			if _, err := r.Save(current); err != nil {
				return fmt.Errorf("failed to save current account: %w", err)
			}
		}
	}

	// Remove current ~/.codex
	if err := os.RemoveAll(r.paths.Home); err != nil {
		return fmt.Errorf("failed to clear ~/.codex: %w", err)
	}

	// Copy account to ~/.codex
	if err := copyDir(accountPath, r.paths.Home); err != nil {
		return fmt.Errorf("failed to activate account: %w", err)
	}

	// Re-setup sharing symlinks if enabled
	shareManager := sharing.NewManager()
	if err := shareManager.LoadConfig(); err == nil && shareManager.IsEnabled() {
		_ = shareManager.SetupSymlinks()
	}

	// Update state
	if err := r.saveState(name); err != nil {
		return err
	}

	return nil
}

// Current returns the currently active account name.
func (r *DirectoryRepository) Current() (string, error) {
	state, err := r.loadState()
	if err != nil {
		return "", nil
	}
	return state.Current, nil
}

// State tracks the current and previous accounts.
type State struct {
	Current  string `json:"current"`
	Previous string `json:"previous"`
}

func (r *DirectoryRepository) loadState() (*State, error) {
	data, err := os.ReadFile(r.paths.StateFile())
	if err != nil {
		return &State{}, nil
	}
	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return &State{}, nil
	}
	return &state, nil
}

func (r *DirectoryRepository) saveState(current string) error {
	state, _ := r.loadState()
	state.Previous = state.Current
	state.Current = current

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	if err := r.paths.EnsureDirs(); err != nil {
		return err
	}

	return os.WriteFile(r.paths.StateFile(), data, 0644)
}

// copyDir recursively copies a directory.
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate destination path
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relPath)

		// Handle symlinks
		if info.Mode()&os.ModeSymlink != 0 {
			link, err := os.Readlink(path)
			if err != nil {
				return err
			}
			return os.Symlink(link, dstPath)
		}

		// Handle directories
		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		// Copy file
		return copyFile(path, dstPath)
	})
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
