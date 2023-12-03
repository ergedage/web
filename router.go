// Code generated by hertztool
package main

import (
	"code.byted.org/middleware/hertz/pkg/app"
	"code.byted.org/middleware/hertz/pkg/app/server"
	"context"
	"github.com/ergedage/web/biz/handler"
	"github.com/ergedage/web/gateway"
	"github.com/ergedage/web/pb_gen/message"
)

func register(r *server.Hertz) {
	{
		_zhong := r.Group("zhong")
		{
			_biz := _zhong.Group("biz")
			{
				_get_project := _biz.Group("get_project")

				{
					_get_project.POSTEX("v1", func(ctx context.Context, c *app.RequestContext) {
						gateway.WrapHandler(ctx, c, &handler.GetProjectHandler{}, &message.RespOfGetProject{})
					}, "get_project")
				}
			}
		}

	}

}