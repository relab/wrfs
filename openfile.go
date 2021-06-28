package wrfs

import (
	"io"
	"os"
)

// WriteFile is a file that can be written to.
type WriteFile interface {
	File
	io.Writer
}

// Write writes len(p) bytes from p to the file.
// It returns the number of bytes written from p (0 <= n <= len(p))
// and any error encountered that caused the write to stop early.
func Write(file File, p []byte) (n int, err error) {
	if file, ok := file.(io.Writer); ok {
		return file.Write(p)
	}
	return 0, ErrUnsupported
}

// Seek sets the offset for the next Read or Write on file to offset,
// interpreted according to whence: 0 means relative to the origin of the file,
// 1 means relative to the current offset, and 2 means relative to the end.
// It returns the new offset and an error, if any.
func Seek(file File, offset int64, whence int) (int64, error) {
	if file, ok := file.(io.Seeker); ok {
		return file.Seek(offset, whence)
	}
	return 0, ErrUnsupported
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
	return nil, &PathError{Op: "open", Path: name, Err: ErrUnsupported}
}

// Create creates or truncates the named file. If the file already exists,
// it is truncated. If the file does not exist, it is created with mode 0666
// (before umask). If successful, methods on the returned File can
// be used for I/O; the associated file descriptor has mode O_RDWR.
func Create(fsys FS, name string) (WriteFile, error) {
	file, err := OpenFile(fsys, name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	return file.(WriteFile), err
}
