package handler

import (
	"code.byted.org/bytelingo/goutil/errcode/lingo_err"
	"github.com/ergedage/web/pb_gen/message"
	"github.com/ergedage/web/util"

	"github.com/ergedage/web/gateway"
	"google.golang.org/protobuf/proto"
)

type GetProjectHandler struct {
}

func (h *GetProjectHandler) PreProcess(requester *gateway.Requester, requestModel proto.Message) (responseModel proto.Message) {
	request := requestModel.(*message.ReqOfGetProject)
	if request.GetProjectId() <= 0 {
		return util.FillResponse(&message.RespOfGetProject{}, nil, lingo_err.ParamErr)
	}

	return nil
}

func (h *GetProjectHandler) Process(requester *gateway.Requester, requestModel proto.Message) (responseModel proto.Message) {
	//request := requestModel.(*message.ReqOfGetProject)
	responseModel = &message.RespOfGetProject{}

	return util.FillResponse(&message.RespOfGetProject{
		ProjectName: "new_project",
		UrlLink:     "link",
		Img:         "image",
	}, nil, nil)
}

func (h *GetProjectHandler) PostProcess(requester *gateway.Requester, requestModel proto.Message, responseModel proto.Message) {
	return
}
