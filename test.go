package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	// Method 1: Using os.Getenv() to read individual config variables
	databaseURL := os.Getenv("DATABASE_URL")
	apiKey := os.Getenv("API_KEY")

	fmt.Println("Database URL:", databaseURL)
	fmt.Println("API Key:", apiKey)

	// Method 2: Retrieve all environment variables
	fmt.Println("\nAll Environment Variables:")
	for _, env := range os.Environ() {
		fmt.Println(env)
	}

	// Method 3: Using a map to store all environment variables
	envMap := make(map[string]string)
	for _, env := range os.Environ() {
		key, value, found := strings.Cut(env, "=")
		if found {
			envMap[key] = value
		}
	}

	// Accessing a specific variable from the map
	if specificVar, exists := envMap["SPECIFIC_CONFIG"]; exists {
		fmt.Println("\nSpecific Config Variable:", specificVar)
	}

	fmt.Println("hello world")
}
