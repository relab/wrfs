package wrfs

// Chmod is a file with a Chmod method.
type ChmodFile interface {
	File

	// Chmod changes the mode of the file to mode.
	Chmod(mode FileMode) error
}

// ChmodFS is a file system with a Chmod method.
type ChmodFS interface {
	FS

	// Chmod changes the mode of the named file to mode.
	// If the file is a symbolic link, it changes the mode of the link's target.
	Chmod(name string, mode FileMode) error
}

// Chmod changes the mode of the named file to mode.
// If the file is a symbolic link, it changes the mode of the link's target.
func Chmod(fsys FS, name string, mode FileMode) error {
	if fsys, ok := fsys.(ChmodFS); ok {
		return fsys.Chmod(name, mode)
	}

	// Open the file and attempt to call chmod on it.
	file, err := fsys.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()

	if file, ok := file.(ChmodFile); ok {
		return file.Chmod(mode)
	}

	return &PathError{Op: "chmod", Path: name, Err: ErrNotSupported}
}
