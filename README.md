# WRFS

Package `wrfs` implements several [extension interfaces][iofs] for the `io/fs` package that make it possible to create and modify files.
Each extension interface has one or more helper functions that make use of it.
The helper functions accept any FS implementation, but return `ErrUnsupported` if the FS cannot support the requested operation.

`wrfs` re-exports the types and functions found in the `io/fs` package and can thus be used as a drop-in replacement for the `io/fs` package.

[iofs]: https://go.googlesource.com/proposal/+/master/design/draft-iofs.md#extension-interfaces-and-the-extension-pattern
