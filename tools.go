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
		replacement := strings.Repeat("*", len(password))
		return strings.ReplaceAll(fullCommand, password, replacement)
	}
	return fullCommand
}

func generateUDID() (string, error) {
	bytes := make([]byte, 16) // 16 字节（128 位）
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
