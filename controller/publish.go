package controller

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/config"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/RaymondCode/simple-demo/response"
	"github.com/RaymondCode/simple-demo/service"
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

	// 增加视频压缩上传逻辑
	playUrl, err := service.CompressAndUploadVideo(c, data, &author)
	if err != nil {
		response.CommonResp(c, 1, err.Error())
		return
	}
	video.PlayUrl = playUrl

	//获得封面并且保存封面图片于服务器
	video.CoverUrl, err = utils.ExtractFirstFrame(config.SERVER_RESOURCES+video.PlayUrl, finalName)
	if err != nil {
		response.CommonResp(c, 1, err.Error())
		return
	}

	video.Title = title
	video.AuthorID = author.Id

	err = service.PublishVideo(video, author)
	if err != nil {
		response.CommonResp(c, 1, err.Error())
	} else {
		response.CommonResp(c, 0, finalName+" uploaded successfully")
	}
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	userId, err := strconv.Atoi(c.Query("user_id"))
	if err != nil {
		response.CommonResp(c, 1, "操作失败")
	}

	//获取发布列表
	var RespVideoList []model.RespVideo
	RespVideoList, err = service.GetPublishVideoList(int64(userId))
	if err != nil {
		response.CommonResp(c, 1, err.Error())
	} else {
		response.VideoListResponseFun(c, response.Response{StatusCode: 0}, RespVideoList)
	}

}
