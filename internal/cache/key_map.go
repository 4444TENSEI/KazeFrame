// Resdis键名映射表, 方便统一管理维护
package cache

// 用户在线状态相关: User...
// 邮箱验证码相关: Email...
// IP频繁限制中间件缓存相关: RequestCount
const (
	UserOnlineKey           = "user:online:"
	UserLastSeenKey         = "user:last_seen:"
	EmailCaptchaKey         = "email:captcha:"
	EmailRegisterCaptchaKey = "email:captcha:register:"
	EmailForgetCaptchaKey   = "email:captcha:forget:"
	RequestCount            = "request_count:"
)
