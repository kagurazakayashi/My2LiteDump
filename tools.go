package main

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
)

// ReplaceIfContains 用 '*' 替换字符串 A 中的所有 B
// 示例: fmt.Println(ReplaceIfContains("hello world", "world"))  // 输出: "hello *****"
func ReplaceIfContains(fullCommand string, password string) string {
	if len(fullCommand) > 0 && strings.Contains(fullCommand, password) {
		replacement := "*****" //strings.Repeat("*", len(password))
		return strings.ReplaceAll(fullCommand, password, replacement)
	}
	return fullCommand
}

// MaskPassword 手動解析並替換 MySQL 連線字串中的密碼部分
func MaskPassword(dsn string) string {
	// 正则匹配 "root:密码@tcp(127.0.0.1:3306)/"
	// re := regexp.MustCompile(`([^:]+):([^@]+)@tcp\(([^)]+)\)/`)
	// 替换密码部分为 "*****"
	// return re.ReplaceAllString(dsn, `${1}:*****@tcp(${3})/`)

	// 找到使用者名稱和密碼的分隔符 ":"
	userPassEnd := strings.Index(dsn, "@tcp(")
	if userPassEnd == -1 {
		return dsn // 如果格式不對，直接返回原字串
	}

	// 找到使用者名稱部分
	userPass := dsn[:userPassEnd] // 例如 "root:12345"
	colonIndex := strings.Index(userPass, ":")
	if colonIndex == -1 {
		return dsn // 沒有找到 ":"，返回原字串
	}

	// 構造新的 DSN
	maskedUserPass := userPass[:colonIndex+1] + "*****" // 替換密碼部分
	return maskedUserPass + dsn[userPassEnd:]           // 重新拼接
}

func generateUDID() (string, error) {
	bytes := make([]byte, 16) // 16 字节（128 位）
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
