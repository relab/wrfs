package wrfs_test

import (
	"testing"

	. "github.com/Raytar/wrfs"
)

type chmodOnly struct {
	ChmodFS
}

func (chmodOnly) Open(name string) (File, error) {
	return nil, &PathError{Op: "open", Path: name, Err: ErrNotSupported}
}

func TestChmod(t *testing.T) {
	mapfs := testFS()
	check := func(fsys FS, want FileMode) {
		err := Chmod(fsys, "hello.txt", want)
		if err != nil {
			t.Fatal(err)
		}
		if got := mapfs["hello.txt"].Mode & ModePerm; got != want {
			t.Errorf("Wrong permissions for file: got: %#o, want: %#o", got, want)
		}
	}
	check(chmodOnly{mapfs}, 0777)
	check(openOnly{mapfs}, 0444)
}
