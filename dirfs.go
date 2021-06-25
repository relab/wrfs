package wrfs

import (
	"os"
	"time"
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
	return &subFS{fsys: hostFS{}, dir: dir}
}

type hostFS struct{}

func (hostFS) Chmod(name string, mode FileMode) error {
	return os.Chmod(name, mode)
}

func (hostFS) Chown(name string, uid, gid int) error {
	return os.Chown(name, uid, gid)
}

func (hostFS) Chtimes(name string, atime, mtime time.Time) error {
	return os.Chtimes(name, atime, mtime)
}

func (hostFS) Mkdir(path string, perm FileMode) error {
	return os.Mkdir(path, perm)
}

func (hostFS) Open(name string) (File, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err // nil fs.File
	}
	return f, nil
}

func (hostFS) OpenFile(name string, flag int, perm FileMode) (File, error) {
	file, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (hostFS) Stat(name string) (FileInfo, error) {
	fi, err := os.Stat(name)
	if err != nil {
		return nil, err
	}
	return fi, nil
}

func (hostFS) Lstat(name string) (FileInfo, error) {
	fi, err := os.Lstat(name)
	if err != nil {
		return nil, err
	}
	return fi, nil
}

func (hostFS) Readlink(name string) (string, error) {
	link, err := os.Readlink(name)
	if err != nil {
		return "", err
	}
	return link, nil
}

func (hostFS) Remove(name string) error {
	return os.Remove(name)
}

func (hostFS) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (hostFS) Rename(oldpath, newpath string) error {
	return os.Rename(oldpath, newpath)
}

func (hostFS) SameFile(fi1, fi2 FileInfo) bool {
	return os.SameFile(fi1, fi2)
}

func (hostFS) Symlink(oldname, newname string) error {
	return os.Symlink(oldname, newname)
}

func (hostFS) Link(oldname, newname string) error {
	return os.Link(oldname, newname)
}

func (hostFS) Truncate(name string, size int64) error {
	return os.Truncate(name, size)
}
