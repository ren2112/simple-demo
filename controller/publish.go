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
	tokenStr := c.PostForm("token")
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
	_, claims, _ := common.ParseToken(tokenStr)
	finalName := fmt.Sprintf("%d_%s", claims.UserId, filename)

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
	video.PlayUrl = "http://" + serverIp + ":8080/static/" + fmt.Sprintf("%d_%s", claims.UserId, filename)
	video.CoverUrl, err = utils.ExtractFirstFrame(video.PlayUrl, finalName, c)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	video.Title = title
	video.AuthorID = claims.UserId

	common.DB.Create(&video)
	fmt.Println("filename:", filename, finalName+" uploaded successfully")
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
	return
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Query("user_id"))
	videoList := []model.Video{}
	common.DB.Preload("Author").Model(&videoList).Where("author_id=?", int64(userId)).Find(&videoList)
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videoList,
	})
}
