// Package account provides account management for Codex CLI.
package account

import (
	"time"
)

// Account represents a Codex CLI account.
type Account struct {
	Name      string    `json:"name"`
	Email     string    `json:"email,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewAccount creates a new account with the given name.
func NewAccount(name string) *Account {
	now := time.Now()
	return &Account{
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Repository defines the interface for account storage.
type Repository interface {
	// List returns all saved accounts.
	List() ([]*Account, error)

	// Get retrieves an account by name.
	Get(name string) (*Account, error)

	// Save stores the current ~/.codex as the given account.
	Save(name string) (*Account, error)

	// Delete removes an account.
	Delete(name string) error

	// Activate switches to the given account.
	Activate(name string) error

	// Current returns the currently active account name.
	Current() (string, error)
}
