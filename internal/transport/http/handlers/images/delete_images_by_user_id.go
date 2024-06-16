package images

import (
	"beyimtech-test/internal/transport/http/api"

	"github.com/valyala/fasthttp"
)

func (c *Controller) deleteImagesByUserIDHandler(ctx *fasthttp.RequestCtx) {
	userID := ctx.UserValue("user_id").(string)

	err := c.imagesService.DeleteImagesByUserID(ctx, userID)
	if err != nil {
		api.SendError(ctx, err)
		return
	}

	api.SendData(ctx, api.NewEmptySuccessResponse())
}
