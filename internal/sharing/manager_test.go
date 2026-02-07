package sharing_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/delhombre/cxa/internal/sharing"
)

func TestManager_EnableDisable(t *testing.T) {
	tmpDir := t.TempDir()
	homeDir := filepath.Join(tmpDir, ".codex")
	
	// Create ~/.codex structure
	if err := os.MkdirAll(filepath.Join(homeDir, "sessions"), 0755); err != nil {
		t.Fatalf("failed to create sessions dir: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(homeDir, "sqlite"), 0755); err != nil {
		t.Fatalf("failed to create sqlite dir: %v", err)
	}

	os.Setenv("HOME", tmpDir)
	defer os.Unsetenv("HOME")

	manager := sharing.NewManager()

	// Initially disabled
	if manager.IsEnabled() {
		t.Error("sharing should be disabled initially")
	}

	// Enable
	if err := manager.Enable(false); err != nil {
		t.Fatalf("Enable failed: %v", err)
	}

	// Reload and check
	manager2 := sharing.NewManager()
	if err := manager2.LoadConfig(); err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if !manager2.IsEnabled() {
		t.Error("sharing should be enabled")
	}

	if manager2.GetMode() != sharing.ModeGlobal {
		t.Errorf("expected mode Global, got %s", manager2.GetMode())
	}

	// Disable
	if err := manager2.Disable(); err != nil {
		t.Fatalf("Disable failed: %v", err)
	}

	manager3 := sharing.NewManager()
	if err := manager3.LoadConfig(); err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if manager3.IsEnabled() {
		t.Error("sharing should be disabled after Disable()")
	}
}

func TestManager_SymlinkCreation(t *testing.T) {
	tmpDir := t.TempDir()
	homeDir := filepath.Join(tmpDir, ".codex")
	sharedDir := filepath.Join(tmpDir, "codex-data", "shared")
	
	// Create ~/.codex with sessions
	sessionsDir := filepath.Join(homeDir, "sessions")
	if err := os.MkdirAll(sessionsDir, 0755); err != nil {
		t.Fatalf("failed to create sessions dir: %v", err)
	}
	
	// Create a test file in sessions
	testFile := filepath.Join(sessionsDir, "test.json")
	if err := os.WriteFile(testFile, []byte(`{"test": true}`), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	os.Setenv("HOME", tmpDir)
	defer os.Unsetenv("HOME")

	manager := sharing.NewManager()
	if err := manager.Enable(false); err != nil {
		t.Fatalf("Enable failed: %v", err)
	}

	// Check that sessions is now a symlink
	info, err := os.Lstat(sessionsDir)
	if err != nil {
		t.Fatalf("failed to stat sessions: %v", err)
	}

	if info.Mode()&os.ModeSymlink == 0 {
		t.Error("sessions should be a symlink after Enable()")
	}

	// Check that the file was migrated
	migratedFile := filepath.Join(sharedDir, "sessions", "test.json")
	if _, err := os.Stat(migratedFile); os.IsNotExist(err) {
		t.Error("test file should have been migrated to shared location")
	}
}
