package wrfs

import (
	"path/filepath"
)

// RemoveFS is a file system that supports the Remove function.
type RemoveFS interface {
	FS

	// Remove removes the named file or (empty) directory.
	Remove(name string) error
}

// Remove removes the named file or (empty) directory.
func Remove(fsys FS, name string) error {
	if fsys, ok := fsys.(RemoveFS); ok {
		return fsys.Remove(name)
	}
	return &PathError{Op: "remove", Path: name, Err: ErrUnsupported}
}

// RemoveAllFS is a file system that supports the RemoveAll function.
type RemoveAllFS interface {
	FS

	// RemoveAll removes path and any children it contains.
	RemoveAll(path string) error
}

// RemoveAll removes path and any children it contains.
func RemoveAll(fsys FS, path string) error {
	if fsys, ok := fsys.(RemoveAllFS); ok {
		return fsys.RemoveAll(path)
	}

	fi, err := Stat(fsys, path)
	if err != nil {
		return err
	}

	// Check if we are removing a file or a directory.
	if !fi.IsDir() {
		return Remove(fsys, path)
	}

	files, err := ReadDir(fsys, path)
	if err != nil {
		return err
	}

	for _, fi := range files {
		if fi.IsDir() {
			if err = RemoveAll(fsys, filepath.Join(path, fi.Name())); err != nil {
				return err
			}
		} else if err = Remove(fsys, filepath.Join(path, fi.Name())); err != nil {
			return err
		}
	}

	return Remove(fsys, path)
}
