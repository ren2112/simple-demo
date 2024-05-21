package service

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/config"
	"github.com/RaymondCode/simple-demo/model"
	pb "github.com/RaymondCode/simple-demo/rpc-service/proto"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func ToRespVideo(video model.Video) model.RespVideo {
	respVideo := model.RespVideo{
		Id:            video.Id,
		Author:        ToRespUser(video.Author),
		Title:         video.Title,
		PlayUrl:       video.PlayUrl,
		CoverUrl:      video.CoverUrl,
		FavoriteCount: video.FavoriteCount,
		CommentCount:  video.CommentCount,
		IsFavorite:    video.IsFavorite,
	}
	return respVideo
}

func ToProtoVideo(video model.Video) pb.Video {
	protoUser := ToProtoUser(video.Author)
	respVideo := pb.Video{
		Id:            video.Id,
		Author:        &protoUser,
		Title:         video.Title,
		PlayUrl:       video.PlayUrl,
		CoverUrl:      video.CoverUrl,
		FavoriteCount: video.FavoriteCount,
		CommentCount:  video.CommentCount,
		IsFavorite:    video.IsFavorite,
	}
	return respVideo
}

func FeedVideoList(latestTime int64) (videoList []*pb.Video, err error) {
	// 查询数据库
	var modelVidelList []model.Video
	err = common.DB.Model(&model.Video{}).
		Preload("Author").
		Order("created_at DESC").
		Limit(config.VIDEO_STREAM_BATCH_SIZE).
		Where("created_at < ?", latestTime).
		Find(&modelVidelList).Error
	if err != nil {
		return nil, err
	}
	for _, v := range modelVidelList {
		protoVideo := ToProtoVideo(v)
		videoList = append(videoList, &protoVideo)
	}
	return videoList, nil
}

func JudgeFavorite(userId int64, videoId int64) (bool, error) {
	var isFavorite bool
	err := common.DB.Model(&model.Favorite{}).
		Where("user_id = ? AND video_id = ?", userId, videoId).
		Pluck("is_favorite", &isFavorite).Error
	return isFavorite, err
}

func JudgeRelation(authorId, userId int64) error {
	var follow model.Follow
	err := common.DB.Where("user_id = ? AND follower_user_id = ?", authorId, userId).First(&follow).Error
	return err
}

func PublishVideo(video model.Video, author model.User) error {
	var err error

	// 开始事务
	tx := common.DB.Begin()

	// 创建视频
	if err = tx.Create(&video).Error; err != nil {
		// 如果创建视频时出现错误，回滚事务
		tx.Rollback()
		// 返回错误
		return err
	}

	// 更新作者的work_count字段
	author.WorkCount++

	// 使用UpdateColumn更新作品计数字段
	if err = tx.Model(&author).UpdateColumn("work_count", author.WorkCount).Error; err != nil {
		// 如果更新作者信息时出现错误，回滚事务
		tx.Rollback()
		// 返回错误
		return err
	}

	// 提交事务
	tx.Commit()

	return err
}

func GetPublishVideoList(userId int64) ([]model.RespVideo, error) {
	var videoList []model.Video
	RespVideoList := []model.RespVideo{}
	err := common.DB.Preload("Author").Model(&videoList).Where("author_id=?", userId).Find(&videoList).Error
	//转化为响应专用
	for _, v := range videoList {
		RespVideoList = append(RespVideoList, ToRespVideo(v))
	}
	return RespVideoList, err
}

// CompressAndUploadVideo 压缩视频并上传至指定目录
func CompressAndUploadVideo(c *gin.Context, data *multipart.FileHeader, author *model.User) (string, error) {
	// 确保public/tmp_videos/目录存在
	tmpVideosPath := "./public/tmp_videos/"
	err := os.MkdirAll(tmpVideosPath, 0755)
	if err != nil && !os.IsExist(err) {
		return "", err
	}
	defer func() {
		// 清理临时文件，注意实际应用中需谨慎处理，避免删除非预期文件
		_ = os.RemoveAll(tmpVideosPath)
	}()

	// 使用public下的临时目录
	tempDir := tmpVideosPath
	originalFilePath := filepath.Join(tempDir, data.Filename)

	// 保存上传的原始视频到临时目录
	if err = c.SaveUploadedFile(data, originalFilePath); err != nil {
		return "", err
	}

	// 定义压缩后的视频文件名
	compressedFilePath := filepath.Join(tempDir, "compressed_"+data.Filename)

	// 使用FFmpeg命令压缩视频
	ffmpegArgs := []string{
		"-i", originalFilePath,
		"-c:v", "libx264",
		"-preset", "medium",
		"-crf", "23",
		"-c:a", "aac",
		"-b:a", "128k",
		compressedFilePath,
	}
	cmd := exec.Command("ffmpeg", ffmpegArgs...)
	if err = cmd.Run(); err != nil {
		return "", err
	}

	// 上传压缩后的视频到public/videos目录
	finalName := fmt.Sprintf("%d_%s", author.Id, strings.TrimSuffix(data.Filename, filepath.Ext(data.Filename)))
	saveFile := filepath.Join("./public/videos/", finalName+".mp4")
	if err = os.Rename(compressedFilePath, saveFile); err != nil {
		return "", err
	}

	return "videos/" + finalName + ".mp4", nil
}
