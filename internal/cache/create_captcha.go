// 验证码生成并储存到Redis
package cache

import (
	"KazeFrame/internal/config"
	"context"
	"fmt"
	"time"

	"golang.org/x/exp/rand"
)

// 生成验证码并存储到Redis
func CreateEmailCaptcha(email string, captchaType string) (string, error) {
	rand.Seed(uint64(time.Now().UnixNano()))
	Captcha := fmt.Sprintf("%06d", rand.Intn(900000)+100000)
	ctx := context.Background()
	key := fmt.Sprintf("%s%s:%s", EmailCaptchaKey, captchaType, email)
	if err := config.GetRedis().Set(ctx, key, Captcha, 10*time.Minute).Err(); err != nil {
		return "", fmt.Errorf("存储验证码到Redis失败: %w", err)
	}
	return Captcha, nil
}
