package images

import (
	"beyimtech-test/internal/services/images"

	"github.com/fasthttp/router"
)

type Controller struct {
	imagesService images.Service
}

func New(
	imagesService images.Service,
) *Controller {
	return &Controller{
		imagesService: imagesService,
	}
}

func (c *Controller) Init(r *router.Router) {
	group := r.Group("/api/v1")

	group.POST("/images", c.uploadImageHandler)

	group.GET("/users/{user_id}/images", c.listImagesByUserIDHandler)
	group.GET("/images", c.listImagesHandler)

	group.DELETE("/users/{user_id}/images", c.deleteImagesByUserIDHandler)
	group.DELETE("/images", c.deleteImagesByURLHandler)
}
