package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
)

func ExtractFirstFrame(videoPath, finalName string) (string, error) {
	framePath := filepath.Join("public", "covers", finalName+".jpg") // 定义提取的帧图像保存路径

	cmd := exec.Command("ffmpeg", "-y", "-i", videoPath, "-vframes", "1", framePath)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("ffmpeg command failed:", err)
		fmt.Println("ffmpeg command stderr:", stderr.String())
		return "", err
	}

	return "covers/" + finalName + ".jpg", nil // 返回提取的帧图像的 URL
}
