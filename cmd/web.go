package cmd

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sukhera/uptime-monitor/internal/shared/logger"
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Serve the web dashboard",
	Long: `Serve the React web dashboard that provides:

- Real-time status dashboard
- Service health visualization
- Incident management interface
- Dark/light theme support
- Responsive design

The web dashboard connects to the API server for data.
Example:
  status-page web --port 3000 --api-url http://localhost:8080`,
	Run: runWeb,
}

var (
	webPort   string
	apiURL    string
	staticDir string
)

func init() {
	rootCmd.AddCommand(webCmd)

	// Web-specific flags
	webCmd.Flags().StringVarP(&webPort, "port", "p", "3000", "Port to serve the web dashboard on")
	webCmd.Flags().StringVar(&apiURL, "api-url", "http://localhost:8080", "API server URL")
	webCmd.Flags().StringVar(&staticDir, "static-dir", "./web/react-status-page/dist", "Static files directory")

	// Bind flags to viper
	bindFlagToViper(webCmd, "web.port", "port")
	bindFlagToViper(webCmd, "web.api_url", "api-url")
	bindFlagToViper(webCmd, "web.static_dir", "static-dir")
}

func runWeb(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	log := logger.Get()

	// Get configuration values from Viper
	webPort := viper.GetString("web.port")
	apiURL := viper.GetString("web.api_url")
	staticDir := viper.GetString("web.static_dir")

	// Check if static directory exists
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		log.Error(ctx, "Static directory not found", nil, logger.Fields{
			"static_dir": staticDir,
			"instructions": []string{
				"cd web/react-status-page && npm run build",
				"or use Docker: docker-compose up web",
			},
		})
		os.Exit(1)
	}

	// Create file server
	fs := http.FileServer(http.Dir(staticDir))

	// Create server with SPA routing support
	mux := http.NewServeMux()

	// Serve static files with enhanced security
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the requested path
		requestedPath := r.URL.Path

		// Security: Use a strict allowlist approach to completely avoid user-controlled data in path construction
		allowedFiles := map[string]string{
			"/":            "index.html",
			"/index.html":  "index.html",
			"/favicon.ico": "favicon.ico",
		}

		// Look up the file in the allowlist
		safeFile, ok := allowedFiles[requestedPath]
		if !ok {
			// For any other path, serve index.html for SPA routing
			safeFile = "index.html"
		}

		filePath := filepath.Join(staticDir, safeFile)

		// Security: Ensure the resolved path is within the static directory using filepath.Rel
		absStaticDir, err := filepath.Abs(staticDir)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		absFilePath, err := filepath.Abs(filePath)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		relPath, err := filepath.Rel(absStaticDir, absFilePath)
		if err != nil || strings.HasPrefix(relPath, "..") || strings.Contains(relPath, string(os.PathSeparator)+"..") {
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		// Check if the file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
			return
		}
		fs.ServeHTTP(w, r)
	}))

	// Create server
	server := &http.Server{
		Addr:              ":" + webPort,
		Handler:           mux,
		ReadHeaderTimeout: 30 * time.Second, // Prevent Slowloris attacks
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Info(ctx, "Shutting down web server", nil)

		// Create context with timeout for graceful shutdown
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Error(ctx, "Error during graceful shutdown, forcing close", err, nil)
			if err := server.Close(); err != nil {
				log.Error(ctx, "Error forcing server close", err, nil)
			}
		} else {
			log.Info(ctx, "Web server shutdown completed gracefully", nil)
		}
	}()

	// Start server
	log.Info(ctx, "Starting web server", logger.Fields{
		"port":          webPort,
		"static_dir":    staticDir,
		"dashboard_url": "http://localhost:" + webPort,
		"api_url":       apiURL,
	})

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(ctx, "Failed to start web server", err, logger.Fields{"port": webPort})
	}
}
