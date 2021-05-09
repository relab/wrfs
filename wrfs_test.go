package wrfs_test

import (
	"time"

	. "github.com/Raytar/wrfs"
	"github.com/Raytar/wrfs/wrfstest"
)

type openOnly struct {
	FS
}

func testFS() wrfstest.MapFS {
	return wrfstest.MapFS{
		"hello.txt": &wrfstest.MapFile{
			Data:    []byte("hello world"),
			Mode:    0666,
			ModTime: time.Now(),
			Sys: &wrfstest.MapFileSys{
				ATime: time.Now(),
				Uid:   1,
				Gid:   1,
			},
		},
		"subdir/goodbye.txt": &wrfstest.MapFile{
			Data:    []byte("goodbye"),
			Mode:    0666,
			ModTime: time.Now(),
			Sys: &wrfstest.MapFileSys{
				ATime: time.Now(),
				Uid:   1,
				Gid:   1,
			},
		},
	}
}
