package base

import (
	"io"
	"os"
	"strings"
)

// StringFirstUpper 字符串首字母大写
func StringFirstUpper(s string) string {
	if s == "" {
		return ""
	}

	return strings.ToUpper(s[:1]) + s[1:]
}

// StringFirstLower 字符串首字母小写
func StringFirstLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// Copy 复制文件
// @param dstPath 复制内容到目标
// @param srcPath 复制内容来源
// @return error
func Copy(dstPath string, srcPath string) error {
	out, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer out.Close()

	in, err := os.OpenFile(dstPath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer in.Close()

	_, err = io.Copy(in, out)
	return err
}
