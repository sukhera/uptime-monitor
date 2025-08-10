package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sukhera/uptime-monitor/internal/shared/logger"
)

var (
	cfgFile   string
	verbose   bool
	version   string
	commit    string
	buildDate string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "status-page",
	Short: "A comprehensive status page monitoring system",
	Long: `Status Page Starter is a complete monitoring solution that provides:

- Real-time service health monitoring
- RESTful API for status data
- Web dashboard for status visualization
- Database persistence for historical data
- Docker support for easy deployment

Features:
- Monitor multiple services simultaneously
- Configurable health check intervals
- REST API for integration
- Web dashboard with real-time updates
- Database persistence with MongoDB
- Docker containerization
- Comprehensive testing suite`,
	Version: version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// SetVersion sets the version for the root command
func SetVersion(v string) {
	version = v
	rootCmd.Version = v
}

// SetBuildInfo sets all build information
func SetBuildInfo(v, c, bd string) {
	version = v
	commit = c
	buildDate = bd
	rootCmd.Version = v
}

// GetBuildInfo returns the build information
func GetBuildInfo() (string, string, string) {
	return version, commit, buildDate
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Local flags
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".status-page" (without extension).
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	// Set environment variable prefix and enable automatic env binding
	viper.SetEnvPrefix("STATUS_PAGE")
	viper.AutomaticEnv()

	// Map environment variables to viper keys for compatibility
	setupEnvBindings()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if verbose {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}
}

// setupEnvBindings maps common environment variables to viper keys
func setupEnvBindings() {
	// Server/API environment variables
	_ = viper.BindEnv("api.port", "PORT", "API_PORT")
	_ = viper.BindEnv("server.port", "PORT", "SERVER_PORT")
	_ = viper.BindEnv("server.read_timeout", "READ_TIMEOUT")
	_ = viper.BindEnv("server.write_timeout", "WRITE_TIMEOUT")
	_ = viper.BindEnv("server.idle_timeout", "IDLE_TIMEOUT")

	// Database environment variables
	_ = viper.BindEnv("database.url", "MONGO_URI", "DB_URL", "DATABASE_URL")
	_ = viper.BindEnv("database.name", "DB_NAME", "DATABASE_NAME")
	_ = viper.BindEnv("database.timeout", "DB_TIMEOUT")

	// Logging environment variables
	_ = viper.BindEnv("logging.level", "LOG_LEVEL")
	_ = viper.BindEnv("logging.json", "LOG_JSON")

	// Checker environment variables
	_ = viper.BindEnv("checker.interval", "CHECK_INTERVAL", "CHECKER_INTERVAL")

	// Web environment variables
	_ = viper.BindEnv("web.port", "WEB_PORT")
	_ = viper.BindEnv("web.api_url", "API_URL")
	_ = viper.BindEnv("web.static_dir", "STATIC_DIR")
}

// bindFlagToViper binds a cobra command flag to viper with error handling
func bindFlagToViper(cmd *cobra.Command, viperKey, flagName string) {
	ctx := context.Background()
	log := logger.Get()

	if err := viper.BindPFlag(viperKey, cmd.Flags().Lookup(flagName)); err != nil {
		log.Fatal(ctx, fmt.Sprintf("Failed to bind %s flag", viperKey), err, nil)
	}
}
