package util

import (
	"code.byted.org/bytelingo/goutil/errcode"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"time"
)

const (
	RespErrNoName = "err_no"
	RespErrTips   = "err_tips"
	RespData      = "data"
	RespTS        = "ts"
)

func FillResponse(resp proto.Message, dataMsg proto.Message, err error) proto.Message {
	fields := resp.ProtoReflect().Descriptor().Fields()
	errNo := fields.ByName(RespErrNoName)
	lingoErr := errcode.OfError(err)
	if errNo != nil {
		resp.ProtoReflect().Set(errNo, protoreflect.ValueOf(lingoErr.Code()))
	}
	errMsg := fields.ByName(RespErrTips)
	if errMsg != nil {
		resp.ProtoReflect().Set(errMsg, protoreflect.ValueOf(lingoErr.Msg()))
	}
	if dataMsg != nil {
		data := fields.ByName(RespData)
		if data != nil {
			resp.ProtoReflect().Set(data, protoreflect.ValueOf(dataMsg.ProtoReflect()))
		}
	}
	ts := fields.ByName(RespTS)
	if ts != nil {
		resp.ProtoReflect().Set(ts, protoreflect.ValueOfInt64(time.Now().Unix()))
	}
	return resp
}
