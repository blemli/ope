package main

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestParseOpeURL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"triple slash", "ope:///tmp", "/tmp", false},
		{"with localhost", "ope://localhost/tmp", "/tmp", false},
		{"nested path", "ope:///home/user/docs", "/home/user/docs", false},
		{"path with spaces", "ope:///tmp/my%20folder", "/tmp/my folder", false},
		{"wrong scheme", "http://example.com", "", true},
		{"empty", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseOpeURL(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseOpeURL(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseOpeURL(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestExpandPath(t *testing.T) {
	home, _ := os.UserHomeDir()

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"absolute path", "/tmp", "/tmp", false},
		{"tilde", "~/Documents", filepath.Join(home, "Documents"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExpandPath(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExpandPath(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ExpandPath(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestExpandPathGlob(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("glob test uses /tmp")
	}

	// Create a temp file to glob for
	dir := t.TempDir()
	f, err := os.Create(filepath.Join(dir, "testfile.txt"))
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	got, err := ExpandPath(filepath.Join(dir, "*.txt"))
	if err != nil {
		t.Fatalf("ExpandPath glob: %v", err)
	}
	if got != filepath.Join(dir, "testfile.txt") {
		t.Errorf("ExpandPath glob = %q, want %q", got, filepath.Join(dir, "testfile.txt"))
	}
}

func TestCheckSecurity(t *testing.T) {
	cfg := &Config{
		Blocked: []string{"*.exe", "*.bat"},
		Allowed: []string{"readme.txt", "*.pdf"},
	}

	tests := []struct {
		name string
		path string
		want SecurityAction
	}{
		{"blocked exe", "/tmp/virus.exe", ActionBlock},
		{"blocked bat", "/tmp/script.bat", ActionBlock},
		{"blocked case insensitive", "/tmp/VIRUS.EXE", ActionBlock},
		{"allowed exact", "/tmp/readme.txt", ActionAllow},
		{"allowed glob", "/tmp/document.pdf", ActionAllow},
		{"unknown file", "/tmp/photo.jpg", ActionAsk},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cfg.CheckSecurity(tt.path)
			if got != tt.want {
				t.Errorf("CheckSecurity(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestCheckSecurityDirectory(t *testing.T) {
	cfg := &Config{
		Blocked: []string{},
		Allowed: []string{},
	}

	dir := t.TempDir()
	got := cfg.CheckSecurity(dir)
	if got != ActionAllow {
		t.Errorf("CheckSecurity(dir) = %v, want ActionAllow", got)
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if len(cfg.Blocked) == 0 {
		t.Error("DefaultConfig should have blocked extensions")
	}
	if len(cfg.Allowed) != 0 {
		t.Error("DefaultConfig should have empty allowed list")
	}
}
