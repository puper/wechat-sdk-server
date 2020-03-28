package errors

import (
	"github.com/puper/wechat-sdk-server/errors/errcodes"
	"github.com/puper/wechat-sdk-server/errors/errmsgs"

	terrors "github.com/puper/ppgo/errors"
)

func ParamError(err error) *terrors.Error {
	return terrors.Trace(err, 3).SetType(terrors.TypeParamError).
		SetLevel(terrors.LevelDebug).
		SetCode(errcodes.ParamError).
		SetMessage(errmsgs.ParamError)
}

func ServerError(err error) *terrors.Error {
	return terrors.Trace(err, 3).SetType(terrors.TypeServerError).
		SetLevel(terrors.LevelError).
		SetCode(errcodes.ServerError).
		SetMessage(errmsgs.ServerError)
}

type MultiError struct {
	Errors map[string][]string
}

func (this *MultiError) AddError(field, msg string) {
	if _, ok := this.Errors[field]; !ok {
		this.Errors[field] = []string{}
	}
	this.Errors[field] = append(this.Errors[field], msg)
}
