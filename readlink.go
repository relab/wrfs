package wrfs

// ReadlinkFS is a file system that supports the Readlink function.
type ReadlinkFS interface {
	// Readlink returns the destination of the named symbolic link.
	Readlink(name string) (string, error)
}

// Readlink returns the destination of the named symbolic link.
func Readlink(fsys FS, name string) (string, error) {
	if fsys, ok := fsys.(ReadlinkFS); ok {
		return fsys.Readlink(name)
	}
	return "", &PathError{Op: "readlink", Path: name, Err: ErrNotSupported}
}
