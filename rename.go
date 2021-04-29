package wrfs

import "os"

type RenameFS interface {
	FS

	// Rename renames (moves) oldpath to newpath.
	// If newpath already exists and is not a directory, Rename replaces it.
	Rename(oldpath, newpath string) error
}

// Rename renames (moves) oldpath to newpath.
// If newpath already exists and is not a directory, Rename replaces it.
func Rename(fsys FS, oldpath, newpath string) error {
	if fsys, ok := fsys.(RenameFS); ok {
		return fsys.Rename(oldpath, newpath)
	}
	return &os.LinkError{Op: "rename", Old: oldpath, New: newpath, Err: ErrNotSupported}
}
