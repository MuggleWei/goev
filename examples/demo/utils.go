package demo

import (
	"encoding/json"

	"github.com/MuggleWei/goev"
	log "github.com/sirupsen/logrus"
)

func SessionWrite(session goev.Session, msgId uint32, obj interface{}) error {
	payload, err := json.Marshal(obj)
	if err != nil {
		log.Errorf("failed json marshal: %v", err)
		return err
	}

	out, err := session.GetCodec().Encode(msgId, payload)
	if err != nil {
		log.Errorf("failed encode: %v", err)
		return err
	}

	_, err = session.GetConn().Write(out)
	if err != nil {
		log.Errorf("failed write")
		return err
	}

	return nil
}
