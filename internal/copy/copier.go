package copy

import (
	"copier/internal/log"
	"copier/internal/skip"

	"io"
	"os"
	"path/filepath"
)

// CopyDir recursively copies a directory tree, skipping certain directories/files, logging skips and errors.
func CopyDir(src, dst string, logger *log.Logger) error {
	return copyDirInternal(src, dst, logger, nil)
}

// CopyDirWithIgnore recursively copies a directory tree, skipping via .copyignore and default rules.
func CopyDirWithIgnore(src, dst string, logger *log.Logger, ignoreFile string) error {
	var matcher *skip.IgnoreMatcher
	if im, err := skip.LoadIgnoreMatcher(ignoreFile); err == nil {
		matcher = im
	}
	return copyDirInternal(src, dst, logger, matcher)
}

func copyDirInternal(src, dst string, logger *log.Logger, matcher *skip.IgnoreMatcher) error {
	var copyErr error

	err := filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		rel, _ := filepath.Rel(src, path)
		dstPath := filepath.Join(dst, rel)

		if err != nil {
			logger.LogError(path, err)
			copyErr = err
			return nil // continue
		}

		// .copyignore pattern skip
		if matcher != nil && matcher.ShouldIgnore(rel) {
			logger.LogSkip(path)
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if d.IsDir() {
			if skip.ShouldSkipDir(d.Name()) {
				logger.LogSkip(path)
				return filepath.SkipDir
			}
			if path != src {
				info, ierr := d.Info()
				perm := os.FileMode(0755)
				if ierr == nil {
					perm = info.Mode().Perm()
				}
				if mkerr := os.MkdirAll(dstPath, perm); mkerr != nil {
					logger.LogError(dstPath, mkerr)
					copyErr = mkerr
				}
			}
			return nil
		}

		if skip.ShouldSkipFile(d.Name()) {
			logger.LogSkip(path)
			return nil
		}

		// Symlink: copy as link
		if d.Type()&os.ModeSymlink != 0 {
			target, lerr := os.Readlink(path)
			if lerr != nil {
				logger.LogError(path, lerr)
				copyErr = lerr
				return nil
			}
			if lerr := os.Symlink(target, dstPath); lerr != nil {
				logger.LogError(dstPath, lerr)
				copyErr = lerr
			}
			return nil
		}

		// Regular file
		if ferr := copyFile(path, dstPath, d, logger); ferr != nil {
			logger.LogError(path, ferr)
			copyErr = ferr
		}
		return nil
	})
	if err != nil {
		return err
	}
	return copyErr
}

func copyFile(src, dst string, d os.DirEntry, logger *log.Logger) error {
	srcF, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcF.Close()

	info, err := d.Info()
	if err != nil {
		return err
	}

	dstF, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	defer dstF.Close()

	_, err = io.Copy(dstF, srcF)
	if err != nil {
		return err
	}

	// Preserve permissions
	if cherr := os.Chmod(dst, info.Mode()); cherr != nil {
		logger.LogError(dst, cherr)
	}
	return nil
}
