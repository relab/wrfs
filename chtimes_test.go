package wrfs_test

import (
	"testing"
	"time"

	. "github.com/Raytar/wrfs"
	"github.com/Raytar/wrfs/wrfstest"
)

type chtimesOnly struct {
	ChtimesFS
}

func (chtimesOnly) Open(name string) (File, error) {
	return nil, &PathError{Op: "chtimes", Path: name, Err: ErrNotSupported}
}

func TestChtimes(t *testing.T) {
	mapfs := testFS()
	check := func(fsys FS, atime, mtime time.Time) {
		err := Chtimes(fsys, "hello.txt", atime, mtime)
		if err != nil {
			t.Fatal(err)
		}
		if got := mapfs["hello.txt"].Sys.(*wrfstest.MapFileSys).ATime; !got.Equal(atime) {
			t.Errorf("Wrong atime: got: %v, want: %v", got, atime)
		}
		if got := mapfs["hello.txt"].ModTime; !got.Equal(mtime) {
			t.Errorf("Wrong mtime: got: %v, want: %v", got, mtime)
		}
	}
	check(chtimesOnly{mapfs}, time.Now().Add(time.Hour), time.Now().Add(time.Minute))
	check(openOnly{mapfs}, time.Now().Add(time.Minute), time.Now().Add(time.Hour))
}
