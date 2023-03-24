package goev

import "bufio"

type Codec interface {
	// Encode function
	// :param msgId: message ID
	// :param payload: payload bytes
	// :return
	//   - []byte: encode result
	//   - error: error message
	Encode(msdId interface{}, payload []byte) ([]byte, error)

	// Decode function
	// :param reader: reader
	// :return
	//   - interface{}: message header
	//   - []byte: message payload
	//   - error: error message
	Decode(reader *bufio.Reader) (interface{}, []byte, error)
}
