package wrfs

import "time"

// ChtimesFile is a file with a Chtimes method.
type ChtimesFile interface {
	File

	// Chtimes changes the access and modification times of the file,
	// similar to the Unix utime() or utimes() functions.
	Chtimes(atime time.Time, mtime time.Time) error
}

// ChtimesFS is a file system with a Chtimes method.
type ChtimesFS interface {
	FS

	// Chtimes changes the access and modification times of the named file,
	// similar to the Unix utime() or utimes() functions.
	Chtimes(name string, atime time.Time, mtime time.Time) error
}

// Chtimes changes the access and modification times of the named file,
// similar to the Unix utime() or utimes() functions.
func Chtimes(fsys FS, name string, atime time.Time, mtime time.Time) (err error) {
	if fsys, ok := fsys.(ChtimesFS); ok {
		return fsys.Chtimes(name, atime, mtime)
	}
	file, err := fsys.Open(name)
	defer safeClose(file, &err)
	if err != nil {
		return err
	}
	if file, ok := file.(ChtimesFile); ok {
		return file.Chtimes(atime, mtime)
	}
	return &PathError{Op: "chtimes", Path: name, Err: ErrNotSupported}
}
