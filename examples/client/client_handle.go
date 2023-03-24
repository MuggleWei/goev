package main

import (
	"encoding/json"
	"net"
	"time"

	goev "github.com/MuggleWei/goev"
	demo "github.com/MuggleWei/goev/examples/demo"
	log "github.com/sirupsen/logrus"
)

type ClientHandle struct {
	dispatcher *demo.Dispatcher
	evloop     *goev.Evloop
}

func NewClientHandle() *ClientHandle {
	handle := &ClientHandle{
		dispatcher: demo.NewDispatcher(),
	}

	handle.dispatcher.RegisterCallback(demo.MSG_JSON_PONG, handle.OnPong)

	return handle
}

func (this *ClientHandle) OnAddSession(session goev.Session) {
	userData, _ := session.GetUserData().(*ClientUserData)
	log.Infof("connected, addr=%v", userData.servAddr)
}
func (this *ClientHandle) OnClose(session goev.Session, err error) {
	userData := session.GetUserData().(*ClientUserData)
	log.Infof("disconnected, addr=%v", userData.servAddr)
	go func(servAddr string) {
		for {
			time.Sleep(3 * time.Second)
			log.Infof("try reconnect")

			conn, err := net.Dial("tcp", servAddr)
			if err != nil {
				log.Errorf("failed dail addr: %v", servAddr)
				continue
			}
			log.Infof("success dial: %v", servAddr)

			session := &ClientSession{
				conn: conn,
				codec: &demo.BytesCodec{
					MaxPayloadLimit: 512 * 1024,
				},
				userData: &ClientUserData{
					servAddr: servAddr,
				},
			}
			this.evloop.AddSession(session)

			break
		}
	}(userData.servAddr)
}
func (this *ClientHandle) onMessage(session goev.Session, hdr interface{}, payload []byte) {
	this.dispatcher.HandleMessage(session, hdr, payload)
}
func (this *ClientHandle) onTimer(sessions map[goev.Session]time.Time) {
	for session := range sessions {
		ts := time.Now()
		req := &demo.MsgPing{
			Sec:  uint64(ts.Unix()),
			NSec: uint32(ts.Nanosecond()),
		}

		log.Debugf("send msg ping: %+v", req)
		demo.SessionWrite(session, demo.MSG_JSON_PING, req)
	}
}

func (this *ClientHandle) OnPong(session goev.Session, payload []byte) {
	var req demo.MsgPong
	err := json.Unmarshal(payload, &req)
	if err != nil {
		log.Errorf("failed unmarshal json: %v", err)
		conn := session.GetConn()
		conn.Close()
		return
	}
	log.Debugf("recv msg pong: %+v", req)
}
