// Package wrfs provides extension interfaces for the io/fs package that enable creating, updating and removing files
// and directories. wrfs can work as a drop-in replacement for the io/fs package, as it re-exports the types and
// functions defined by io/fs.
package wrfs

import (
	"errors"
	"io"
)

// TODO: replace this with errors.ErrUnsupported: https://github.com/golang/go/issues/41198

var ErrUnsupported = errors.New("unsupported operation")

// safeClose closes an io.Closer and stores the error in errPtr
func safeClose(closer io.Closer, errPtr *error) {
	err := closer.Close()
	// don't overwrite an existing error
	if *errPtr == nil {
		*errPtr = err
	}
}
