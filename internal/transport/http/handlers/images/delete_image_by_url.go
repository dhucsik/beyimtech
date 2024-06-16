package images

import (
	"beyimtech-test/internal/transport/http/api"

	"github.com/valyala/fasthttp"
)

func (c *Controller) deleteImagesByURLHandler(ctx *fasthttp.RequestCtx) {
	var req deleteByURLReq
	if err := api.GetData(ctx, &req); err != nil {
		api.SendError(ctx, err)
		return
	}

	err := c.imagesService.DeleteImageByURL(ctx, req.URL)
	if err != nil {
		api.SendError(ctx, err)
		return
	}

	api.SendData(ctx, api.NewEmptySuccessResponse())
}
