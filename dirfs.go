package wrfs

import (
	"os"
	"runtime"
)

// DirFS returns a file system (an fs.FS) for the tree of files rooted at the directory dir.
//
// Note that DirFS("/prefix") only guarantees that the Open calls it makes to the
// operating system will begin with "/prefix": DirFS("/prefix").Open("file") is the
// same as os.Open("/prefix/file"). So if /prefix/file is a symbolic link pointing outside
// the /prefix tree, then using DirFS does not stop the access any more than using
// os.Open does. DirFS is therefore not a general substitute for a chroot-style security
// mechanism when the directory tree contains arbitrary content.
func DirFS(dir string) FS {
	return dirFS(dir)
}

func containsAny(s, chars string) bool {
	for i := 0; i < len(s); i++ {
		for j := 0; j < len(chars); j++ {
			if s[i] == chars[j] {
				return true
			}
		}
	}
	return false
}

type dirFS string

func (dir dirFS) validPath(name string) (string, error) {
	// On Windows, we will not allow back slashes or colons.
	if !ValidPath(name) || runtime.GOOS == "windows" && containsAny(name, `\:`) {
		return "", &PathError{Op: "open", Path: name, Err: ErrInvalid}
	}
	return string(dir) + "/" + name, nil
}

func (dir dirFS) Mkdir(path string, perm FileMode) error {
	full, err := dir.validPath(path)
	if err != nil {
		return err
	}
	return os.Mkdir(full, perm)
}

func (dir dirFS) Open(name string) (File, error) {
	full, err := dir.validPath(name)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(full)
	if err != nil {
		return nil, err // nil fs.File
	}
	return f, nil
}

func (dir dirFS) OpenFile(name string, flag int, perm FileMode) (File, error) {
	full, err := dir.validPath(name)
	if err != nil {
		return nil, err
	}
	file, err := os.OpenFile(full, flag, perm)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (dir dirFS) Remove(name string) error {
	full, err := dir.validPath(name)
	if err != nil {
		return err
	}
	return os.Remove(full)
}

func (dir dirFS) RemoveAll(path string) error {
	full, err := dir.validPath(path)
	if err != nil {
		return err
	}
	return os.RemoveAll(full)
}

func (dir dirFS) Rename(oldpath, newpath string) error {
	oldfull, err := dir.validPath(oldpath)
	if err != nil {
		return err
	}
	newfull, err := dir.validPath(newpath)
	if err != nil {
		return err
	}
	return os.Rename(oldfull, newfull)
}

func (dir dirFS) Symlink(oldname, newname string) error {
	oldfull, err := dir.validPath(oldname)
	if err != nil {
		return err
	}
	newfull, err := dir.validPath(newname)
	if err != nil {
		return err
	}
	return os.Symlink(oldfull, newfull)
}

func (dir dirFS) Truncate(name string, size int64) error {
	full, err := dir.validPath(name)
	if err != nil {
		return err
	}
	return os.Truncate(full, size)
}
