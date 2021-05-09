package wrfs_test

import (
	"testing"

	. "github.com/Raytar/wrfs"
)

type mkdirOnly struct {
	MkdirFS
}

type mkdirAllOnly struct {
	MkdirAllFS
}

func TestMkdir(t *testing.T) {
	mapfs := testFS()

	err := Mkdir(mapfs, "newdir", 0777)
	if err != nil {
		t.Fatal(err)
	}
	info, err := Stat(mapfs, "newdir")
	if err != nil {
		t.Fatal(err)
	}

	if !info.IsDir() {
		t.Error("Is not a directory")
	}
	if perm := info.Mode() & ModePerm; perm != 0777 {
		t.Errorf("Wrong permissions: got: %#o ,want: %#o", perm, 0777)
	}
}

func TestMkdirAll(t *testing.T) {
	mapfs := testFS()
	checkDir := func(fsys FS, dir string, mode FileMode) {
		info, err := Stat(fsys, dir)
		if err != nil {
			t.Fatal(err)
		}
		if !info.IsDir() {
			t.Errorf("Expected %s to be a directory", dir)
		}
		if got := info.Mode() & ModePerm; got != mode {
			t.Errorf("Wrong permissions for directory: got: %#o, want: %#o", got, mode)
		}
	}
	check := func(fsys FS) {
		err := MkdirAll(fsys, "foo/bar/baz/.", 0775)
		if err != nil {
			t.Fatal(err)
		}
		checkDir(fsys, "foo", 0775)
		checkDir(fsys, "foo/bar", 0775)
		checkDir(fsys, "foo/bar/baz", 0775)
	}
	check(mkdirAllOnly{mapfs})
	// get a fresh FS
	mapfs = testFS()
	check(mkdirOnly{mapfs})
}
