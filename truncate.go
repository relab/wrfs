package wrfs

import "os"

// TruncateFile is a file that supports the Truncate method.
type TruncateFile interface {
	File

	// Truncate changes the size of the file.
	Truncate(size int64) error
}

// TruncateFS is a file system with a Truncate method.
type TruncateFS interface {
	FS

	// Truncate changes the size of the named file.
	Truncate(name string, size int64) error
}

// Truncate changes the size of the named file.
func Truncate(fsys FS, name string, size int64) (err error) {
	if fsys, ok := fsys.(TruncateFS); ok {
		return fsys.Truncate(name, size)
	}

	file, err := OpenFile(fsys, name, os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer safeClose(file, &err)

	if file, ok := file.(TruncateFile); ok {
		return file.Truncate(size)
	}

	// We could try to manually truncate the file if the fs supports OpenFile,
	// but that would be very inefficient.
	return &PathError{Op: "truncate", Path: name, Err: ErrUnsupported}
}
