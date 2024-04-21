package handlers

import "io"

type Handler interface {
	Handle(conn io.ReadWriteCloser) error
}
