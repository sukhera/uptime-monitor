package cmd

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	if err := viper.BindPFlag("web.port", webCmd.Flags().Lookup("port")); err != nil {
		log.Fatalf("Failed to bind web.port flag: %v", err)
	}
	if err := viper.BindPFlag("web.api_url", webCmd.Flags().Lookup("api-url")); err != nil {
		log.Fatalf("Failed to bind web.api_url flag: %v", err)
	}
	if err := viper.BindPFlag("web.static_dir", webCmd.Flags().Lookup("static-dir")); err != nil {
		log.Fatalf("Failed to bind web.static_dir flag: %v", err)
	}
}

func runWeb(cmd *cobra.Command, args []string) {
	// Check if static directory exists
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		log.Printf("Static directory not found: %s", staticDir)
		log.Println("Please build the React app first:")
		log.Println("  cd web/react-status-page && npm run build")
		log.Println("  or use Docker: docker-compose up web")
		os.Exit(1)
	}

	// Create file server
	fs := http.FileServer(http.Dir(staticDir))

	// Create server with SPA routing support
	mux := http.NewServeMux()

	// Serve static files
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate and sanitize the requested path
		requestedPath := r.URL.Path

		// Prevent path traversal attacks by ensuring the path doesn't contain ".."
		if strings.Contains(requestedPath, "..") {
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		// Ensure the path starts with "/" and normalize it
		if !strings.HasPrefix(requestedPath, "/") {
			requestedPath = "/" + requestedPath
		}

		// Additional security check: ensure the resolved path is within the static directory
		absStaticDir, err := filepath.Abs(staticDir)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Use a more secure approach: only allow specific file types and patterns
		// This prevents any potential path injection while still supporting SPA routing
		cleanPath := strings.TrimPrefix(requestedPath, "/")

		// Only allow specific file patterns for security
		allowedPatterns := []string{
			"",           // Root path
			"index.html", // Main SPA file
			"assets/",    // Asset directory
			"static/",    // Static files directory
		}

		isAllowed := false
		for _, pattern := range allowedPatterns {
			if cleanPath == pattern || strings.HasPrefix(cleanPath, pattern) {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			// For any other path, serve index.html for SPA routing
			http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
			return
		}

		// Construct the file path safely after validation
		filePath := filepath.Join(staticDir, cleanPath)

		// Final security check: ensure the resolved path is within the static directory
		absFilePath, err := filepath.Abs(filePath)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Ensure the file path is within the static directory
		if !strings.HasPrefix(absFilePath, absStaticDir) {
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		// Check if the file exists using the validated path
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			// File doesn't exist, serve index.html for SPA routing
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

		log.Println("Shutting down web server...")
		if err := server.Close(); err != nil {
			log.Printf("Error shutting down server: %v", err)
		}
	}()

	// Start server
	log.Printf("Starting web server on port %s", webPort)
	log.Printf("Serving static files from: %s", staticDir)
	log.Printf("Web dashboard: http://localhost:%s", webPort)
	log.Printf("API server should be running at: %s", apiURL)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start web server: %v", err)
	}
}
