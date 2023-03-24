package main

import (
	"net"

	goev "github.com/MuggleWei/goev"
)

type ClientUserData struct {
	servAddr string
}

type ClientSession struct {
	conn     net.Conn
	codec    goev.Codec
	userData *ClientUserData
}

func (this *ClientSession) GetConn() net.Conn {
	return this.conn
}
func (this *ClientSession) GetUserData() interface{} {
	return this.userData
}
func (this *ClientSession) GetCodec() goev.Codec {
	return this.codec
}
