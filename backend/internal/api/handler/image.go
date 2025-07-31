package handler

import (
	"context"
	"fmt"
	"path/filepath"

	// "fmt"
	"net/http"
	// "os"
	// "path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/guobinqiu/llm-qwen-emo/internal/model"
	"github.com/guobinqiu/llm-qwen-emo/pkg/oss"

	"github.com/gin-gonic/gin"
	// "github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 上传本地 (非oss)
// func UploadImage(mongoClient *mongo.Client) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		userID := c.GetString("user_id")

// 		// 获取上传的文件
// 		file, err := c.FormFile("file")
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "未选择文件"})
// 			return
// 		}

// 		// 构建本地保存路径
// 		dir := "./images"
// 		os.MkdirAll(dir, 0755) // 确保目录存在

// 		filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(file.Filename))
// 		localPath := filepath.Join(dir, filename)

// 		// 保存图片到本地
// 		if err := c.SaveUploadedFile(file, localPath); err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "保存图片失败"})
// 			return
// 		}

// 		// 写入 MongoDB
// 		coll := mongoClient.Database("video_tasks").Collection("images")
// 		image := model.Image{
// 			ID:        uuid.NewString(),
// 			UserID:    userID,
// 			URL:       os.Getenv("APP_BASE_URL") + "/images/" + filename,
// 			Filename:  filename,
// 			CreatedAt: time.Now(),
// 		}
// 		_, err = coll.InsertOne(context.TODO(), image)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库写入失败"})
// 			return
// 		}

// 		c.JSON(http.StatusOK, gin.H{"image": image})
// 	}
// }

func UploadImage(mongoClient *mongo.Client, ossClient *oss.OSSClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")

		fileHeader, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "未选择文件"})
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "文件读取失败"})
			return
		}
		defer file.Close()

		filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(fileHeader.Filename))
		objectKey := "images/" + filename

		// 上传到 OSS
		url, err := ossClient.UploadReader(objectKey, file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "OSS 上传失败: " + err.Error()})
			return
		}

		image := model.Image{
			ID:        uuid.NewString(),
			UserID:    userID,
			URL:       url,
			Filename:  filename,
			CreatedAt: time.Now(),
		}

		coll := mongoClient.Database("video_tasks").Collection("images")
		_, err = coll.InsertOne(context.TODO(), image)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库写入失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"image": image})
	}
}

func ListImages(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		coll := client.Database("video_tasks").Collection("images")
		cur, err := coll.Find(context.TODO(),
			bson.M{"user_id": userID, "deleted_at": bson.M{"$exists": false}},
			options.Find().SetSort(bson.M{"created_at": -1}),
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
			return
		}
		defer cur.Close(context.TODO())
		var imgs []model.Image
		for cur.Next(context.TODO()) {
			var img model.Image
			cur.Decode(&img)
			imgs = append(imgs, img)
		}
		c.JSON(http.StatusOK, gin.H{"images": imgs})
	}
}

func DeleteImage(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		id := c.Param("id")
		coll := client.Database("video_tasks").Collection("images")
		_, err := coll.UpdateOne(context.TODO(),
			bson.M{"_id": id, "user_id": userID},
			bson.M{"$set": bson.M{"deleted_at": time.Now()}})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "删除失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "删除成功"})
	}
}
