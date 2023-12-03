package gateway

type RespOfModel interface {
	GetErrNo() int32
	GetErrTips() string
	String() string
}
