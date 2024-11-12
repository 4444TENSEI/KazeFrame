package model

// 用户登录接口-请求体所需参数
// 在此设置登录数据校验, alphanum代表字母数字组合
type UserLoginPayload struct {
	Username   string `json:"username" binding:"omitempty,alphanum,min=2,max=20" comment:"用户名"`
	Email      string `json:"email" binding:"omitempty,email,max=50" comment:"邮箱"`
	Password   string `json:"password" binding:"required,alphanum,min=4,max=20" comment:"密码"`
	RememberMe bool   `json:"remember_me" binding:"omitempty" comment:"设置refresh_token为30天有效期, 实现30天内免登录"`
}

// 用户注册接口-请求体所需参数
// 在此设置注册数据校验, alphanum代表字母数字组合
type UserRegisterPayload struct {
	Nickname        string `json:"nickname" binding:"omitempty,min=2,max=20" comment:"昵称"`
	Username        string `json:"username" binding:"required,alphanum,min=2,max=20" comment:"用户名"`
	Email           string `json:"email" binding:"required,email,max=50" comment:"邮箱"`
	Password        string `json:"password" binding:"required,alphanum,min=4,max=20" comment:"密码"`
	RegisterCaptcha string `json:"register_captcha" binding:"omitempty,min=6,max=6" comment:"邮箱验证码"`
}

// 用户个人资料更新接口-请求体所需参数
type UserUpdatePayload struct {
	Nickname      string `json:"nickname" binding:"omitempty,min=2,max=20" comment:"昵称"`
	Signature     string `json:"signature" binding:"omitempty,min=0,max=56" comment:"个性签名"`
	Gender        string `json:"gender" binding:"omitempty,min=0,max=4" comment:"性别"`
	AvatarUrl     string `json:"avatar_url" binding:"omitempty,url,min=0,max=56" comment:"头像url"`
	BackgroundUrl string `json:"background_url" binding:"omitempty,url,min=0,max=56" comment:"背景图url"`
}

// 删除接口通用-请求体所需参数
type DeletePayload struct {
	Field string `json:"field" binding:"required"`
	Value []any  `json:"value" binding:"required"`
}

// 用户找回密码接口-请求体所需参数
type UserForgetPayload struct {
	Email            string `json:"email" binding:"required,email"`
	NewPassword      string `json:"new_password" binding:"required"`
	ForgetPswCaptcha string `json:"forget_captcha" binding:"required"`
}

// 用户注册、找回邮件接口-请求体所需参数
type EmailCaptchaPayload struct {
	Email string `json:"email" binding:"required,email"`
}

// 清理Redis缓存接口-请求体所需参数
type ClearCachePayload struct {
	CacheKey []string `json:"cache_key" binding:"required"`
}

// 删除用户接口-请求体所需参数
type DeleteUserResponse struct {
	OkCount  int      `json:"ok_count"`
	ErrCount int      `json:"err_count"`
	OkValue  []string `json:"ok_value"`
	ErrValue []string `json:"fail_value"`
}
