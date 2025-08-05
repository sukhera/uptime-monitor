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

	// Serve static files with enhanced security
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the requested path and normalize it
		requestedPath := r.URL.Path

		// Normalize the path to prevent path traversal attacks
		// Remove any ".." sequences and normalize separators
		cleanPath := filepath.Clean(requestedPath)

		// Additional security: ensure the path doesn't contain ".." after cleaning
		if strings.Contains(cleanPath, "..") {
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		// Ensure the path starts with "/"
		if !strings.HasPrefix(cleanPath, "/") {
			cleanPath = "/" + cleanPath
		}

		// Remove leading slash for file system operations
		relativePath := strings.TrimPrefix(cleanPath, "/")

		// Security: Only allow specific file patterns and extensions
		allowedExtensions := []string{".html", ".css", ".js", ".json", ".ico", ".png", ".jpg", ".jpeg", ".gif", ".svg", ".woff", ".woff2", ".ttf", ".eot"}
		allowedPrefixes := []string{"assets/", "static/", ""}

		// Check if the path has an allowed extension or prefix
		isAllowed := false

		// Check for allowed extensions
		for _, ext := range allowedExtensions {
			if strings.HasSuffix(relativePath, ext) {
				isAllowed = true
				break
			}
		}

		// Check for allowed prefixes (for directories)
		if !isAllowed {
			for _, prefix := range allowedPrefixes {
				if relativePath == prefix || strings.HasPrefix(relativePath, prefix) {
					isAllowed = true
					break
				}
			}
		}

		// If not allowed, serve index.html for SPA routing
		if !isAllowed {
			http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
			return
		}

		// Construct the file path using only validated components
		filePath := filepath.Join(staticDir, relativePath)

		// Security: Ensure the resolved path is within the static directory
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

		// Ensure the file path is within the static directory
		if !strings.HasPrefix(absFilePath, absStaticDir) {
			http.Error(w, "Invalid path", http.StatusBadRequest)
			return
		}

		// Check if the file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			// File doesn't exist, serve index.html for SPA routing
			http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
			return
		}

		// Serve the file using the file server
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
