package handler

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/guobinqiu/llm-qwen-emo/internal/model"
	"github.com/guobinqiu/llm-qwen-emo/pkg/oss"
	"github.com/guobinqiu/llm-qwen-emo/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateTask(mongoClient *mongo.Client, ossClient *oss.OSSClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		var req struct {
			Name          string `json:"name"`
			ImageID       string `json:"image_id"`
			AudioID       string `json:"audio_id"`
			SegmentSecond int    `json:"segment_second"`
		}
		if err := c.ShouldBindJSON(&req); err != nil || req.ImageID == "" || req.AudioID == "" {
			c.JSON(400, gin.H{"error": "参数错误"})
			return
		}

		imgColl := mongoClient.Database("video_tasks").Collection("images")
		audioColl := mongoClient.Database("video_tasks").Collection("audios")

		var img model.Image
		if err := imgColl.FindOne(context.TODO(), bson.M{"_id": req.ImageID, "user_id": userID}).Decode(&img); err != nil {
			c.JSON(400, gin.H{"error": "图片不存在"})
			return
		}

		var audio model.Audio
		if err := audioColl.FindOne(context.TODO(), bson.M{"_id": req.AudioID, "user_id": userID}).Decode(&audio); err != nil {
			c.JSON(400, gin.H{"error": "音频不存在"})
			return
		}

		task := model.Task{
			ID:        uuid.NewString(),
			UserID:    userID,
			Name:      req.Name,
			ImageID:   img.ID,
			AudioID:   audio.ID,
			ImageURL:  img.URL,
			CreatedAt: time.Now(),
			SubTasks:  []model.SubTask{},
			Status:    model.TaskStatusRunning,
		}

		audioDir := "audios/" + task.ID //本地存储目录
		os.MkdirAll(audioDir, 0755)

		audioFiles, err := utils.SliceAudio("audios/"+audio.Filename, audioDir, req.SegmentSecond)
		if err != nil {
			fmt.Println("音频切割失败:", err)
			c.JSON(500, gin.H{"error": "音频切割失败"})
			return
		}

		for _, audioFile := range audioFiles {
			ossURL, err := ossClient.Upload(audioFile, audioFile)
			if err != nil {
				c.JSON(500, gin.H{"error": "音频片段上传失败"})
				return
			}
			task.SubTasks = append(task.SubTasks, model.SubTask{
				AudioURL:   ossURL,
				TaskID:     "",
				TaskStatus: model.TaskStatusPending,
				SubmitTime: time.Now(),
				VideoURL:   "",
			})
		}
		taskColl := mongoClient.Database("video_tasks").Collection("tasks")
		_, err = taskColl.InsertOne(context.TODO(), task)
		if err != nil {
			c.JSON(500, gin.H{"error": "任务创建失败"})
			return
		}
		go processTask(mongoClient, task.ID)
		c.JSON(200, gin.H{"task": task})
	}
}

func processTask(mongoClient *mongo.Client, taskID string) {
	taskColl := mongoClient.Database("video_tasks").Collection("tasks")
	var task model.Task
	if err := taskColl.FindOne(context.TODO(), bson.M{"_id": taskID}).Decode(&task); err != nil {
		return
	}

	// 1. 图片检测
	faceBBox, extBBox, err := utils.CheckImage(task.ImageURL)
	if err != nil {
		task.Status = model.TaskStatusFailed
		task.Message = fmt.Sprintf("图片检测失败: %v", err)
		taskColl.UpdateByID(context.TODO(), task.ID, bson.M{
			"$set": bson.M{
				"message": err.Error(),
			},
		})
		return
	}

	// 2. 为每个子任务调用生成视频接口
	for i := range task.SubTasks {
		fmt.Println(i)
		videoTaskID, err := utils.GenerateVideo(task.ImageURL, task.SubTasks[i].AudioURL, faceBBox, extBBox)
		task.SubTasks[i].TaskID = videoTaskID
		if err != nil {
			task.SubTasks[i].TaskStatus = model.TaskStatusFailed
			task.SubTasks[i].Message = fmt.Sprintf("视频生成失败: %v", err)
			task.Status = model.TaskStatusFailed
			task.Message = fmt.Sprintf("视频生成失败: %v", err)
			break
		} else {
			task.SubTasks[i].TaskStatus = model.TaskStatusPending
		}
	}
	taskColl.UpdateByID(context.TODO(), task.ID, bson.M{"$set": bson.M{"sub_tasks": task.SubTasks, "status": task.Status, "message": task.Message}})

	// 3. 轮询任务状态直到全部完成
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	timeout := time.After(24 * time.Hour)

loop:
	for {
		select {
		case <-timeout:
			break loop
		case <-ticker.C:
			allDone := true
			for i := range task.SubTasks {
				sub := &task.SubTasks[i]
				if sub.TaskStatus == model.TaskStatusSucceeded || sub.TaskStatus == model.TaskStatusFailed {
					continue
				}
				status, videoURL, code, message, err := utils.QueryTaskStatus(sub.TaskID)
				if err != nil {
					allDone = false
					continue
				}

				sub.TaskStatus = status
				sub.Code = code
				sub.Message = message
				if status == model.TaskStatusSucceeded {
					sub.VideoURL = videoURL
				}

				// 中间状态
				if status != model.TaskStatusSucceeded && status != model.TaskStatusFailed {
					allDone = false
				}
			}

			taskColl.UpdateByID(context.TODO(), task.ID, bson.M{"$set": bson.M{"sub_tasks": task.SubTasks}})

			if allDone {
				break loop
			}
		}
	}

	// 4. 判断是否所有子任务都成功
	allSucceeded := true
	for _, sub := range task.SubTasks {
		if sub.TaskStatus != model.TaskStatusSucceeded {
			allSucceeded = false
			break
		}
	}

	if !allSucceeded {
		taskColl.UpdateByID(context.TODO(), task.ID, bson.M{"$set": bson.M{
			"status":  model.TaskStatusFailed,
			"message": "存在失败的子任务取消视频合成",
		}})
		return
	}

	// 5. 合成视频
	videoDir := "videos/" + task.ID
	os.MkdirAll(videoDir, 0755)

	var videoPaths []string
	for i, sub := range task.SubTasks {
		localPath := filepath.Join(videoDir, fmt.Sprintf("part_%d.mp4", i))
		if err := utils.DownloadFile(sub.VideoURL, localPath); err != nil {
			taskColl.UpdateByID(context.TODO(), task.ID, bson.M{"$set": bson.M{
				"status":  model.TaskStatusFailed,
				"message": fmt.Sprintf("视频片段下载失败: %v", err),
			}})
			return
		}
		videoPaths = append(videoPaths, localPath)
	}

	outFile := filepath.Join(videoDir, "full.mp4")
	if err = utils.MergeVideos(videoPaths, outFile); err != nil {
		taskColl.UpdateByID(context.TODO(), task.ID, bson.M{"$set": bson.M{
			"status":  model.TaskStatusFailed,
			"message": fmt.Sprintf("视频合成失败: %v", err),
		}})
		return
	}

	taskColl.UpdateByID(context.TODO(), task.ID, bson.M{"$set": bson.M{
		"status":    model.TaskStatusSucceeded,
		"video_url": os.Getenv("APP_BASE_URL") + "/" + outFile,
	}})
}

func ListTasks(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		coll := client.Database("video_tasks").Collection("tasks")
		cur, err := coll.Find(context.TODO(),
			bson.M{"user_id": userID, "deleted_at": bson.M{"$exists": false}},
			options.Find().SetSort(bson.M{"created_at": -1}),
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
			return
		}
		defer cur.Close(context.TODO())
		var tasks []model.Task
		for cur.Next(context.TODO()) {
			var task model.Task
			cur.Decode(&task)
			tasks = append(tasks, task)
		}
		c.JSON(http.StatusOK, gin.H{"tasks": tasks})
	}
}

func GetTask(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		id := c.Param("id")
		coll := client.Database("video_tasks").Collection("tasks")
		var task model.Task
		if err := coll.FindOne(context.TODO(), bson.M{"_id": id, "user_id": userID, "deleted_at": bson.M{"$exists": false}}).Decode(&task); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"task": task})
	}
}

func DeleteTask(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		id := c.Param("id")
		coll := client.Database("video_tasks").Collection("tasks")
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
