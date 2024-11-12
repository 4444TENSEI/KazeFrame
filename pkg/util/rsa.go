// RSA公钥、私钥生成工具, 并且还有转base64工具
// 工具使用流程:
// 1. 调用GetPublicKey()函数获取公钥, 将公钥转为base64发给前端用于客户端进行数据加密
// 2. 前端接收到bser64编码的RSA公钥进行解码, 用来加密要传输的数据, 最后转base64发给后端
// 3. 后端调用DecryptCiphertext()对于接收到的密文进行base64解密, 使用私钥解密数据, 得到客户端传过来的数据原文
package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"io/ioutil"
	"os"
)

// 定义密钥自动生成和读取路径
const (
	// RSA私钥存放路径
	privateKeyPath = "./key/private.pem"
	// RSA公钥存放路径
	publicKeyPath = "./key/public.pem"
)

// 读取公钥文件，传输过程需要转为Base64字符串不然前端不好解析
func GetPublicKey() (string, error) {
	pubKey, err := readRsaKey(publicKeyPath, "public")
	if err != nil {
		return "", err
	}
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return "", err
	}
	pubKeyPem := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	}
	pubKeyPemBytes := pem.EncodeToMemory(pubKeyPem)
	return encodeBase64Ciphertext(pubKeyPemBytes)
}

// 接收Base64编码的密文，进行解码后使用私钥解密
func DecryptCiphertext(encryptedData string) (string, error) {
	privateKey, err := readRsaKey(privateKeyPath, "private")
	if err != nil {
		return "", err
	}
	decodedData, err := decodeBase64Ciphertext(encryptedData)
	if err != nil {
		return "", err
	}
	decryptedData, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey.(*rsa.PrivateKey), decodedData, nil)
	if err != nil {
		return "", err
	}
	return string(decryptedData), nil
}

// 传入RSA公钥, 进行Base64加密用于发给前端
func encodeBase64Ciphertext(data []byte) (string, error) {
	return base64.StdEncoding.EncodeToString(data), nil
}

// 接收Base64密文, 解码为字节数组
func decodeBase64Ciphertext(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}

// 根据密钥类型读取或生成密钥对
func readRsaKey(filePath string, keyType string) (interface{}, error) {
	// 读取密钥文件
	pemData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, err
	}
	// 根据密钥类型解析密钥
	if keyType == "public" {
		pub, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		return pub, nil
	} else if keyType == "private" {
		privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err == nil {
			return privateKey, nil
		}
		pkcs8Key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		rsaKey, ok := pkcs8Key.(*rsa.PrivateKey)
		if !ok {
			return nil, err
		}
		return rsaKey, nil
	}
	return nil, err
}

// 生成RSA密钥对工具, 可以在程序初始化时调用, 检查私钥文件是否存在, 如果密钥对不存在则自动生成
func createRsaKey() (*rsa.PublicKey, error) {
	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		pubKey, err := saveRsaKey()
		if err != nil {
			return nil, err
		}
		return pubKey, nil
	}
	pubKey, err := readRsaKey(publicKeyPath, "public")
	if err != nil {
		return nil, err
	}
	return pubKey.(*rsa.PublicKey), nil
}

// 生成密钥对并保存到文件
func saveRsaKey() (*rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	// 保存私钥
	privateKeyPem := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	privateKeyPemBytes := pem.EncodeToMemory(privateKeyPem)
	err = ioutil.WriteFile(privateKeyPath, privateKeyPemBytes, 0600)
	if err != nil {
		return nil, err
	}
	// 保存公钥
	pubKey := &privateKey.PublicKey
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return nil, err
	}
	pubKeyPem := &pem.Block{Type: "PUBLIC KEY", Bytes: pubKeyBytes}
	pubKeyPemBytes := pem.EncodeToMemory(pubKeyPem)
	err = ioutil.WriteFile(publicKeyPath, pubKeyPemBytes, 0644)
	if err != nil {
		return nil, err
	}
	return pubKey, nil
}
