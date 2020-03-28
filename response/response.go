package response

import (
	"fmt"
	"net/http"

	"github.com/puper/orderedmap"

	"github.com/puper/wechat-sdk-server/errors/errcodes"
	"github.com/puper/wechat-sdk-server/errors/errmsgs"

	terrors "github.com/puper/ppgo/errors"

	"github.com/kataras/iris/v12"
)

type Response struct {
	context    iris.Context
	RequestID  string
	statusCode int
	code       int
	message    string
	details    interface{}
	result     interface{}
}

func New(ctx iris.Context) *Response {
	return &Response{
		context:    ctx,
		statusCode: http.StatusOK,
		message:    errmsgs.SUCCESS,
	}
}

func (this *Response) Send() {
	m := orderedmap.New()
	m.Set("code", fmt.Sprintf("%06d", this.code))
	m.Set("message", this.message)
	m.Set("data", this.result)
	this.context.JSON(m)
}

func (this *Response) StatusCode(statusCode int) *Response {
	this.statusCode = statusCode
	return this
}

func (this *Response) Code(code int) *Response {
	this.code = code
	return this
}

func (this *Response) Message(message ...string) *Response {
	if len(message) > 0 {
		this.message = fmt.Sprintf(message[0], func() []interface{} {
			result := []interface{}{}
			for _, row := range message[1:] {
				result = append(result, row)
			}
			return result
		}()...)
	}
	return this
}

func (this *Response) Details(details interface{}) *Response {
	this.details = details
	return this
}

func (this *Response) Result(result interface{}) *Response {
	this.result = result
	return this
}

func (this *Response) ParamError(message ...string) *Response {
	this.Code(errcodes.ParamError).StatusCode(http.StatusBadRequest)
	if len(message) > 0 {
		this.Message(message...)
	} else {
		this.Message(errmsgs.ParamError)
	}
	return this
}

func (this *Response) ServerError(message ...string) *Response {
	this.Code(errcodes.ServerError).StatusCode(http.StatusInternalServerError)
	if len(message) > 0 {
		this.Message(message...)
	} else {
		this.Message(errmsgs.ServerError)
	}
	return this
}

func (this *Response) Error(err error) *Response {
	if terr, ok := err.(*terrors.Error); ok {
		return this.Code(terr.Code).
			Message(terr.Message).
			Details(terr.Details).
			StatusCode(parseErrorTypeToStatusCode(terr.Type))
	}
	return this.ServerError()
}

var (
	errTypeStatusCodeMap = map[string]int{
		terrors.TypeUnset:            http.StatusInternalServerError,
		terrors.TypeServerError:      http.StatusInternalServerError,
		terrors.TypeDirtyData:        http.StatusInternalServerError,
		terrors.TypeIllegal:          http.StatusBadRequest,
		terrors.TypeParamError:       http.StatusBadRequest,
		terrors.TypeNotFound:         http.StatusNotFound,
		terrors.TypePermissionDenied: http.StatusForbidden,
		terrors.TypeUnauthorized:     http.StatusUnauthorized,
	}
)

func parseErrorTypeToStatusCode(errType string) int {
	return errTypeStatusCodeMap[errType]
}
