package main

import (
	"encoding/json"
	"time"

	goev "github.com/MuggleWei/goev"
	demo "github.com/MuggleWei/goev/examples/demo"
	log "github.com/sirupsen/logrus"
)

type ServerHandle struct {
	dispatcher *demo.Dispatcher
}

func NewServerHandle() *ServerHandle {
	handle := &ServerHandle{
		dispatcher: demo.NewDispatcher(),
	}

	handle.dispatcher.RegisterCallback(demo.MSG_JSON_PING, handle.OnPing)

	return handle
}

func (this *ServerHandle) OnAddSession(session goev.Session) {
	userData, _ := session.GetUserData().(*ServerUserData)

	log.Infof("connected, addr=%v", userData.remoteAddr)
	userData.lastUpdateTime = time.Now()
}
func (this *ServerHandle) OnClose(session goev.Session, err error) {
	userData, _ := session.GetUserData().(*ServerUserData)
	log.Infof("disconnected, addr=%v", userData.remoteAddr)
}
func (this *ServerHandle) onMessage(session goev.Session, hdr interface{}, payload []byte) {
	this.dispatcher.HandleMessage(session, hdr, payload)
}
func (this *ServerHandle) onTimer(sessions map[goev.Session]time.Time) {
	currTime := time.Now()
	for session := range sessions {
		userData, _ := session.GetUserData().(*ServerUserData)
		elapsed := currTime.Sub(userData.lastUpdateTime)
		if elapsed > time.Second*15 {
			conn := session.GetConn()
			conn.Close()
		}
	}
}

func (this *ServerHandle) OnPing(session goev.Session, payload []byte) {
	var req demo.MsgPing
	err := json.Unmarshal(payload, &req)
	if err != nil {
		log.Errorf("failed unmarshal json: %v", err)
		conn := session.GetConn()
		conn.Close()
		return
	}
	log.Debugf("recv msg ping: %+v", req)

	userData := session.GetUserData().(*ServerUserData)
	userData.lastUpdateTime = time.Now()

	log.Debugf("send msg pong: %+v", req)
	demo.SessionWrite(session, demo.MSG_JSON_PONG, &demo.MsgPong{
		Sec:  req.Sec,
		NSec: req.NSec,
	})
}
