package main

import (
	"log"
	"time"

	"github.com/sukhera/uptime-monitor/internal/shared/config"
)

// This example demonstrates the benefits of the functional options pattern
// over the previous configuration approach.

func main() {
	log.Println("=== Functional Options Pattern Demo ===")

	// Example 1: Default configuration
	log.Println("\n1. Default Configuration:")
	defaultCfg := config.New()
	log.Printf("   Server Port: %s", defaultCfg.Server.Port)
	log.Printf("   Database: %s/%s", defaultCfg.Database.URI, defaultCfg.Database.Name)
	log.Printf("   Check Interval: %v", defaultCfg.Checker.Interval)

	// Example 2: Environment-based configuration
	log.Println("\n2. Environment-based Configuration:")
	envCfg := config.New(config.FromEnvironment())
	log.Printf("   Server Port: %s", envCfg.Server.Port)
	log.Printf("   Database: %s/%s", envCfg.Database.URI, envCfg.Database.Name)
	log.Printf("   Check Interval: %v", envCfg.Checker.Interval)

	// Example 3: Custom configuration with options
	log.Println("\n3. Custom Configuration with Options:")
	customCfg := config.New(
		config.WithServerPort("9090"),
		config.WithDatabase("mongodb://custom:27017", "custom_db", 15*time.Second),
		config.WithLogging("debug", true),
		config.WithCheckerInterval(5*time.Minute),
	)
	log.Printf("   Server Port: %s", customCfg.Server.Port)
	log.Printf("   Database: %s/%s", customCfg.Database.URI, customCfg.Database.Name)
	log.Printf("   Log Level: %s (JSON: %t)", customCfg.Logging.Level, customCfg.Logging.JSON)
	log.Printf("   Check Interval: %v", customCfg.Checker.Interval)

	// Example 4: Partial configuration (only override what you need)
	log.Println("\n4. Partial Configuration:")
	partialCfg := config.New(
		config.FromEnvironment(),      // Start with environment defaults
		config.WithServerPort("8081"), // Override only the port
	)
	log.Printf("   Server Port: %s", partialCfg.Server.Port)
	log.Printf("   Database: %s/%s", partialCfg.Database.URI, partialCfg.Database.Name)

	// Example 5: Configuration validation
	log.Println("\n5. Configuration Validation:")
	invalidCfg := config.New(
		config.WithServerPort(""),     // Invalid: empty port
		config.WithCheckerInterval(0), // Invalid: zero interval
	)
	if err := invalidCfg.Validate(); err != nil {
		log.Printf("   Validation Error: %v", err)
	} else {
		log.Println("   Configuration is valid")
	}

	// Example 6: Backward compatibility
	log.Println("\n6. Backward Compatibility:")
	legacyCfg := config.Load() // Old way still works
	log.Printf("   Server Port: %s", legacyCfg.Server.Port)
	log.Printf("   Database: %s/%s", legacyCfg.Database.URI, legacyCfg.Database.Name)

	log.Println("\n=== Demo Complete ===")
}

// Benefits demonstrated:
//
// 1. **Flexibility**: Easy to create different configurations for different scenarios
// 2. **Composability**: Options can be combined in any order
// 3. **Testability**: No need to set environment variables for testing
// 4. **Validation**: Built-in validation prevents runtime errors
// 5. **Backward Compatibility**: Existing code continues to work
// 6. **Clarity**: Configuration intent is explicit and readable
// 7. **Maintainability**: Easy to add new configuration options
