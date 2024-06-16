package images

import (
	"beyimtech-test/internal/transport/http/api"

	"github.com/valyala/fasthttp"
)

func (c *Controller) uploadImageHandler(ctx *fasthttp.RequestCtx) {
	var req uploadImageRequest
	if err := api.GetData(ctx, &req); err != nil {
		api.SendError(ctx, err)
		return
	}

	c.imagesService.SendFile(ctx, req.userID, req.file)
	api.SendData(ctx, api.NewEmptySuccessResponse())
}
