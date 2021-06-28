package wrfs

// LstatFS is a file system that supports the Lstat operation.
type LstatFS interface {
	// Lstat returns a FileInfo describing the named file.
	// If the file is a symbolic link, the returned FileInfo describes the symbolic link.
	// Lstat makes no attempt to follow the link.
	Lstat(name string) (FileInfo, error)
}

// Lstat returns a FileInfo describing the named file.
// If the file is a symbolic link, the returned FileInfo describes the symbolic link.
// Lstat makes no attempt to follow the link.
func Lstat(fsys FS, name string) (info FileInfo, err error) {
	if fsys, ok := (fsys.(LstatFS)); ok {
		return fsys.Lstat(name)
	}
	return nil, &PathError{Op: "lstat", Path: name, Err: ErrUnsupported}
}
