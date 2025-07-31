package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	api "github.com/guobinqiu/llm-qwen-emo/internal/api"
	"github.com/guobinqiu/llm-qwen-emo/pkg/oss"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}

	ossClient, err := oss.NewOSSClient(
		os.Getenv("OSS_ENDPOINT"),
		os.Getenv("OSS_ACCESS_KEY_ID"),
		os.Getenv("OSS_ACCESS_KEY_SECRET"),
		os.Getenv("OSS_BUCKET"),
	)
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	config := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}
	r.Use(cors.New(config))

	r.Static("/images", "./images")
	r.Static("/audios", "./audios")
	r.Static("/videos", "./videos")

	api.RegisterRoutes(r, mongoClient, ossClient)
	r.Run(":8080")
}
