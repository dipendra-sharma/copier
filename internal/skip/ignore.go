package skip

import (
	"os"
	"strings"
	"github.com/sabhiram/go-gitignore"
)

type IgnoreMatcher struct {
	matcher *ignore.GitIgnore
}

// LoadIgnoreMatcher loads patterns from the given ignore file path.
func LoadIgnoreMatcher(ignorePath string) (*IgnoreMatcher, error) {
	data, err := os.ReadFile(ignorePath)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	matcher := ignore.CompileIgnoreLines(lines...)
	return &IgnoreMatcher{matcher: matcher}, nil
}

// ShouldIgnore returns true if the given path should be ignored.
func (im *IgnoreMatcher) ShouldIgnore(relPath string) bool {
	if im == nil || im.matcher == nil {
		return false
	}
	return im.matcher.MatchesPath(relPath)
}
