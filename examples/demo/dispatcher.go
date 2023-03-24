package demo

import (
	"errors"
	"fmt"

	"github.com/MuggleWei/goev"
)

type DemoCallback func(goev.Session, []byte)

type Dispatcher struct {
	callbacks map[uint32]DemoCallback
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		callbacks: make(map[uint32]DemoCallback),
	}
}

func (this *Dispatcher) RegisterCallback(msgId uint32, cb DemoCallback) error {
	_, found := this.callbacks[msgId]
	if found {
		return errors.New(
			fmt.Sprintf("repeated register callback with message id: %v", msgId))
	}
	this.callbacks[msgId] = cb
	return nil
}

func (this *Dispatcher) HandleMessage(session goev.Session, hdr interface{}, payload []byte) error {
	msgHdr, ok := hdr.(*MsgHdr)
	if !ok {
		return errors.New("invalid head")
	}

	cb, found := this.callbacks[msgHdr.MsgId]
	if !found {
		return errors.New(fmt.Sprintf("msg_id[%v] not register", msgHdr.MsgId))
	}
	cb(session, payload)

	return nil
}
