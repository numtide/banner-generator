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
		Short: "Generate and upload PNG banners for GitHub repositories",
		Long: `banner-cli is a tool for generating PNG banners from SVG templates
and uploading them to GitHub repositories for social media previews.`,
	}

	// Global flags
	var (
		configPath  string
		githubToken string
		outputPath  string
	)

	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "Path to configuration file")
	rootCmd.PersistentFlags().StringVar(&githubToken, "token", "", "GitHub API token (overrides config)")

	// Generate command
	var generateCmd = &cobra.Command{
		Use:   "generate [owner/repo]",
		Short: "Generate a PNG banner for a repository",
		Args:  cobra.ExactArgs(1),
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

			return generator.GeneratePNG(args[0], outputPath)
		},
	}

	generateCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output path for PNG file (default: banner.png)")
	rootCmd.AddCommand(generateCmd)

	// Upload command
	var uploadCmd = &cobra.Command{
		Use:   "upload [owner/repo] [png-file]",
		Short: "Upload a PNG banner to a GitHub repository",
		Args:  cobra.ExactArgs(2),
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

			uploader := cli.NewUploader(appConfig.GitHub.Token)
			return uploader.UploadBanner(args[0], args[1])
		},
	}

	rootCmd.AddCommand(uploadCmd)

	// Generate and upload command (combined)
	var (
		branch     string
		commitMsg  string
		targetPath string
	)

	var generateUploadCmd = &cobra.Command{
		Use:   "generate-upload [owner/repo]",
		Short: "Generate and upload a PNG banner in one step",
		Args:  cobra.ExactArgs(1),
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

			// Generate PNG to temporary file
			tmpFile, err := os.CreateTemp("", "banner-*.png")
			if err != nil {
				return fmt.Errorf("failed to create temporary file: %w", err)
			}
			tmpPath := tmpFile.Name()
			if err := tmpFile.Close(); err != nil {
				return fmt.Errorf("failed to close temporary file: %w", err)
			}
			defer func() {
				if err := os.Remove(tmpPath); err != nil {
					// Log error but don't fail - temp file cleanup is not critical
					fmt.Fprintf(os.Stderr, "warning: failed to remove temporary file: %v\n", err)
				}
			}()

			if err := generator.GeneratePNG(args[0], tmpPath); err != nil {
				return fmt.Errorf("failed to generate PNG: %w", err)
			}

			// Upload to repository
			uploader := cli.NewUploader(appConfig.GitHub.Token)
			if err := uploader.UploadBannerWithOptions(args[0], tmpPath, branch, targetPath, commitMsg); err != nil {
				return fmt.Errorf("failed to upload banner: %w", err)
			}

			fmt.Println("Banner generated and uploaded successfully!")
			return nil
		},
	}

	generateUploadCmd.Flags().StringVar(&branch, "branch", "main", "Target branch for upload")
	generateUploadCmd.Flags().StringVar(&commitMsg, "message", "Update repository banner", "Commit message")
	generateUploadCmd.Flags().StringVar(&targetPath, "path", "assets/banner.png", "Target path in repository")
	rootCmd.AddCommand(generateUploadCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// loadConfig loads the application configuration
func loadConfig(configPath string) (*config.AppConfig, error) {
	loader := config.NewConfigLoader()
	return loader.LoadConfig(configPath)
}
