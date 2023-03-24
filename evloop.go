package goev

import (
	"bufio"
	"errors"
	"time"
)

const EV_MSG_TYPE_ATTACH = 1
const EV_MSG_TYPE_DETACH = 2
const EV_MSG_TYPE_MSG = 3

type EvloopMsg struct {
	evMsgType uint32
	session   Session
	hdr       interface{}
	payload   []byte
	err       error
}

type CallbackAddSession func(session Session)
type CallbackClose func(session Session, err error)
type CallbackMessage func(session Session, hdr interface{}, payload []byte)
type CallbackTimer func(sessions map[Session]time.Time)

type Evloop struct {
	sessions  map[Session]time.Time
	ch        chan *EvloopMsg
	tick      time.Duration
	onAddConn CallbackAddSession
	onClose   CallbackClose
	onMessage CallbackMessage
	onTimer   CallbackTimer
}

func NewEvloop() *Evloop {
	return &Evloop{
		sessions:  make(map[Session]time.Time),
		ch:        make(chan *EvloopMsg),
		tick:      5 * time.Second,
		onAddConn: nil,
		onClose:   nil,
		onMessage: nil,
		onTimer:   nil,
	}
}

func (this *Evloop) SetCallbackOnAddConn(cb CallbackAddSession) {
	this.onAddConn = cb
}
func (this *Evloop) SetCallbackOnClose(cb CallbackClose) {
	this.onClose = cb
}
func (this *Evloop) SetCallbackOnMessage(cb CallbackMessage) {
	this.onMessage = cb
}
func (this *Evloop) SetCallbackOnTimer(cb CallbackTimer) {
	this.onTimer = cb
}
func (this *Evloop) SetTimerTick(tick time.Duration) {
	this.tick = tick
}

func (this *Evloop) Run() {
	timer := time.Tick(this.tick)
	for {
		select {
		case msg := <-this.ch:
			switch msg.evMsgType {
			case EV_MSG_TYPE_ATTACH:
				this.sessions[msg.session] = time.Now()
				if this.onAddConn != nil {
					this.onAddConn(msg.session)
				}
			case EV_MSG_TYPE_DETACH:
				if this.onClose != nil {
					this.onClose(msg.session, msg.err)
				}
				delete(this.sessions, msg.session)
				msg.session.GetConn().Close()
			case EV_MSG_TYPE_MSG:
				if this.onMessage != nil {
					this.onMessage(msg.session, msg.hdr, msg.payload)
				}
			}
		case <-timer:
			if this.onTimer != nil {
				this.onTimer(this.sessions)
			}
		}
	}
}

func (this *Evloop) AddSession(session Session) error {
	if session.GetConn() == nil {
		return errors.New("no conn in session")
	}
	if session.GetCodec() == nil {
		return errors.New("no codec in session")
	}
	go this.handleSession(session)
	return nil
}

func (this *Evloop) handleSession(session Session) {
	this.ch <- &EvloopMsg{
		evMsgType: EV_MSG_TYPE_ATTACH,
		session:   session,
		hdr:       nil,
		payload:   nil,
		err:       nil,
	}

	reader := bufio.NewReader(session.GetConn())
	for {
		hdr, payload, err := session.GetCodec().Decode(reader)
		if err != nil {
			this.ch <- &EvloopMsg{
				evMsgType: EV_MSG_TYPE_DETACH,
				session:   session,
				hdr:       nil,
				payload:   nil,
				err:       err,
			}
			break
		}

		this.ch <- &EvloopMsg{
			evMsgType: EV_MSG_TYPE_MSG,
			session:   session,
			hdr:       hdr,
			payload:   payload,
			err:       nil,
		}
	}
}
