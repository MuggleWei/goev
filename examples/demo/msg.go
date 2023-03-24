package demo

const (
	MSG_ID_NULL = iota
	MSG_JSON_PING
	MSG_JSON_PONG
	MSG_JSON_REQ_LOGIN
	MSG_JSON_RSP_LOGIN
	MSG_JSON_REQ_SUM
	MSG_JSON_RSP_SUM
)

type MsgPing struct {
	Sec  uint64 // second
	NSec uint32 // nano-second
}
type MsgPong struct {
	Sec  uint64 // second
	NSec uint32 // nano-second
}
type MsgReqLogin struct {
	User   string // user name
	Passwd string // password
}
type MsgRspLogin struct {
	ErrId  uint32 // 0 - success, otherwise - error id
	ErrMsg string // error message
}
