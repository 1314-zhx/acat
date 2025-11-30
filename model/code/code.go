package code

const (
	// 其它业务逻辑错误码
	Success = 200
	Error   = 400
	// 用户业务错误码 100x
	// login
	InvalidParam         = 1001 // 无效参数
	PhoneOrPasswordError = 1002 // 电话或密码错误
	InvalidPhoneForm     = 1003 // 手机号格式错误
	VCodeError           = 1004 // 验证码错误

	// register
	NameExists                = 2002 // 名字已存在
	PhoneExists               = 2003 // 手机号已存在
	EmailFormError            = 2004 // 邮箱格式不正确
	PasswordUnequalRePassword = 2005

	MissMustInfo = 3001 // 缺少必需字段

	// post
	FirstViewNotPass = 4001
)
