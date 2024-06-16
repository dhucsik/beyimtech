package images

import (
	"beyimtech-test/internal/entity"
	"beyimtech-test/internal/transport/http/api"
	"encoding/json"
	"log"
	"mime/multipart"
	"strconv"

	"github.com/valyala/fasthttp"
)

type uploadImageRequest struct {
	file   *multipart.FileHeader
	userID string
}

func (r *uploadImageRequest) GetParameters(ctx *fasthttp.RequestCtx) error {
	userID := ctx.FormValue("user_id")
	file, err := ctx.FormFile("file")
	if err != nil {
		return err
	}

	r.file = file
	r.userID = string(userID)

	return nil
}

type listImagesRequest struct {
	limit  int
	offset int
}

func (r *listImagesRequest) GetParameters(ctx *fasthttp.RequestCtx) error {
	limitStr := ctx.QueryArgs().Peek("limit")
	limit, err := strconv.Atoi(string(limitStr))
	if err != nil {
		log.Println(err)
		limit = 20
	}

	offsetStr := ctx.QueryArgs().Peek("offset")
	offset, err := strconv.Atoi(string(offsetStr))
	if err != nil {
		log.Println(err)
		offset = 0
	}

	r.limit = limit
	r.offset = offset
	return nil
}

type deleteByURLReq struct {
	URL string `json:"url"`
}

func (r *deleteByURLReq) GetParameters(ctx *fasthttp.RequestCtx) error {
	return json.Unmarshal(ctx.Request.Body(), r)
}

type listImagesResponse struct {
	api.Response
	Images []*entity.Image `json:"images"`
}
