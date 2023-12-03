package gateway

import (
	"code.byted.org/gopkg/logs"
	"context"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"

	"google.golang.org/protobuf/proto"

	"code.byted.org/bytelingo/goutil/log"
	"code.byted.org/middleware/hertz/pkg/app"
	"code.byted.org/middleware/hertz_ext/v2/binding"
)

// MarshalOptions ref https://pkg.go.dev/google.golang.org/protobuf/encoding/protojson
var MarshalOptions = protojson.MarshalOptions{
	UseProtoNames:   true,
	UseEnumNumbers:  true,
	EmitUnpopulated: true,
}

// WrapHandler 接口通用接入层，维护非具体业务功能的公共逻辑
func WrapHandler(ctx context.Context, request *app.RequestContext, hdl Handler, reqModel proto.Message) {
	requester := createRequester(ctx, request)
	if err := binding.BindAndValidate(requester.request, reqModel); err != nil {
		log.ErrorWithKey(ctx, "frame parse req fail", "err", err)
	}

	respModel := hdl.PreProcess(requester, reqModel)
	if respModel != nil {
		log.WarnWithKey(ctx, "handler pre-process failed")
	} else {
		respModel = hdl.Process(requester, reqModel)
	}
	hdl.PostProcess(requester, reqModel, respModel)
	logs.CtxInfo(ctx, "header=%+v, query=%+v", requester.GetHeaders(), requester.GetQueries())

	switch requester.getBodyCodec() {
	case bodyProtobuf:
		request.ProtoBuf(http.StatusOK, respModel)
	case bodyJSON:
		// 前端处理int64会有精度丢失问题，约定特殊header，序列化时，按照string处理
		if requester.GetUseJsonPb() {
			bytes, _ := MarshalOptions.Marshal(respModel)
			request.Data(http.StatusOK, "application/json", bytes)
		} else {
			request.JSON(http.StatusOK, respModel)
		}
	default:
		log.WarnWithKey(ctx, "unexpected content-type from request")
		if requester.GetUseJsonPb() {
			bytes, _ := MarshalOptions.Marshal(respModel)
			request.Data(http.StatusOK, "", bytes)
		} else {
			request.JSON(http.StatusOK, respModel)
		}
	}
}

type Handler interface {
	// PreProcess 前置处理，包含但不限于参数校验。
	// - responseModel 为 nil 时，进入Process；
	// - 否则，会作为结果直接返回；
	PreProcess(requester *Requester, requestModel proto.Message) (responseModel proto.Message)
	// Process 业务逻辑。
	Process(requester *Requester, requestModel proto.Message) (responseModel proto.Message)
	// PostProcess 后置处理，但不允许在这里对 responseModel 再进行写操作；
	// 框架保证无论 PreProcess 和 Process 结果如何，都会进入到 PostProcess。
	PostProcess(requester *Requester, requestModel proto.Message, responseModel proto.Message)
}
