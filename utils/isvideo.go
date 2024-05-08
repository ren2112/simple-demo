package utils

import (
	"path/filepath"
	"strings"
)

func IsVideoFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	// 假设视频文件扩展名为 ".mp4", ".avi", ".mov" 等
	videoExts := []string{".mp4", ".avi", ".mov"} // 添加其他视频格式的扩展名
	for _, videoExt := range videoExts {
		if ext == videoExt {
			return true
		}
	}
	return false
}
