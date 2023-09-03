package middlewares

type ResCode int32

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy

	CodeNeedLogin
	CodeInvalidToken
	CodeInvalidRequest
	CodeIsfollowFailed
	CodeGetFolloweeCountFailed
	CodeGetFollowerCountFailed
	CodeGetFolloweeListFailed
	CodeGetFollowerListFailed
	CodeGetFriendListFailed
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:                "success",
	CodeInvalidParam:           "请求参数错误",
	CodeUserExist:              "用户名已存在",
	CodeUserNotExist:           "用户名不存在",
	CodeInvalidPassword:        "用户名或密码错误",
	CodeServerBusy:             "服务繁忙",
	CodeNeedLogin:              "需要登录",
	CodeInvalidToken:           "无效的token",
	CodeInvalidRequest:         "无效请求",
	CodeIsfollowFailed:         "判断是否关注目标时出错",
	CodeGetFolloweeCountFailed: "获取关注数量时出错",
	CodeGetFollowerCountFailed: "获取粉丝数量时出错",
	CodeGetFolloweeListFailed:  "获取关注列表时出错",
	CodeGetFollowerListFailed:  "获取粉丝列表时出错",
	CodeGetFriendListFailed:    "获取好友列表时出错",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
