package types

import "io"

type WriteCloserWithError interface {
	io.WriteCloser
	CloseWithError(e error) error
}
