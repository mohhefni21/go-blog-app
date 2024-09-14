package responsepkg

import (
	"mohhefni/go-blog-app/infra/errorpkg"

	"github.com/labstack/echo/v4"
)

type Response struct {
	HttpCode int         `json:"-"`
	Status   string      `json:"status"`
	Message  string      `json:"message,omitempty"`
	Data     interface{} `json:"data,omitempty"`
	Query    interface{} `json:"query,omitempty"`
	Error    string      `json:"error,omitempty"`
	// CKEEditor format
	Url string `json:"url,omitempty"`
}

func NewResponse(params ...func(*Response) *Response) Response {
	var resp = Response{
		Status: "success",
	}

	for _, param := range params {
		param(&resp)
	}

	return resp
}

func WithStatus(err error) func(*Response) *Response {
	return func(r *Response) *Response {
		var receivedErr error

		receivedErr, ok := errorpkg.ErrorMapping[err.Error()]

		if !ok {
			receivedErr = errorpkg.ErrorGeneral
		}

		myError, ok := receivedErr.(errorpkg.Error)
		if !ok {
			myError = errorpkg.ErrorGeneral
		}

		r.Status = "fail"
		r.HttpCode = myError.HttpCode
		r.Message = myError.Message

		if myError == errorpkg.ErrorGeneral {
			r.Status = "error"
			r.Error = err.Error()
		}

		return r
	}

}

func WithHttpCode(httpCode int) func(*Response) *Response {
	return func(r *Response) *Response {
		r.HttpCode = httpCode
		return r
	}
}

func WithMessage(message string) func(*Response) *Response {
	return func(r *Response) *Response {
		r.Message = message
		return r
	}
}

func WithData(data interface{}) func(*Response) *Response {
	return func(r *Response) *Response {
		r.Data = data
		return r
	}
}

func WithQuery(query interface{}) func(*Response) *Response {
	return func(r *Response) *Response {
		r.Query = query
		return r
	}
}

func (r Response) Send(ctx echo.Context) error {
	return ctx.JSON(r.HttpCode, r)
}

func WithUrlCKEditor(url string) func(*Response) *Response {
	return func(r *Response) *Response {
		r.Url = url
		return r
	}
}
