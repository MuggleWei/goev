package main

import (
	"net"
	"time"

	goev "github.com/MuggleWei/goev"
)

type ServerUserData struct {
	lastUpdateTime time.Time
	remoteAddr     string
}

type ServerSession struct {
	conn     net.Conn
	codec    goev.Codec
	userData *ServerUserData
}

func (this *ServerSession) GetConn() net.Conn {
	return this.conn
}
func (this *ServerSession) GetUserData() interface{} {
	return this.userData
}
func (this *ServerSession) GetCodec() goev.Codec {
	return this.codec
}
