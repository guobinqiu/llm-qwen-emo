package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI           string
	OssEndpoint        string
	OssBucket          string
	OssAccessKeyID     string
	OssAccessKeySecret string
	DashScopeApiKey    string
}

func LoadConfig() *Config {
	_ = godotenv.Load()
	return &Config{
		MongoURI:           os.Getenv("MONGO_URI"),
		OssEndpoint:        os.Getenv("OSS_ENDPOINT"),
		OssBucket:          os.Getenv("OSS_BUCKET"),
		OssAccessKeyID:     os.Getenv("OSS_ACCESS_KEY_ID"),
		OssAccessKeySecret: os.Getenv("OSS_ACCESS_KEY_SECRET"),
		DashScopeApiKey:    os.Getenv("DASHSCOPE_API_KEY"),
	}
}
