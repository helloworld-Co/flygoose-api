package comm

import (
	"errors"
	"github.com/kataras/iris/v12"
	"strconv"
	"strings"
)

var (
	CodeSuccess         = 1
	CodeFailed          = 0
	CodeParamError      = 1001 //request param error
	CodePhoneExist      = 1002 //account is exist
	CodeServerError     = 1003 //server internal error
	CodeTokenExpired    = 1004 // token expired
	CodeUserNotExist    = 1005 //user not exist
	CodePasswordError   = 1006 //password error
	CodeProfileNotExist = 1007 //profile not exist
	CodeNotExist        = 2000 //不存在
	CodeUserBlock       = 2001 //用户被封
	CodeFrequentlyReq   = 2002 //请求过于频繁
)
var (
	IrisFormNotExist = "www.helloworld.net_is_developer_owner_technical_community"
)

type GooseEgg struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type BaseComponent struct {
	Ctx iris.Context
}

// HandleError 解决post请求，form表单出出错的 "schema: invalid path" 错误问题
func (c *BaseComponent) HandleError(ctx iris.Context, err error) {
	// to ignore any "schema: invalid path" you can check the error type
	// as we do here and don't stop the execution.
	if iris.IsErrPath(err) {
		return
	}

	// on any other error, stop the execution of the handler chain.
	ctx.StopExecution()
}

func (c *BaseComponent) PostFormIntDefault(key string) (int, bool) {
	val := c.Ctx.FormValueDefault(key, IrisFormNotExist)
	if val == IrisFormNotExist {
		return -1, false
	} else {
		val = strings.TrimSpace(val)
		v, err := strconv.Atoi(val)
		if err != nil {
			return -1, true
		} else {
			return v, true
		}
	}
}

func (c *BaseComponent) PostFormInt(key string) (int, error) {
	s := strings.TrimSpace(c.Ctx.FormValue(key))
	if s == "" {
		return 0, errors.New(key + "参数错误")
	}

	if value, err := strconv.Atoi(s); err != nil {
		return 0, err
	} else {
		return value, nil
	}
}

func (c *BaseComponent) PostFormInt64Default(key string) (int64, bool) {
	val := c.Ctx.FormValueDefault(key, IrisFormNotExist)
	if val == IrisFormNotExist {
		return -1, false
	} else {
		val = strings.TrimSpace(val)
		v, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return -1, true
		} else {
			return v, true
		}
	}
}

func (c *BaseComponent) PostFormInt64(key string) (int64, error) {
	s := strings.TrimSpace(c.Ctx.FormValue(key))
	if s == "" {
		return 0, errors.New(key + "参数错误")
	}

	if value, err := strconv.ParseInt(s, 10, 64); err != nil {
		return 0, err
	} else {
		return value, nil
	}
}

// PostFormStringDefault 返回 IrisFormNotExist 时，说明 key 不存在
func (c *BaseComponent) PostFormStringDefault(key string) (string, bool) {
	val := c.Ctx.FormValueDefault(key, IrisFormNotExist)
	if val == IrisFormNotExist {
		return "", false
	} else {
		return strings.TrimSpace(val), true
	}
}

func (c *BaseComponent) PostFormString(key string) string {
	return strings.TrimSpace(c.Ctx.FormValue(key))
}

func (c *BaseComponent) RespJson(code int, data interface{}, message string) {
	c.Ctx.JSON(GooseEgg{
		Code:    code,
		Data:    data,
		Message: message,
	})
}

func (c *BaseComponent) RespSuccess(data interface{}, message string) {
	c.Ctx.JSON(GooseEgg{
		Code:    CodeSuccess,
		Data:    data,
		Message: message,
	})
}

func (c *BaseComponent) RespFailed(data interface{}, message string) {
	c.Ctx.JSON(GooseEgg{
		Code:    CodeFailed,
		Data:    data,
		Message: message,
	})
}

func (c *BaseComponent) RespFailedMessage(message string) {
	c.Ctx.JSON(GooseEgg{
		Code:    CodeFailed,
		Data:    nil,
		Message: message,
	})
}

func (c *BaseComponent) RespFailedParam() {
	c.Ctx.JSON(GooseEgg{
		Code:    CodeFailed,
		Data:    nil,
		Message: "参数错误",
	})
}

func (c *BaseComponent) RespFailedSession() {
	c.Ctx.JSON(GooseEgg{
		Code:    CodeFailed,
		Data:    nil,
		Message: "会话已过期，请重新登录",
	})
}

func (c *BaseComponent) RespSuccessMessage(message string) {
	c.Ctx.JSON(GooseEgg{
		Code:    CodeSuccess,
		Data:    nil,
		Message: message,
	})
}

func (c *BaseComponent) NewDict() map[string]interface{} {
	return make(map[string]interface{}, 0)
}
