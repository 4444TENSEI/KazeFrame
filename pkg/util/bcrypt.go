// 密码散列加密工具, 负责密码加密存入数据库和登录校验
package util

import "golang.org/x/crypto/bcrypt"

// 密码散列加密
func BcryptPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// 校验密码, 比较哈希值
func ComparePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
