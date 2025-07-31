package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/guobinqiu/llm-qwen-emo/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 上传本地 (非oss)
func UploadAudio(mongoClient *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")

		// 获取上传的文件
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "未选择文件"})
			return
		}

		// 构建本地保存路径
		dir := "./audios"
		os.MkdirAll(dir, 0755) // 确保目录存在

		filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(file.Filename))
		localPath := filepath.Join(dir, filename)

		// 保存文件到本地
		if err := c.SaveUploadedFile(file, localPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "保存本地文件失败"})
			return
		}

		// 可选：提取音频时长
		duration := 0

		// 写入 MongoDB
		coll := mongoClient.Database("video_tasks").Collection("audios")
		audio := model.Audio{
			ID:        uuid.NewString(),
			UserID:    userID,
			URL:       os.Getenv("APP_BASE_URL") + "/audios/" + filename,
			Filename:  filename,
			Duration:  duration,
			CreatedAt: time.Now(),
			Size:      file.Size,
		}
		_, err = coll.InsertOne(context.TODO(), audio)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库写入失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"audio": audio})
	}
}

func ListAudios(mongoClient *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		coll := mongoClient.Database("video_tasks").Collection("audios")
		cur, err := coll.Find(context.TODO(),
			bson.M{"user_id": userID, "deleted_at": bson.M{"$exists": false}},
			options.Find().SetSort(bson.M{"created_at": -1}),
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
			return
		}
		defer cur.Close(context.TODO())
		var audios []model.Audio
		for cur.Next(context.TODO()) {
			var audio model.Audio
			cur.Decode(&audio)
			audios = append(audios, audio)
		}
		c.JSON(http.StatusOK, gin.H{"audios": audios})
	}
}

func DeleteAudio(mongoClient *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		id := c.Param("id")
		coll := mongoClient.Database("video_tasks").Collection("audios")
		res, err := coll.DeleteOne(context.TODO(), bson.M{"_id": id, "user_id": userID})
		if err != nil || res.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "删除失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "删除成功"})
	}
}
