package wrfs

// SameFileFS is a file system that supports the SameFile function.
type SameFileFS interface {
	// SameFile reports whether fi1 and fi2 describe the same file.
	// For example, on Unix this means that the device and inode fields of the two underlying structures are identical;
	// on other systems the decision may be based on the path names.
	SameFile(fi1, fi2 FileInfo) bool
}

// SameFile reports whether fi1 and fi2 describe the same file.
// For example, on Unix this means that the device and inode fields of the two underlying structures are identical;
// on other systems the decision may be based on the path names.
func SameFile(fsys FS, fi1, fi2 FileInfo) bool {
	if fsys, ok := fsys.(SameFileFS); ok {
		return fsys.SameFile(fi1, fi2)
	}
	return false
}
