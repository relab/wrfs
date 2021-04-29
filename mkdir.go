package wrfs

import (
	"os"
	"syscall"
)

// MkdirFS is a file system that supports the Mkdir function.
type MkdirFS interface {
	FS

	// Mkdir creates a new directory with the specified name and permission bits.
	Mkdir(name string, perm FileMode) error
}

// Mkdir creates a new directory with the specified name and permission bits.
func Mkdir(fsys FS, name string, perm FileMode) error {
	if fsys, ok := fsys.(MkdirFS); ok {
		return fsys.Mkdir(name, perm)
	}
	return &PathError{Op: "mkdir", Path: name, Err: ErrNotSupported}
}

type MkdirAllFS interface {
	FS

	// MkdirAll creates a directory named path, along with any necessary parents, and returns nil,
	// or else returns an error. The permission bits perm (before umask) are used for all
	// directories that MkdirAll creates. If path is already a directory, MkdirAll does nothing
	// and returns nil.
	MkdirAll(path string, perm FileMode) error
}

// MkdirAll creates a directory named path, along with any necessary parents, and returns nil,
// or else returns an error. The permission bits perm (before umask) are used for all
// directories that MkdirAll creates. If path is already a directory, MkdirAll does nothing
// and returns nil.
func MkdirAll(fsys FS, path string, perm FileMode) error {
	if fsys, ok := fsys.(MkdirAllFS); ok {
		return fsys.MkdirAll(path, perm)
	}

	fsys, ok := fsys.(MkdirFS)
	if !ok {
		return &PathError{Op: "mkdir", Path: path, Err: ErrNotSupported}
	}

	// Based on os.MkdirAll

	// Fast path: if we can tell whether path is a directory or file, stop with success or error.
	dir, err := Stat(fsys, path)
	if err == nil {
		if dir.IsDir() {
			return nil
		}
		return &PathError{Op: "mkdir", Path: path, Err: syscall.ENOTDIR}
	}

	// Slow path: make sure parent exists and then call Mkdir for path.
	i := len(path)
	for i > 0 && os.IsPathSeparator(path[i-1]) { // Skip trailing path separator.
		i--
	}

	j := i
	for j > 0 && !os.IsPathSeparator(path[j-1]) { // Scan backward over element.
		j--
	}

	if j > 1 {
		// Create parent.
		err = MkdirAll(fsys, path[:j-1], perm)
		if err != nil {
			return err
		}
	}

	// Parent now exists; invoke Mkdir and use its result.
	err = Mkdir(fsys, path, perm)
	if err != nil {
		// Handle arguments like "foo/." by
		// double-checking that directory doesn't exist.
		dir, err1 := Stat(fsys, path)
		if err1 == nil && dir.IsDir() {
			return nil
		}
		return err
	}
	return nil
}
