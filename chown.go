package wrfs

// ChownFile is a file with a Chown method.
type ChownFile interface {
	File

	// Chown changes the numeric uid and gid of the file.
	// A uid or gid of -1 means to not change that value.
	Chown(uid, gid int) error
}

// ChownFS is a file system that supports the Chown function.
type ChownFS interface {
	FS

	// Chown changes the numeric uid and gid of the named file.
	// If the file is a symbolic link, it changes the uid and gid of the link's target.
	// A uid or gid of -1 means to not change that value.
	Chown(name string, uid, gid int) error
}

// Chown changes the numeric uid and gid of the named file.
// If the file is a symbolic link, it changes the uid and gid of the link's target.
// A uid or gid of -1 means to not change that value.
func Chown(fsys FS, name string, uid, gid int) (err error) {
	if fsys, ok := fsys.(ChownFS); ok {
		return fsys.Chown(name, uid, gid)
	}

	// Open the file and attempt to call chown on it.
	file, err := fsys.Open(name)
	if err != nil {
		return err
	}
	defer safeClose(file, &err)

	if file, ok := file.(ChownFile); ok {
		return file.Chown(uid, gid)
	}

	return &PathError{Op: "chown", Path: name, Err: ErrNotSupported}
}

// LchownFS is a file system that supports the Lchown function.
type LchownFS interface {
	// Lchown changes the numeric uid and gid of the named file.
	// If the file is a symbolic link, it changes the uid and gid of the link itself.
	Lchown(name string, uid, gid int) error
}

// Lchown changes the numeric uid and gid of the named file.
// If the file is a symbolic link, it changes the uid and gid of the link itself.
func Lchown(fsys FS, name string, uid, gid int) (err error) {
	if fsys, ok := fsys.(LchownFS); ok {
		return fsys.Lchown(name, uid, gid)
	}
	return &PathError{Op: "chown", Path: name, Err: ErrNotSupported}
}
