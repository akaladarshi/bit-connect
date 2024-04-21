package handlers

import (
	"net"
)

type Handler interface {
	Handle(conn net.Conn) error
}
