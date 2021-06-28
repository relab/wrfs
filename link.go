package wrfs

import "os"

// LinkFS is a file system that supports the Link function.
type LinkFS interface {
	// Link creates newname as a hard link to the oldname file.
	Link(oldname, newname string) error
}

// Link creates newname as a hard link to the oldname file.
func Link(fsys FS, oldname, newname string) error {
	if fsys, ok := fsys.(LinkFS); ok {
		return fsys.Link(oldname, newname)
	}
	return &os.LinkError{Op: "link", Old: oldname, New: newname, Err: ErrUnsupported}
}
