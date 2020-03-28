package errors

import (
	"fmt"

	"github.com/puper/tracerr"
	"github.com/sirupsen/logrus"
)

var defaultLogger = logrus.New()

var DefaultLoggerFunc = func() *logrus.Logger {
	return defaultLogger
}

// error type
const (
	TypeUnset            = "unset"
	TypeServerError      = "serverError"      //系统错误
	TypeDirtyData        = "dirtyData"        //出现了脏数据
	TypeIllegal          = "illegal"          //非法请求，正常流程不应该存在的问题
	TypeParamError       = "paramError"       //请求错误
	TypeNotFound         = "notFound"         //找不懂资源
	TypePermissionDenied = "permissionDenied" //没有权限
	TypeUnauthorized     = "unauthorized"     //未登录
)

// error level
const (
	LevelPanic = "panic"
	LevelFatal = "fatal"
	LevelError = "error"
	LevelWarn  = "warn"
	LevelInfo  = "info"
	LevelDebug = "debug"
)

const (
	CodeUnset = -1
)

type Error struct {
	Type    string
	Level   string
	Code    int
	Message string
	Params  []interface{}
	Details []interface{}
	CauseBy tracerr.Error
}

func Trace(err error, skip ...int) *Error {
	if terr, ok := err.(*Error); ok {
		return terr
	}
	if len(skip) > 0 {
		return New().SetCauseBy(err, skip[0])
	}
	return New().SetCauseBy(err, 2)
}

func New() *Error {
	return &Error{}
}

func (this *Error) Reset() *Error {
	this.Type = TypeUnset
	this.Level = LevelDebug
	this.Code = CodeUnset
	this.Message = ""
	this.Params = nil
	this.Details = nil
	this.CauseBy = nil
	return this
}

func (this *Error) SetType(type_ string) *Error {
	this.Type = type_
	return this
}

func (this *Error) SetLevel(level string) *Error {
	this.Level = level
	return this
}

func (this *Error) SetCode(code int) *Error {
	this.Code = code
	return this
}

func (this *Error) SetMessage(message ...string) *Error {
	if len(message) > 0 {
		this.Message = fmt.Sprintf(message[0], func() (result []interface{}) {
			for _, v := range message[1:] {
				result = append(result, v)
			}
			return result
		}()...)
	}
	return this
}

func (this *Error) SetParams(params ...interface{}) *Error {
	this.Params = params
	return this
}

func (this *Error) SetDetails(details ...interface{}) *Error {
	this.Details = details
	return this
}

func (this *Error) SetCauseBy(err error, skip int) *Error {
	if terr, ok := err.(tracerr.Error); ok {
		this.CauseBy = terr
	} else {
		this.CauseBy = tracerr.WrapSkip(err, skip)
	}
	return this
}

func (this *Error) Log(logger ...*logrus.Logger) *Error {
	if len(logger) > 0 {
		return this.log(logger[0], false)
	}
	return this.log(DefaultLoggerFunc(), false)
}

func (this *Error) LogWithTrace(logger ...*logrus.Logger) *Error {
	if len(logger) > 0 {
		return this.log(logger[0], true)
	}
	return this.log(DefaultLoggerFunc(), true)
}

func (this *Error) log(logger *logrus.Logger, withTrace bool) *Error {
	entry := logger.WithFields(logrus.Fields{
		"causeBy": this.CauseBy,
		"type":    this.Type,
		"params":  this.Params,
		"code":    this.Code,
		"details": this.Details,
	})
	if withTrace && this.CauseBy != nil {
		entry = entry.WithField("causeBy", tracerr.SprintFirst(this.CauseBy, []int{5}, 2, 10))
	}
	if this.Level == LevelDebug {
		entry.Debugln(this.Message)
	} else if this.Level == LevelInfo {
		entry.Infoln(this.Message)
	} else if this.Level == LevelWarn {
		entry.Warnln(this.Message)
	} else if this.Level == LevelError {
		entry.Errorln(this.Message)
	} else if this.Level == LevelFatal {
		entry.Fatalln(this.Message)
	} else if this.Level == LevelPanic {
		entry.Panicln(this.Message)
	}
	return this
}

func (this *Error) Error() string {
	return this.Message
}
