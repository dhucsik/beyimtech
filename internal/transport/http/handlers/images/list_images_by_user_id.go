package images

import (
	"beyimtech-test/internal/transport/http/api"

	"github.com/valyala/fasthttp"
)

func (c *Controller) listImagesByUserIDHandler(ctx *fasthttp.RequestCtx) {
	userID := ctx.UserValue("user_id").(string)

	images, err := c.imagesService.GetImagesByUserID(ctx, userID)
	if err != nil {
		api.SendError(ctx, err)
		return
	}

	api.SendData(ctx, listImagesResponse{
		Response: api.NewEmptySuccessResponse(),
		Images:   images,
	})
}
