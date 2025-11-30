package code

var mapMsg = map[int]string{
	// 通用
	Success: "成功",
	Error:   "失败",
	// 用户
	InvalidParam:         "无效参数-用户登录",
	PhoneOrPasswordError: "手机号或密码错误-用户登录",
	InvalidPhoneForm:     "手机号格式错误-用户登录",
	VCodeError:           "验证码错误-用户登录",

	NameExists:                "名字已存在-用户注册",
	PhoneExists:               "手机号已存在-用户注册",
	EmailFormError:            "邮箱格式不正确-用户注册",
	PasswordUnequalRePassword: "密码与第二次不一样",

	MissMustInfo: "缺少必需字段",

	FirstViewNotPass: "一面未通过，无法参加二面",
	// 管理员
}

// GetMsg 提供给包外的获取业务码内容的函数
func GetMsg(code int) string {
	if msg, ok := mapMsg[code]; ok {
		return msg
	}
	return "未知错误"
}
