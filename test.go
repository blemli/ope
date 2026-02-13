package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func setupTestFiles() error {
	var base string
	switch runtime.GOOS {
	case "windows":
		base = filepath.Join(os.TempDir())
	default:
		base = "/tmp"
	}

	// Directories
	dirs := []string{
		filepath.Join(base, "my folder"),
	}

	// Files (empty — just need to exist for security dialog tests)
	files := []string{
		filepath.Join(base, "ope-test-1.txt"),
		filepath.Join(base, "ope-test-2.txt"),
		filepath.Join(base, "notes.txt"),
		filepath.Join(base, "document.pdf"),
		filepath.Join(base, "photo.png"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("mkdir %s: %w", dir, err)
		}
		fmt.Printf("  dir:  %s\n", dir)
	}

	for _, file := range files {
		if err := os.WriteFile(file, []byte{}, 0o644); err != nil {
			return fmt.Errorf("create %s: %w", file, err)
		}
		fmt.Printf("  file: %s\n", file)
	}

	fmt.Printf("\nTest files created in %s\n", base)
	fmt.Println("Open test.html in your browser to run the test suite.")
	return nil
}
