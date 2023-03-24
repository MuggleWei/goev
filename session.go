package goev

import "net"

type Session interface {
	GetConn() net.Conn
	GetCodec() Codec
	GetUserData() interface{}
}
