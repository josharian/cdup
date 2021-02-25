// Package cdup implements searching for a file (or directory) in all parents of a given directory.
// This frequently comes up for files whose presence marks all subdirectories
// as part of a coherent whole, such as go.mod for Go modules or .git for git repositories.
package cdup

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// Find searches dir and all its ancestors directories for a file (or directory) named name.
// It returns the most recent ancestor of dir that contains name.
func Find(dir, name string) (parent string, err error) {
	// There's a lot of duplication here with FindIn.
	// But the amount of work to get os.DirFS to work correctly across platforms
	// is more than the amount of work required to just duplicate the code.
	if dir == "" {
		return "", errors.New("dir cannot be empty")
	}
	if name == "" {
		return "", errors.New("name cannot be empty")
	}
	dir = filepath.Clean(dir)
	for {
		candidate := filepath.Join(dir, name)
		_, err := os.Stat(candidate)
		if err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			// Hit root.
			return "", os.ErrNotExist
		}
		// Pop up a directory.
		dir = parent
	}
}

// FindIn searches dir and all its ancestors directories for a file (or directory) named name.
// It returns the most recent ancestor of dir that contains name.
func FindIn(fsys fs.FS, dir, name string) (parent string, err error) {
	if dir == "" {
		return "", errors.New("dir cannot be empty")
	}
	if name == "" {
		return "", errors.New("name cannot be empty")
	}
	dir = filepath.Clean(dir)
	for {
		candidate := filepath.Join(dir, name)
		if !fs.ValidPath(candidate) {
			return "", fmt.Errorf("invalid path: %q", candidate)
		}
		_, err := fs.Stat(fsys, candidate)
		if err == nil {
			return dir, nil
		}
		if dir == "." || dir == "/" {
			// Hit root.
			return "", os.ErrNotExist
		}
		// Pop up a directory.
		dir = filepath.Dir(dir)
	}
}
