package main

import (
	"fmt"
	"os"

	"github.com/numtide/banner-generator/internal/cli"
	"github.com/numtide/banner-generator/internal/config"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "banner-cli",
		Short: "Generate PNG banners for GitHub repository social previews",
		Long: `banner-cli generates PNG banners from SVG templates for use as
GitHub repository social media previews.

After generating the banner, upload it manually via:
  Repository Settings > Social preview > Edit`,
		SilenceUsage: true,
	}

	// Global flags
	var (
		configPath  string
		githubToken string
		outputPath  string
		noStats     bool
		darkMode    bool
	)

	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "Path to configuration file")
	rootCmd.PersistentFlags().StringVar(&githubToken, "token", "", "GitHub API token (overrides config)")

	// Generate command
	var generateCmd = &cobra.Command{
		Use:   "generate [owner/repo]",
		Short: "Generate a PNG banner for a repository",
		Long: `Generate a PNG banner for a GitHub repository.

After generating, upload the banner as social preview via:
  https://github.com/OWNER/REPO/settings > Social preview > Edit`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load configuration
			appConfig, err := loadConfig(configPath)
			if err != nil {
				return fmt.Errorf("failed to load configuration: %w", err)
			}

			// Override token if provided
			if githubToken != "" {
				appConfig.GitHub.Token = githubToken
			}

			generator, err := cli.NewGeneratorWithConfig(appConfig)
			if err != nil {
				return fmt.Errorf("failed to initialize generator: %w", err)
			}

			repoPath := args[0]
			if err := generator.GeneratePNG(repoPath, outputPath, noStats, darkMode); err != nil {
				return err
			}

			// Print instructions for setting social preview
			fmt.Println()
			fmt.Println("To set as social preview, go to:")
			fmt.Printf("  https://github.com/%s/settings\n", repoPath)
			fmt.Println("Then scroll to 'Social preview' and click 'Edit' to upload the generated PNG.")

			return nil
		},
	}

	generateCmd.Flags().StringVarP(&outputPath, "output", "o", "banner.png", "Output path for PNG file")
	generateCmd.Flags().BoolVar(&noStats, "no-stats", false, "Omit stars, forks, and language from banner")
	generateCmd.Flags().BoolVar(&darkMode, "dark", false, "Use dark color scheme (default is light)")
	rootCmd.AddCommand(generateCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// loadConfig loads the application configuration
func loadConfig(configPath string) (*config.AppConfig, error) {
	loader := config.NewConfigLoader()
	return loader.LoadConfig(configPath)
}
