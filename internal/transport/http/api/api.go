package api

import (
	apiErrors "beyimtech-test/internal/errors"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

type Validator interface {
	Validate() error
}

type ParametersGetter interface {
	GetParameters(ctx *fasthttp.RequestCtx) error
}

func GetData(ctx *fasthttp.RequestCtx, data any) error {
	if params, ok := data.(ParametersGetter); ok {
		if err := params.GetParameters(ctx); err != nil {
			return err
		}
	}

	if validator, ok := data.(Validator); ok {
		return validator.Validate()
	}

	return nil
}

func SendData(ctx *fasthttp.RequestCtx, data any) {
	ctx.SetContentType("application/json")
	v, err := json.Marshal(data)
	if err != nil {
		// log
		return
	}

	ctx.Response.SetBody(v)
}

func SendError(ctx *fasthttp.RequestCtx, err error) {
	ctx.SetContentType("application/json")
	var apiErr apiErrors.APIError
	ok := errors.As(err, &apiErr)
	if ok {
		resp := Response{
			Status: "error",
			Error:  apiErr.Error(),
			Code:   apiErr.InnerCode,
		}
		v, _ := json.Marshal(&resp)
		ctx.SetStatusCode(apiErr.Code)
		ctx.SetBody(v)
		return
	}

	resp := Response{
		Status: "error",
		Error:  err.Error(),
	}

	v, _ := json.Marshal(&resp)
	ctx.SetStatusCode(http.StatusInternalServerError)
	ctx.SetBody(v)
}
