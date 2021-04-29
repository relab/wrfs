package wrfs

// TruncateFS is a file system with a Truncate method.
type TruncateFS interface {
	FS

	// Truncate changes the size of the named file.
	Truncate(name string, size int64) error
}

// Truncate changes the size of the named file.
func Truncate(fsys FS, name string, size int64) error {
	if fsys, ok := fsys.(TruncateFS); ok {
		return fsys.Truncate(name, size)
	}
	// We could try to manually truncate the file if the fs supports OpenFile,
	// but that would be very inefficient.
	return &PathError{Op: "truncate", Path: name, Err: ErrNotSupported}
}
