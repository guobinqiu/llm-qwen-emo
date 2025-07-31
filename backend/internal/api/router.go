package app

import (
	"github.com/guobinqiu/llm-qwen-emo/internal/api/handler"
	"github.com/guobinqiu/llm-qwen-emo/internal/api/middleware"
	"github.com/guobinqiu/llm-qwen-emo/pkg/oss"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(r *gin.Engine, mongoClient *mongo.Client, ossClient *oss.OSSClient) {
	r.POST("/api/register", handler.Register(mongoClient))
	r.POST("/api/login", handler.Login(mongoClient))
	r.GET("/api/user", middleware.AuthMiddleware(), handler.GetUserInfo(mongoClient))

	// r.POST("/api/images/upload", middleware.AuthMiddleware(), handler.UploadImage(mongoClient))
	r.POST("/api/images/upload", middleware.AuthMiddleware(), handler.UploadImage(mongoClient, ossClient))
	r.GET("/api/images", middleware.AuthMiddleware(), handler.ListImages(mongoClient))
	r.DELETE("/api/images/:id", middleware.AuthMiddleware(), handler.DeleteImage(mongoClient))

	r.POST("/api/audios/upload", middleware.AuthMiddleware(), handler.UploadAudio(mongoClient))
	r.GET("/api/audios", middleware.AuthMiddleware(), handler.ListAudios(mongoClient))
	r.DELETE("/api/audios/:id", middleware.AuthMiddleware(), handler.DeleteAudio(mongoClient))

	r.POST("/api/tasks", middleware.AuthMiddleware(), handler.CreateTask(mongoClient, ossClient))
	r.GET("/api/tasks", middleware.AuthMiddleware(), handler.ListTasks(mongoClient))
	r.GET("/api/tasks/:id", middleware.AuthMiddleware(), handler.GetTask(mongoClient))
	r.DELETE("/api/tasks/:id", middleware.AuthMiddleware(), handler.DeleteTask(mongoClient))
}
