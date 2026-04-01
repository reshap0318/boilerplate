package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/reshap0318/go-boilerplate/internal/helpers"
)

func main() {
	// Find .env file
	envPath, err := findEnvFile()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Generate secure JWT secret using helper
	secret, err := helpers.GenerateRandomString(32)
	if err != nil {
		log.Fatalf("Error generating key: %v", err)
	}

	// Read .env file
	content, err := os.ReadFile(envPath)
	if err != nil {
		log.Fatalf("Error reading .env file: %v", err)
	}

	// Update JWT_SECRET
	updated := updateJWTSecret(string(content), secret)

	// Write back to .env file
	if err := os.WriteFile(envPath, []byte(updated), 0644); err != nil {
		log.Fatalf("Error writing .env file: %v", err)
	}

	fmt.Println("✅ JWT_SECRET generated successfully!")
	fmt.Printf("📁 Updated: %s\n", envPath)
	fmt.Printf("🔑 Secret: %s\n", secret)
}

// findEnvFile searches for .env file in current directory and parent directories
func findEnvFile() (string, error) {
	// Check current directory first
	if _, err := os.Stat(".env"); err == nil {
		return ".env", nil
	}

	// Get absolute path of current directory
	absPath, err := filepath.Abs(".")
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Search in parent directories (max 5 levels)
	current := absPath
	for i := 0; i < 4; i++ {
		envPath := filepath.Join(current, ".env")
		if _, err := os.Stat(envPath); err == nil {
			return envPath, nil
		}
		current = filepath.Dir(current)
	}

	return "", fmt.Errorf(".env file not found (create one or run from project root)")
}

// updateJWTSecret updates or adds JWT_SECRET in .env content
func updateJWTSecret(content, secret string) string {
	lines := strings.Split(content, "\n")
	found := false

	for i, line := range lines {
		if strings.HasPrefix(line, "JWT_SECRET=") {
			lines[i] = "JWT_SECRET=" + secret
			found = true
			break
		}
	}

	// If JWT_SECRET not found, add it
	if !found {
		lines = append(lines, "JWT_SECRET="+secret)
	}

	return strings.Join(lines, "\n")
}
