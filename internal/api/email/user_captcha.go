package email

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"KazeFrame/internal/cache"
	"KazeFrame/internal/model"
	"KazeFrame/internal/service"
	"KazeFrame/pkg/util"

	"github.com/gin-gonic/gin"
)

// 用户注册-邮件验证码发送接口
func SendRegisterCaptcha(c *gin.Context) {
	// 绑定请求所需的参数
	var sendCaptchaPayload model.EmailCaptchaPayload
	if err := c.ShouldBindJSON(&sendCaptchaPayload); err != nil {
		util.Rsp(c, 400, "请求参数错误, "+err.Error())
		return
	}
	// 检查邮箱是否已存在
	existingUser := service.IsExistingData("email", sendCaptchaPayload.Email)
	if existingUser != nil {
		util.Rsp(c, 409, "邮箱已被注册, 请直接登录或找回密码")
		return
	}
	// 这里的第三个参数字符串将会拼接到Redis键名
	sendCaptchaEmail(c, sendCaptchaPayload.Email, "register", "期待您的加入!")
}

// 用户找回密码-邮件验证码发送接口
func SendForgetPswCaptcha(c *gin.Context) {
	// 绑定请求所需的参数
	var sendCaptchaPayload model.EmailCaptchaPayload
	if err := c.ShouldBindJSON(&sendCaptchaPayload); err != nil {
		util.Rsp(c, 400, "请求参数错误, "+err.Error())
		return
	}
	// 检查邮箱是否已存在
	existingUser := service.IsExistingData("email", sendCaptchaPayload.Email)
	if existingUser == nil {
		util.Rsp(c, 409, 4602)
		return
	}
	// 这里的第三个参数字符串将会拼接到Redis键名
	sendCaptchaEmail(c, sendCaptchaPayload.Email, "forget", "查看您的重置验证码!")
}

// 根据用户邮箱和验证码类型, 动态构建验证码邮件标题和内容
func sendCaptchaEmail(c *gin.Context, email string, captchaType string, emailTitle string) {
	EmailCaptcha, err := cache.CreateEmailCaptcha(email, captchaType)
	if err != nil {
		util.Rsp(c, 500, err.Error())
		return
	}
	err = CaptchaHTML(email, EmailCaptcha, captchaType, emailTitle)
	if err != nil {
		util.Rsp(c, 500, err.Error())
		return
	}
	util.Rsp(c, 200, 4500)
}

// 调用service.SendMail()通用发送邮件服务, 读取HTML模板文件, 发送验证码HTML邮件
func CaptchaHTML(userEmail string, Captcha string, captchaType string, emailTitle string) error {
	// 验证码模板资源路径
	templatePath := "./static/server/template/email/captcha.html"
	captchaTypeMap := map[string]string{
		"register": "注册",
		"forget":   "密码重置",
	}
	templateData, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("读取邮件模板失败: %w", err)
	}
	tmpl, err := template.New("register_template").Parse(string(templateData))
	if err != nil {
		return fmt.Errorf("解析HTML模板失败: %w", err)
	}
	var tpl bytes.Buffer
	CaptchaType := captchaTypeMap[captchaType]
	if err := tmpl.Execute(&tpl, map[string]string{
		"Captcha":     Captcha,
		"CaptchaType": CaptchaType,
		"EmailTitle":  emailTitle,
	}); err != nil {
		return fmt.Errorf("执行HTML模板失败: %w", err)
	}
	emailContent := tpl.String()
	if err := service.SendMail(
		userEmail,
		emailTitle,
		emailContent,
	); err != nil {
		return fmt.Errorf("发送邮件失败, 原因: %w", err)
	}
	return nil
}
