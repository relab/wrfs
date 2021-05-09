package wrfs_test

import (
	"testing"

	. "github.com/Raytar/wrfs"
	"github.com/Raytar/wrfs/wrfstest"
)

type chownOnly struct {
	ChownFS
}

func (chownOnly) Open(name string) (File, error) {
	return nil, &PathError{Op: "open", Path: name, Err: ErrNotSupported}
}

func TestChown(t *testing.T) {
	mapfs := testFS()
	check := func(fsys FS, uid, gid int) {
		err := Chown(fsys, "hello.txt", uid, gid)
		if err != nil {
			t.Fatal(err)
		}
		sys := mapfs["hello.txt"].Sys.(*wrfstest.MapFileSys)
		if sys.Gid != gid {
			t.Errorf("wrong gid: got: %d, want: %d", sys.Gid, gid)
		}
		if sys.Uid != uid {
			t.Errorf("wrong uid: got: %d, want: %d", sys.Uid, uid)
		}
	}
	check(chownOnly{mapfs}, 2, 3)
	check(openOnly{mapfs}, 3, 2)
}
