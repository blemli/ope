package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := os.Args[1]

	switch cmd {
	case "install":
		if err := install(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "uninstall":
		if err := uninstall(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "config":
		path, err := ConfigPath()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Config: %s\n", path)
		cfg, err := LoadConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Blocked: %v\n", cfg.Blocked)
		fmt.Printf("Allowed: %v\n", cfg.Allowed)

	case "version":
		fmt.Printf("ope %s\n", Version)

	default:
		// Treat as ope:// URL
		if err := HandleURL(cmd); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `ope %s — open files and folders from the browser

Usage:
  ope <ope://url>    Open a file or folder
  ope install        Register ope:// URL scheme
  ope uninstall      Unregister ope:// URL scheme
  ope config         Show configuration
  ope version        Print version
`, Version)
}
