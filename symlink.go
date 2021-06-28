package wrfs

import "os"

// SymlinkFS is a file system with a Symlink method.
type SymlinkFS interface {
	FS

	// Symlink creates newname as a symbolic link to oldname.
	Symlink(oldname, newname string) error
}

// Symlink creates newname as a symbolic link to oldname.
func Symlink(fsys FS, oldname, newname string) error {
	if fsys, ok := fsys.(SymlinkFS); ok {
		return fsys.Symlink(oldname, newname)
	}
	return &os.LinkError{Op: "symlink", Old: oldname, New: newname, Err: ErrUnsupported}
}
