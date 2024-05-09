package controller

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"strconv"
)

type VideoListResponse struct {
	Response
	VideoList []model.Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "用户不存在！",
		})
		return
	}
	author := user.(model.User)
	title := c.PostForm("title")
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	var video model.Video

	filename := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d_%s", author.Id, filename)

	//如果不是视频，返回异常
	if utils.IsVideoFile(finalName) == false {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "请上传视频！",
		})
		return
	}

	saveFile := filepath.Join("./public/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	serverIp, err := utils.GetLocalIPv4()
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	video.PlayUrl = "http://" + serverIp + ":8080/static/" + fmt.Sprintf("%d_%s", author.Id, filename)
	video.CoverUrl, err = utils.ExtractFirstFrame(video.PlayUrl, finalName, c)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	video.Title = title
	video.AuthorID = author.Id

	// 开始事务
	tx := common.DB.Begin()

	// 创建视频
	if err := tx.Create(&video).Error; err != nil {
		// 如果创建视频时出现错误，回滚事务
		tx.Rollback()
		// 返回错误
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	// 更新作者的work_count字段
	author.WorkCount++

	// 使用UpdateColumn更新工作计数字段
	if err := tx.Model(&author).UpdateColumn("work_count", author.WorkCount).Error; err != nil {
		// 如果更新作者信息时出现错误，回滚事务
		tx.Rollback()
		// 返回错误
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	// 提交事务
	tx.Commit()

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
	return
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	userId, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "操作失败！",
		})
	}
	videoList := []model.Video{}
	common.DB.Preload("Author").Model(&videoList).Where("author_id=?", int64(userId)).Find(&videoList)
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videoList,
	})
}
