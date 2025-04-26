package skip

import "strings"

// ShouldSkipDir returns true if a directory should be skipped (by name).
func ShouldSkipDir(name string) bool {
	skips := []string{
		".git", "node_modules", "build", "dist", ".cache", ".idea", ".vscode", "target", "venv", "__pycache__",
	}
	for _, skip := range skips {
		if strings.EqualFold(name, skip) {
			return true
		}
	}
	return false
}

// ShouldSkipFile returns true if a file should be skipped (by name).
func ShouldSkipFile(name string) bool {
	skips := []string{".DS_Store"}
	for _, skip := range skips {
		if strings.EqualFold(name, skip) {
			return true
		}
	}
	return false
}
