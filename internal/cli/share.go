package cli

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/delhombre/cxa/internal/sharing"
	"github.com/delhombre/cxa/internal/ui/styles"
	"github.com/spf13/cobra"
)

var shareCmd = &cobra.Command{
	Use:   "share",
	Short: "Manage session sharing between accounts",
	Long:  "Share sessions, threads, and history between accounts while keeping authentication separate.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var shareEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable session sharing",
	RunE: func(cmd *cobra.Command, args []string) error {
		manager := sharing.NewManager()
		if err := manager.LoadConfig(); err != nil {
			return err
		}

		if manager.IsEnabled() {
			fmt.Println(styles.RenderWarning(fmt.Sprintf("Sharing is already enabled (mode: %s)", manager.GetMode())))
			return nil
		}

		fmt.Println()
		fmt.Println(styles.RenderTitle("Session Sharing Setup"))
		fmt.Println()
		fmt.Println("This will share sessions, threads, and history between all your accounts.")
		fmt.Println(styles.MutedStyle.Render("Authentication (auth.json) remains private to each account."))
		fmt.Println()

		// Interactive form
		var includeSettings bool
		var confirmMigrate bool

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Also share settings (config.toml, settings.json)?").
					Value(&includeSettings),
				huh.NewConfirm().
					Title("Migrate existing sessions to shared location?").
					Description("Recommended: keeps your current sessions accessible").
					Value(&confirmMigrate),
			),
		)

		if err := form.Run(); err != nil {
			return err
		}

		fmt.Printf("%s Enabling session sharing...\n", styles.Caret)

		if err := manager.Enable(includeSettings); err != nil {
			fmt.Println(styles.RenderError(err.Error()))
			return err
		}

		fmt.Println(styles.RenderSuccess("Session sharing enabled (global mode)"))
		fmt.Println(styles.MutedStyle.Render("All accounts will now share sessions, threads, and history."))

		return nil
	},
}

var shareDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable session sharing",
	RunE: func(cmd *cobra.Command, args []string) error {
		manager := sharing.NewManager()
		if err := manager.LoadConfig(); err != nil {
			return err
		}

		if !manager.IsEnabled() {
			fmt.Println(styles.MutedStyle.Render("Sharing is already disabled."))
			return nil
		}

		fmt.Println()
		fmt.Println("Disabling sharing will copy current shared data to your account's local storage.")

		var confirm bool
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Continue?").
					Value(&confirm),
			),
		)

		if err := form.Run(); err != nil {
			return err
		}

		if !confirm {
			fmt.Println(styles.MutedStyle.Render("Cancelled."))
			return nil
		}

		if err := manager.Disable(); err != nil {
			fmt.Println(styles.RenderError(err.Error()))
			return err
		}

		fmt.Println(styles.RenderSuccess("Session sharing disabled"))
		fmt.Println(styles.MutedStyle.Render("Your sessions have been copied locally."))

		return nil
	},
}

var shareStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show sharing configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		manager := sharing.NewManager()
		if err := manager.LoadConfig(); err != nil {
			return err
		}

		mode, sharedDir, symlinks := manager.Status()

		fmt.Println()
		fmt.Println(styles.RenderTitle("Sharing Status"))
		fmt.Println()

		// Mode
		modeStr := string(mode)
		if mode == sharing.ModeDisabled {
			modeStr = styles.MutedStyle.Render(modeStr)
		} else {
			modeStr = styles.SuccessStyle.Render(modeStr)
		}
		fmt.Printf("  Mode: %s\n", modeStr)

		if sharedDir != "" {
			fmt.Printf("  Location: %s\n", styles.MutedStyle.Render(sharedDir))
		}

		fmt.Println()
		fmt.Println("  Symlinks:")
		for item, target := range symlinks {
			var status string
			switch target {
			case "(local)":
				status = fmt.Sprintf("  %s %s %s", styles.Circle, item, styles.MutedStyle.Render(target))
			case "(missing)":
				status = fmt.Sprintf("  %s %s %s", styles.CrossMark, item, styles.MutedStyle.Render(target))
			default:
				status = fmt.Sprintf("  %s %s %s %s", styles.CheckMark, item, styles.Arrow, styles.MutedStyle.Render(target))
			}
			fmt.Println(status)
		}
		fmt.Println()

		return nil
	},
}

func init() {
	shareCmd.AddCommand(shareEnableCmd)
	shareCmd.AddCommand(shareDisableCmd)
	shareCmd.AddCommand(shareStatusCmd)
	rootCmd.AddCommand(shareCmd)
}
