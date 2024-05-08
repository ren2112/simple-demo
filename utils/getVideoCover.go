package utils

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"os/exec"
	"path/filepath"
)

func ExtractFirstFrame(videoPath, finalName string, c *gin.Context) (string, error) {
	framePath := filepath.Join("public", "cover", "frame_"+finalName+".jpg") // 定义提取的帧图像保存路径

	cmd := exec.Command("ffmpeg", "-y", "-i", videoPath, "-vframes", "1", framePath)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("ffmpeg command failed:", err)
		fmt.Println("ffmpeg command stderr:", stderr.String())
		return "", err
	}

	ipv4, err := GetLocalIPv4()
	if err != nil {
		return "", err
	}

	return "http://" + ipv4 + ":8080/static/cover/frame_" + finalName + ".jpg", nil // 返回提取的帧图像的 URL
}
