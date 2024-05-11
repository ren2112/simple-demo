package controller

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/assist"
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/response"
	"github.com/RaymondCode/simple-demo/utils"
	"github.com/gin-gonic/gin"
	"path/filepath"
	"strconv"
)

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		response.CommonResp(c, 1, "用户不存在")
		return
	}
	author := user.(model.User)
	title := c.PostForm("title")
	data, err := c.FormFile("data")
	if err != nil {
		response.CommonResp(c, 1, err.Error())
		return
	}
	var video model.Video

	filename := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d_%s", author.Id, filename)

	//如果不是视频，返回异常
	if utils.IsVideoFile(finalName) == false {
		response.CommonResp(c, 1, "请上传视频！")
		return
	}

	saveFile := filepath.Join("./public/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		response.CommonResp(c, 1, err.Error())
		return
	}
	serverIp, err := utils.GetLocalIPv4()
	if err != nil {
		response.CommonResp(c, 1, err.Error())
		return
	}
	video.PlayUrl = "http://" + serverIp + ":8080/static/" + fmt.Sprintf("%d_%s", author.Id, filename)

	//获得封面并且保存封面图片于服务器
	video.CoverUrl, err = utils.ExtractFirstFrame(video.PlayUrl, finalName, c)
	if err != nil {
		response.CommonResp(c, 1, err.Error())
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
		response.CommonResp(c, 1, err.Error())
		return
	}

	// 更新作者的work_count字段
	author.WorkCount++

	// 使用UpdateColumn更新作品计数字段
	if err := tx.Model(&author).UpdateColumn("work_count", author.WorkCount).Error; err != nil {
		// 如果更新作者信息时出现错误，回滚事务
		tx.Rollback()
		// 返回错误
		response.CommonResp(c, 1, err.Error())
		return
	}

	// 提交事务
	tx.Commit()
	response.CommonResp(c, 0, finalName+" uploaded successfully")
	return
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	userId, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		response.CommonResp(c, 1, "操作失败")
	}
	videoList := []model.Video{}
	RespVideoList := []model.RespVideo{}
	common.DB.Preload("Author").Model(&videoList).Where("author_id=?", int64(userId)).Find(&videoList)

	//将videoList转化为响应结构体
	for _, v := range videoList {
		RespVideoList = append(RespVideoList, assist.ToRespVideo(v))
	}
	response.VideoListResponseFun(c, response.Response{StatusCode: 0}, RespVideoList)
}
