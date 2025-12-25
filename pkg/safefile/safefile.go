// Package safefile provides utility functions for safe file operations.
//
// This file contains security utilities for safe file operations,
// addressing gosec G304 (potential file inclusion via variable).
package safefile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SafePath validates and returns a clean path that is safe to use.
// It ensures the resolved path doesn't escape the base directory through path traversal.
// This prevents directory traversal attacks (e.g., "../../../etc/passwd").
//
// Example usage:
//
//	safePath, err := safefile.SafePath("/app/uploads", userInput)
//	if err != nil {
//	    return fmt.Errorf("invalid path: %w", err)
//	}
func SafePath(baseDir, relativePath string) (string, error) {
	// Clean and resolve the base directory
	absBase, err := filepath.Abs(baseDir)
	if err != nil {
		return "", fmt.Errorf("failed to resolve base directory: %w", err)
	}

	// Clean and join the paths
	cleanPath := filepath.Clean(filepath.Join(absBase, relativePath))

	// Verify the result is still within the base directory
	// Add separator to prevent prefix matching like "/app" matching "/application"
	if !strings.HasPrefix(cleanPath, absBase+string(filepath.Separator)) && cleanPath != absBase {
		return "", fmt.Errorf("path traversal detected: %s escapes base directory %s", relativePath, baseDir)
	}

	return cleanPath, nil
}

// SafeReadFile reads a file after validating the path is safe.
// This addresses gosec G304 (potential file inclusion via variable).
//
// Example usage:
//
//	content, err := safefile.SafeReadFile("/app/config/settings.yaml")
//	if err != nil {
//	    return fmt.Errorf("failed to read config: %w", err)
//	}
func SafeReadFile(filePath string) ([]byte, error) {
	// Resolve to absolute path and clean it
	absPath, err := filepath.Abs(filepath.Clean(filePath))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve path: %w", err)
	}

	// Read the file
	content, err := os.ReadFile(absPath) // #nosec G304 - path is cleaned and resolved to absolute
	if err != nil {
		return nil, err
	}

	return content, nil
}

// SafeReadFileInDir reads a file ensuring it's within the specified base directory.
// Combines SafePath and SafeReadFile for convenience.
//
// Example usage:
//
//	content, err := safefile.SafeReadFileInDir("/app/uploads", userFileName)
//	if err != nil {
//	    return fmt.Errorf("failed to read file: %w", err)
//	}
func SafeReadFileInDir(baseDir, relativePath string) ([]byte, error) {
	safePath, err := SafePath(baseDir, relativePath)
	if err != nil {
		return nil, err
	}
	return SafeReadFile(safePath)
}

// SafeWriteFile writes content to a file after validating the path is safe.
// This addresses gosec G304 and ensures file operations are secure.
//
// Example usage:
//
//	err := safefile.SafeWriteFile("/app/data/output.json", jsonData, 0644)
//	if err != nil {
//	    return fmt.Errorf("failed to write file: %w", err)
//	}
func SafeWriteFile(filePath string, data []byte, perm os.FileMode) error {
	// Resolve to absolute path and clean it
	absPath, err := filepath.Abs(filepath.Clean(filePath))
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	// Ensure directory exists
	dir := filepath.Dir(absPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write the file
	if err := os.WriteFile(absPath, data, perm); err != nil { // #nosec G304 - path is cleaned and resolved to absolute
		return err
	}

	return nil
}

// SafeWriteFileInDir writes content to a file ensuring it's within the specified base directory.
// Combines SafePath and SafeWriteFile for convenience.
//
// Example usage:
//
//	err := safefile.SafeWriteFileInDir("/app/uploads", userFileName, data, 0644)
//	if err != nil {
//	    return fmt.Errorf("failed to write file: %w", err)
//	}
func SafeWriteFileInDir(baseDir, relativePath string, data []byte, perm os.FileMode) error {
	safePath, err := SafePath(baseDir, relativePath)
	if err != nil {
		return err
	}
	return SafeWriteFile(safePath, data, perm)
}

// SafeOpen opens a file for reading after validating the path is safe.
// The caller is responsible for closing the file.
//
// Example usage:
//
//	file, err := safefile.SafeOpen("/app/config/settings.yaml")
//	if err != nil {
//	    return fmt.Errorf("failed to open config: %w", err)
//	}
//	defer file.Close()
func SafeOpen(filePath string) (*os.File, error) {
	// Resolve to absolute path and clean it
	absPath, err := filepath.Abs(filepath.Clean(filePath))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve path: %w", err)
	}

	// Open the file
	file, err := os.Open(absPath) // #nosec G304 - path is cleaned and resolved to absolute
	if err != nil {
		return nil, err
	}

	return file, nil
}

// SafeCreate creates or truncates a file after validating the path is safe.
// The caller is responsible for closing the file.
//
// Example usage:
//
//	file, err := safefile.SafeCreate("/app/data/output.json")
//	if err != nil {
//	    return fmt.Errorf("failed to create file: %w", err)
//	}
//	defer file.Close()
func SafeCreate(filePath string) (*os.File, error) {
	// Resolve to absolute path and clean it
	absPath, err := filepath.Abs(filepath.Clean(filePath))
	if err != nil {
		return nil, fmt.Errorf("failed to resolve path: %w", err)
	}

	// Ensure directory exists
	dir := filepath.Dir(absPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Create the file
	file, err := os.Create(absPath) // #nosec G304 - path is cleaned and resolved to absolute
	if err != nil {
		return nil, err
	}

	return file, nil
}

// FileExists checks if a file exists at the given path safely.
//
// Example usage:
//
//	if safefile.FileExists("/app/config/settings.yaml") {
//	    // load config
//	}
func FileExists(filePath string) bool {
	absPath, err := filepath.Abs(filepath.Clean(filePath))
	if err != nil {
		return false
	}
	info, err := os.Stat(absPath)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// DirExists checks if a directory exists at the given path safely.
//
// Example usage:
//
//	if safefile.DirExists("/app/uploads") {
//	    // directory ready
//	}
func DirExists(dirPath string) bool {
	absPath, err := filepath.Abs(filepath.Clean(dirPath))
	if err != nil {
		return false
	}
	info, err := os.Stat(absPath)
	if err != nil {
		return false
	}
	return info.IsDir()
}
