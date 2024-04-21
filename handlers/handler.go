package handlers

import (
	"net"
)

// Handler represents a handler
type Handler interface {
	Handle(conn net.Conn) error
}
