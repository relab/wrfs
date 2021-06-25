// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wrfs

import (
	"errors"
	"path"
	"time"
)

// Sub returns an FS corresponding to the subtree rooted at fsys's dir.
//
// If fs implements SubFS, Sub calls returns fsys.Sub(dir).
// Otherwise, if dir is ".", Sub returns fsys unchanged.
// Otherwise, Sub returns a new FS implementation sub that,
// in effect, implements sub.Open(dir) as fsys.Open(path.Join(dir, name)).
// The implementation also translates calls to ReadDir, ReadFile, and Glob appropriately.
//
// Note that Sub(os.DirFS("/"), "prefix") is equivalent to os.DirFS("/prefix")
// and that neither of them guarantees to avoid operating system
// accesses outside "/prefix", because the implementation of os.DirFS
// does not check for symbolic links inside "/prefix" that point to
// other directories. That is, os.DirFS is not a general substitute for a
// chroot-style security mechanism, and Sub does not change that fact.
func Sub(fsys FS, dir string) (FS, error) {
	if dir == "" || dir == "." {
		return fsys, nil
	}
	if !ValidPath(dir) {
		return nil, &PathError{Op: "sub", Path: dir, Err: errors.New("invalid name")}
	}
	if fsys, ok := fsys.(SubFS); ok {
		return fsys.Sub(dir)
	}
	return &subFS{fsys, dir}, nil
}

type subFS struct {
	fsys FS
	dir  string
}

// fullName maps name to the fully-qualified name dir/name.
func (f *subFS) fullName(op string, name string) (string, error) {
	if !ValidPath(name) {
		return "", &PathError{Op: op, Path: name, Err: errors.New("invalid name")}
	}
	return path.Join(f.dir, name), nil
}

// shorten maps name, which should start with f.dir, back to the suffix after f.dir.
func (f *subFS) shorten(name string) (rel string, ok bool) {
	if name == f.dir {
		return ".", true
	}
	if len(name) >= len(f.dir)+2 && name[len(f.dir)] == '/' && name[:len(f.dir)] == f.dir {
		return name[len(f.dir)+1:], true
	}
	return "", false
}

// fixErr shortens any reported names in PathErrors by stripping dir.
func (f *subFS) fixErr(err error) error {
	if e, ok := err.(*PathError); ok {
		if short, ok := f.shorten(e.Path); ok {
			e.Path = short
		}
	}
	return err
}

func (f *subFS) Open(name string) (File, error) {
	full, err := f.fullName("open", name)
	if err != nil {
		return nil, err
	}
	file, err := f.fsys.Open(full)
	return file, f.fixErr(err)
}

func (f *subFS) Stat(name string) (FileInfo, error) {
	full, err := f.fullName("stat", name)
	if err != nil {
		return nil, err
	}
	fi, err := Stat(f.fsys, full)
	return fi, f.fixErr(err)
}

func (f *subFS) Lstat(name string) (FileInfo, error) {
	full, err := f.fullName("lstat", name)
	if err != nil {
		return nil, err
	}
	fi, err := Lstat(f.fsys, full)
	return fi, f.fixErr(err)
}

func (f *subFS) ReadDir(name string) ([]DirEntry, error) {
	full, err := f.fullName("read", name)
	if err != nil {
		return nil, err
	}
	dir, err := ReadDir(f.fsys, full)
	return dir, f.fixErr(err)
}

func (f *subFS) ReadFile(name string) ([]byte, error) {
	full, err := f.fullName("read", name)
	if err != nil {
		return nil, err
	}
	data, err := ReadFile(f.fsys, full)
	return data, f.fixErr(err)
}

func (f *subFS) Glob(pattern string) ([]string, error) {
	// Check pattern is well-formed.
	if _, err := path.Match(pattern, ""); err != nil {
		return nil, err
	}
	if pattern == "." {
		return []string{"."}, nil
	}

	full := f.dir + "/" + pattern
	list, err := Glob(f.fsys, full)
	for i, name := range list {
		name, ok := f.shorten(name)
		if !ok {
			return nil, errors.New("invalid result from inner fsys Glob: " + name + " not in " + f.dir) // can't use fmt in this package
		}
		list[i] = name
	}
	return list, f.fixErr(err)
}

func (f *subFS) OpenFile(name string, flag int, perm FileMode) (File, error) {
	full, err := f.fullName("open", name)
	if err != nil {
		return nil, err
	}
	file, err := OpenFile(f.fsys, full, flag, perm)
	return file, f.fixErr(err)
}

func (f *subFS) Chmod(name string, mode FileMode) error {
	return f.permAction(name, mode, "chmod", Chmod)
}

func (f *subFS) Chown(name string, uid, gid int) error {
	return f.pathAction(name, "chown", func(fsys FS, path string) error {
		return Chown(fsys, path, uid, gid)
	})
}

func (f *subFS) Lchown(name string, uid, gid int) error {
	return f.pathAction(name, "lchown", func(fsys FS, path string) error {
		return Lchown(fsys, path, uid, gid)
	})
}

func (f *subFS) Chtimes(name string, atime, mtime time.Time) error {
	return f.pathAction(name, "chtimes", func(fsys FS, path string) error {
		return Chtimes(fsys, path, atime, mtime)
	})
}

func (f *subFS) Mkdir(name string, perm FileMode) error {
	return f.permAction(name, perm, "mkdir", Mkdir)
}

func (f *subFS) MkdirAll(path string, perm FileMode) error {
	return f.permAction(path, perm, "mkdir", MkdirAll)
}

func (f *subFS) Readlink(name string) (string, error) {
	full, err := f.fullName("readlink", name)
	if err != nil {
		return "", err
	}
	link, err := Readlink(f.fsys, full)
	if err != nil {
		return "", err
	}
	if link, ok := f.shorten(link); ok {
		return link, nil
	}
	return link, nil
}

func (f *subFS) Remove(name string) error {
	return f.pathAction(name, "remove", Remove)
}

func (f *subFS) RemoveAll(name string) error {
	return f.pathAction(name, "remove", RemoveAll)
}

func (f *subFS) Rename(oldname, newname string) error {
	return f.linkAction(oldname, newname, "rename", Rename)
}

func (f *subFS) SameFile(fi1, fi2 FileInfo) bool {
	return SameFile(f.fsys, fi1, fi2)
}

func (f *subFS) Symlink(oldname, newname string) error {
	return f.linkAction(oldname, newname, "symlink", Symlink)
}

func (f *subFS) Link(oldname, newname string) error {
	return f.linkAction(oldname, newname, "link", Link)
}

func (f *subFS) Truncate(name string, size int64) error {
	return f.pathAction(name, "truncate", func(fsys FS, path string) error {
		return Truncate(fsys, path, size)
	})
}

func (f *subFS) pathAction(path string, name string, action func(fsys FS, path string) error) error {
	full, err := f.fullName(name, path)
	if err != nil {
		return err
	}
	return f.fixErr(action(f.fsys, full))
}

func (f *subFS) permAction(path string, perm FileMode, name string, action func(fsys FS, path string, perm FileMode) error) error {
	return f.pathAction(path, name, func(fsys FS, path string) error {
		return action(fsys, path, perm)
	})
}

func (f *subFS) linkAction(oldPath, newPath string, name string, action func(fsys FS, src string, dest string) error) error {
	oldFull, err := f.fullName(name, oldPath)
	if err != nil {
		return err
	}
	newFull, err := f.fullName(name, newPath)
	if err != nil {
		return err
	}
	return f.fixErr(action(f.fsys, oldFull, newFull))
}
