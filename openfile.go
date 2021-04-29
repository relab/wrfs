package wrfs

import (
	"errors"
	"io"
	"os"
)

var ErrNotSupported = errors.New("operation not supported")

// Write writes len(p) bytes from p to the file.
// It returns the number of bytes written from p (0 <= n <= len(p))
// and any error encountered that caused the write to stop early.
func Write(file File, p []byte) (n int, err error) {
	if file, ok := file.(io.Writer); ok {
		return file.Write(p)
	}
	return 0, ErrNotSupported
}

// OpenFileFS is a file system that supports the OpenFile function.
type OpenFileFS interface {
	FS

	// OpenFile opens the named file with specified flag (O_RDONLY etc.).
	// If the file does not exist, and the O_CREATE flag is passed, it is created with mode perm (before umask).
	// If successful, methods on the returned File can be used for I/O.
	OpenFile(name string, flag int, perm FileMode) (File, error)
}

// OpenFile opens the named file with specified flag (O_RDONLY etc.).
// If the file does not exist, and the O_CREATE flag is passed, it is created with mode perm (before umask).
// If successful, methods on the returned File can be used for I/O.
func OpenFile(fsys FS, name string, flag int, perm FileMode) (File, error) {
	if fsys, ok := fsys.(OpenFileFS); ok {
		return fsys.OpenFile(name, flag, perm)
	}
	if flag == os.O_RDONLY {
		return fsys.Open(name)
	}
	return nil, &PathError{Op: "open", Path: name, Err: ErrNotSupported}
}
