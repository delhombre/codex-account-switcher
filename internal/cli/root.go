// Package cli provides the command-line interface.
package cli

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/delhombre/cxa/internal/storage"
	"github.com/delhombre/cxa/internal/ui/styles"
	"github.com/delhombre/cxa/internal/ui/tui"
	"github.com/spf13/cobra"
)

var (
	repo    = storage.NewDirectoryRepository()
	version string
)

// Execute runs the CLI.
func Execute(v string) error {
	version = v
	return rootCmd.Execute()
}

var rootCmd = &cobra.Command{
	Use:   "cxa",
	Short: "Codex Account Switcher - Manage multiple Codex CLI accounts",
	Long: lipgloss.NewStyle().Foreground(styles.Primary).Render(`
   ___  _  __   _   
  / __|| | \ \ / /  _ \
 | (__ |_|  \ V /| (_) |
  \___|(_)  |_|  \___/

`) + "Manage multiple OpenAI Codex CLI accounts with ease.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// No args = launch TUI
		return tui.Run(repo)
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all saved accounts",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		accounts, err := repo.List()
		if err != nil {
			return err
		}

		current, _ := repo.Current()

		if len(accounts) == 0 {
			fmt.Println(styles.MutedStyle.Render("No accounts saved yet."))
			fmt.Println(styles.MutedStyle.Render("Save your current account with: cxa save <name>"))
			return nil
		}

		fmt.Println(styles.RenderTitle("Saved Accounts"))
		fmt.Println()

		for _, acc := range accounts {
			if acc.Name == current {
				fmt.Printf("  %s %s %s\n",
					styles.Bullet,
					styles.CurrentAccountStyle.Render(acc.Name),
					styles.MutedStyle.Render("(current)"),
				)
			} else {
				fmt.Printf("  %s %s\n",
					styles.Circle,
					acc.Name,
				)
			}
		}
		fmt.Println()

		return nil
	},
}

var switchCmd = &cobra.Command{
	Use:     "switch <name>",
	Short:   "Switch to a different account",
	Aliases: []string{"sw", "use"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		fmt.Printf("%s Switching to %s...\n",
			styles.Caret,
			styles.PrimaryStyle.Render(name),
		)

		if err := repo.Activate(name); err != nil {
			fmt.Println(styles.RenderError(err.Error()))
			return err
		}

		fmt.Println(styles.RenderSuccess(fmt.Sprintf("Switched to %s", name)))
		return nil
	},
}

var saveCmd = &cobra.Command{
	Use:   "save <name>",
	Short: "Save the current ~/.codex as an account",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		fmt.Printf("%s Saving current session as %s...\n",
			styles.Caret,
			styles.PrimaryStyle.Render(name),
		)

		if _, err := repo.Save(name); err != nil {
			fmt.Println(styles.RenderError(err.Error()))
			return err
		}

		fmt.Println(styles.RenderSuccess(fmt.Sprintf("Saved account: %s", name)))
		return nil
	},
}

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show the current active account",
	RunE: func(cmd *cobra.Command, args []string) error {
		current, err := repo.Current()
		if err != nil {
			return err
		}

		if current == "" {
			fmt.Println(styles.MutedStyle.Render("No active account tracked."))
			return nil
		}

		fmt.Printf("%s Current account: %s\n",
			styles.Bullet,
			styles.CurrentAccountStyle.Render(current),
		)
		return nil
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("cxa version %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(switchCmd)
	rootCmd.AddCommand(saveCmd)
	rootCmd.AddCommand(currentCmd)
	rootCmd.AddCommand(versionCmd)

	// Silence usage on errors
	rootCmd.SilenceUsage = true
}
