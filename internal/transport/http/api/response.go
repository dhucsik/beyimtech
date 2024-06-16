package api

type Response struct {
	Status string `json:"status"`
	Code   string `json:"code,omitempty"`
	Error  string `json:"error,omitempty"`
}

func NewEmptySuccessResponse() Response {
	return Response{
		Status: "success",
	}
}
