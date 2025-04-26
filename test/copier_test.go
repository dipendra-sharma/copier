package test

import (
	"copier/internal/copy"
	"copier/internal/log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupTestDir(t *testing.T) (string, func()) {
	dir, err := os.MkdirTemp("", "copier_test_src_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	cleanup := func() { os.RemoveAll(dir) }
	return dir, cleanup
}

func setupTestDst(t *testing.T) (string, func()) {
	dir, err := os.MkdirTemp("", "copier_test_dst_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	cleanup := func() { os.RemoveAll(dir) }
	return dir, cleanup
}

func TestCopy_SkipsAndErrors(t *testing.T) {
	src, cleanupSrc := setupTestDir(t)
	defer cleanupSrc()
	dst, cleanupDst := setupTestDst(t)
	defer cleanupDst()
	logfile := filepath.Join(src, "copy_test.log")
	logger, err := log.NewLogger(logfile)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	defer logger.Close()

	// Create files/dirs to test skip logic
	os.Mkdir(filepath.Join(src, ".git"), 0755)
	os.Mkdir(filepath.Join(src, "node_modules"), 0755)
	os.Mkdir(filepath.Join(src, "build"), 0755)
	os.Mkdir(filepath.Join(src, "keepme"), 0755)
	os.WriteFile(filepath.Join(src, "keepme", "file.txt"), []byte("hello"), 0644)
	os.Chmod(filepath.Join(src, "keepme", "file.txt"), 0644) // ensure readable
	os.WriteFile(filepath.Join(src, ".DS_Store"), []byte("junk"), 0644)
	os.WriteFile(filepath.Join(src, "file.txt"), []byte("rootfile"), 0644)

	// Symlink
	symlinkTarget := filepath.Join(src, "keepme", "file.txt")
	os.Symlink(symlinkTarget, filepath.Join(src, "symlink"))

	// Intentionally cause error (unreadable file)
	errfile := filepath.Join(src, "keepme", "errfile.txt")
	os.WriteFile(errfile, []byte("err"), 0644)
	os.Chmod(errfile, 0000)
	defer os.Chmod(errfile, 0644) // so cleanup works

	err = copy.CopyDir(src, dst, logger)
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	// Check skipped dirs/files
	logData, _ := os.ReadFile(logfile)
	logStr := string(logData)
	for _, skip := range []string{".git", "node_modules", "build", ".DS_Store"} {
		if !strings.Contains(logStr, skip) {
			t.Errorf("skip not logged: %s", skip)
		}
	}

	// Check error logged
	if !strings.Contains(logStr, "errfile.txt") {
		t.Errorf("error not logged for unreadable file")
	}

	// Check that keepme/file.txt copied
	dstFile := filepath.Join(dst, "keepme", "file.txt")
	if info, err := os.Stat(dstFile); err != nil {
		t.Errorf("file not copied: %v", err)
	} else {
		t.Logf("Destination file permissions: %v", info.Mode())
	}
	dstDir := filepath.Join(dst, "keepme")
	if info, err := os.Stat(dstDir); err == nil {
		t.Logf("Destination directory permissions: %v", info.Mode())
	}

	// Symlink copied as link
	link := filepath.Join(dst, "symlink")
	if info, err := os.Lstat(link); err != nil || (info.Mode()&os.ModeSymlink == 0) {
		t.Errorf("symlink not copied as link")
	}

	// Do NOT check that errfile.txt was copied, only that the error was logged (already checked above)
}
