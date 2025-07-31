package handler

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/guobinqiu/llm-qwen-emo/internal/api/middleware"
	"github.com/guobinqiu/llm-qwen-emo/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Login(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
			return
		}
		coll := client.Database("video_tasks").Collection("users")
		var user model.User
		err := coll.FindOne(context.TODO(), bson.M{"username": req.Username}).Decode(&user)
		if err != nil || user.Password != hashPassword(req.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
			return
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, middleware.Claims{
			UserID: user.ID,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			},
		})
		tokenString, _ := token.SignedString(middleware.JwtSecret)
		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	}
}

func Register(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
			return
		}

		if req.Username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "用户名不能为空"})
			return
		}

		if req.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "密码不能为空"})
			return
		}

		if len(req.Username) < 3 || len(req.Username) > 20 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "用户名长度应在3到20个字符之间"})
			return
		}

		if len(req.Password) < 6 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "密码长度至少6位"})
			return
		}

		coll := client.Database("video_tasks").Collection("users")
		var exist model.User
		err := coll.FindOne(context.TODO(), bson.M{"username": req.Username}).Decode(&exist)
		if err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
			return
		}
		user := model.User{
			ID:        uuid.NewString(),
			Username:  req.Username,
			Password:  hashPassword(req.Password),
			CreatedAt: time.Now(),
		}
		_, err = coll.InsertOne(context.TODO(), user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"msg": "注册成功"})
	}
}

func hashPassword(pw string) string {
	h := sha256.New()
	h.Write([]byte(pw))
	return hex.EncodeToString(h.Sum(nil))
}

func GetUserInfo(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		coll := client.Database("video_tasks").Collection("users")
		var user model.User
		err := coll.FindOne(context.TODO(), bson.M{"_id": userID}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": user})
	}
}
