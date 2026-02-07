package storage_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/delhombre/cxa/internal/storage"
)

func TestDirectoryRepository_SaveAndList(t *testing.T) {
	// Setup temp directories
	tmpDir := t.TempDir()
	homeDir := filepath.Join(tmpDir, ".codex")
	dataDir := filepath.Join(tmpDir, "codex-data")
	stateDir := filepath.Join(tmpDir, ".codex-switch")

	// Create fake ~/.codex
	if err := os.MkdirAll(homeDir, 0755); err != nil {
		t.Fatalf("failed to create home dir: %v", err)
	}

	// Create a test file
	authFile := filepath.Join(homeDir, "auth.json")
	if err := os.WriteFile(authFile, []byte(`{"test": true}`), 0644); err != nil {
		t.Fatalf("failed to write auth file: %v", err)
	}

	// Override paths (would need to inject paths in real implementation)
	os.Setenv("HOME", tmpDir)
	defer os.Unsetenv("HOME")

	// Create repository
	repo := storage.NewDirectoryRepository()

	// Save account
	acc, err := repo.Save("test-account")
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	if acc.Name != "test-account" {
		t.Errorf("expected name 'test-account', got '%s'", acc.Name)
	}

	// List accounts
	accounts, err := repo.List()
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(accounts) != 1 {
		t.Errorf("expected 1 account, got %d", len(accounts))
	}

	// Verify file was copied
	savedAuth := filepath.Join(dataDir, "accounts", "test-account", "auth.json")
	if _, err := os.Stat(savedAuth); os.IsNotExist(err) {
		t.Error("auth.json was not saved to account directory")
	}

	// Check current
	current, _ := repo.Current()
	if current != "test-account" {
		t.Errorf("expected current 'test-account', got '%s'", current)
	}

	// Clean up
	_ = dataDir
	_ = stateDir
}

func TestDirectoryRepository_Activate(t *testing.T) {
	tmpDir := t.TempDir()
	homeDir := filepath.Join(tmpDir, ".codex")

	// Create fake ~/.codex with account1 content
	if err := os.MkdirAll(homeDir, 0755); err != nil {
		t.Fatalf("failed to create home dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(homeDir, "marker.txt"), []byte("account1"), 0644); err != nil {
		t.Fatalf("failed to write marker: %v", err)
	}

	os.Setenv("HOME", tmpDir)
	defer os.Unsetenv("HOME")

	repo := storage.NewDirectoryRepository()

	// Save as account1
	if _, err := repo.Save("account1"); err != nil {
		t.Fatalf("Save account1 failed: %v", err)
	}

	// Change marker and save as account2
	if err := os.WriteFile(filepath.Join(homeDir, "marker.txt"), []byte("account2"), 0644); err != nil {
		t.Fatalf("failed to update marker: %v", err)
	}
	if _, err := repo.Save("account2"); err != nil {
		t.Fatalf("Save account2 failed: %v", err)
	}

	// Switch back to account1
	if err := repo.Activate("account1"); err != nil {
		t.Fatalf("Activate failed: %v", err)
	}

	// Verify marker is back to account1
	data, err := os.ReadFile(filepath.Join(homeDir, "marker.txt"))
	if err != nil {
		t.Fatalf("failed to read marker: %v", err)
	}

	if string(data) != "account1" {
		t.Errorf("expected marker 'account1', got '%s'", string(data))
	}

	current, _ := repo.Current()
	if current != "account1" {
		t.Errorf("expected current 'account1', got '%s'", current)
	}
}

func TestDirectoryRepository_Delete(t *testing.T) {
	tmpDir := t.TempDir()
	homeDir := filepath.Join(tmpDir, ".codex")

	if err := os.MkdirAll(homeDir, 0755); err != nil {
		t.Fatalf("failed to create home dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(homeDir, "test.txt"), []byte("test"), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	os.Setenv("HOME", tmpDir)
	defer os.Unsetenv("HOME")

	repo := storage.NewDirectoryRepository()

	// Save account
	if _, err := repo.Save("to-delete"); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Delete account
	if err := repo.Delete("to-delete"); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify account is gone
	accounts, _ := repo.List()
	for _, acc := range accounts {
		if acc.Name == "to-delete" {
			t.Error("account 'to-delete' should have been deleted")
		}
	}
}
