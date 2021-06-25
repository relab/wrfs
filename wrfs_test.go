//go:build aix || darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris
// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris

package wrfs_test

import (
	"os"
	"testing"
	"time"

	. "github.com/Raytar/wrfs"
)

func TestChmod(t *testing.T) {
	fsys := getFS(t)

	testCase := func(fsys FS) {
		fileName := "TestChmod"
		want := FileMode(0750)
		newFile(t, fsys, fileName)

		err := Chmod(fsys, fileName, want)
		check(t, err)

		checkMode(t, fsys, fileName, want)
	}

	t.Run("", func(t *testing.T) { testCase(fsys) })
	t.Run("OpenFileOnly", func(t *testing.T) { testCase(openFileOnly{fsys.(OpenFileFS)}) })
}

func TestChtimes(t *testing.T) {
	fsys := getFS(t)
	fileName := "TestChtimes"
	newFile(t, fsys, fileName)
	before, err := Stat(fsys, fileName)
	check(t, err)

	// only looking at mtime for now, as some systems may not support atime.
	err = Chtimes(fsys, fileName, time.Time{}, before.ModTime().Add(-time.Second))
	check(t, err)

	after, err := Stat(fsys, fileName)
	check(t, err)

	if after.ModTime().After(before.ModTime()) {
		t.Fatalf("ModTime didn't go backwards: was: %v, after: %v", before.ModTime(), after.ModTime())
	}
}

func TestMkdirAll(t *testing.T) {
	testCase := func(fsys FS) {
		dirName := "TestMkdirAll/foo/bar"
		mode := FileMode(0755)

		err := MkdirAll(fsys, dirName, mode)
		check(t, err)

		fi, err := Stat(fsys, dirName)
		check(t, err)

		if !fi.IsDir() {
			t.Error("Is not a directory")
		}

		if got := fi.Mode() & ModePerm; got != mode {
			t.Errorf("Wrong mode for directory: got: %v, want: %v", got, mode)
		}
	}
	fsys := getFS(t)
	t.Run("MkdirAllFS", func(*testing.T) {
		testCase(fsys)
	})
	t.Run("MkdirFS", func(*testing.T) {
		testCase(mkdirOnly{fsys.(MkdirFS)})
	})
}

type mkdirOnly struct {
	MkdirFS
}

func TestRemoveAll(t *testing.T) {
	testCase := func(fsys FS) {
		dirName := "TestRemoveAll"
		subDir := dirName + "/sub"
		file := subDir + "/leaf"
		perm := FileMode(0755)

		err := MkdirAll(fsys, subDir, perm)
		check(t, err)
		newFile(t, fsys, file)

		err = RemoveAll(fsys, dirName)
		check(t, err)

		_, err = Stat(fsys, dirName)

		if err == nil {
			t.Error("expected an error, but got nil")
		}
	}
	fsys := getFS(t)
	t.Run("RemoveAllFS", func(*testing.T) {
		testCase(fsys)
	})
	t.Run("RemoveFS", func(*testing.T) {
		testCase(removeOnly{fsys.(removeOnlyFS)})
	})
}

type removeOnlyFS interface {
	RemoveFS
	MkdirFS
	OpenFileFS
}

type removeOnly struct {
	removeOnlyFS
}

func TestRename(t *testing.T) {
	fsys := getFS(t)
	oldName := "TestRename"
	newName := "TestRename2"
	newFile(t, fsys, oldName)

	err := Rename(fsys, oldName, newName)
	check(t, err)

	_, err = Stat(fsys, newName)
	if err != nil {
		t.Error(err)
	}
}

func TestSymlink(t *testing.T) {
	fsys := getFS(t)
	src := "TestSymlink"
	dest := "TestSymlink2"
	newFile(t, fsys, src)

	err := Symlink(fsys, src, dest)
	check(t, err)

	fi, err := Lstat(fsys, dest)
	check(t, err)

	if fi.Mode()&ModeSymlink == 0 {
		t.Error("this does not look like a symlink")
	}

	link, err := Readlink(fsys, dest)
	check(t, err)

	if link != src {
		t.Errorf("got: %v, want: %v", link, src)
	}
}

func TestLink(t *testing.T) {
	fsys := getFS(t)
	src := "TestSymlink"
	dest := "TestSymlink2"
	newFile(t, fsys, src)

	fi1, err := Stat(fsys, src)
	check(t, err)

	err = Link(fsys, src, dest)
	check(t, err)

	fi2, err := Stat(fsys, dest)
	check(t, err)

	if !SameFile(fsys, fi1, fi2) {
		t.Error("SameFile returned false")
	}
}

func TestTruncate(t *testing.T) {
	fsys := getFS(t)

	testCase := func(fsys FS) {
		file := "TestTruncate"
		size := int64(100)

		newFile(t, fsys, file)

		err := Truncate(fsys, file, size)
		check(t, err)

		fi, err := Stat(fsys, file)
		check(t, err)

		if fi.Size() < size {
			t.Error("File size was not increased")
		}
	}

	t.Run("", func(t *testing.T) { testCase(fsys) })
	t.Run("OpenFileOnly", func(t *testing.T) { testCase(openFileOnly{fsys.(OpenFileFS)}) })
}

type openFileOnly struct {
	OpenFileFS
}

func getFS(t *testing.T) FS {
	dir := t.TempDir()
	dirFS := DirFS(dir)
	err := Mkdir(dirFS, "subfs", 0755)
	check(t, err)
	subFS, err := Sub(dirFS, "subfs")
	check(t, err)
	return subFS
}

func newFile(t *testing.T, fsys FS, fileName string) {
	file, err := OpenFile(fsys, fileName, os.O_RDONLY|os.O_CREATE|os.O_TRUNC, 0777)
	check(t, err)
	check(t, file.Close())
}

func checkMode(t *testing.T, fsys FS, fileName string, want FileMode) {
	fi, err := Stat(fsys, fileName)
	check(t, err)

	if fi.Mode() != want {
		t.Errorf("Wrong file mode: got: %v, want: %v", fi.Mode(), want)
	}
}

func check(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}
