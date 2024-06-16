package images

import (
	"beyimtech-test/internal/transport/http/api"

	"github.com/valyala/fasthttp"
)

func (c *Controller) listImagesHandler(ctx *fasthttp.RequestCtx) {
	var req listImagesRequest
	if err := api.GetData(ctx, &req); err != nil {
		api.SendError(ctx, err)
		return
	}

	images, err := c.imagesService.GetImages(ctx, req.limit, req.offset)
	if err != nil {
		api.SendError(ctx, err)
		return
	}

	api.SendData(ctx, listImagesResponse{
		Response: api.NewEmptySuccessResponse(),
		Images:   images,
	})
}
