package model

import "time"

const (
	TaskStatusPending        = "PENDING"
	TaskStatusPreProcessing  = "PRE-PROCESSING"
	TaskStatusRunning        = "RUNNING"
	TaskStatusPostProcessing = "POST-PROCESSING"
	TaskStatusSucceeded      = "SUCCEEDED"
	TaskStatusFailed         = "FAILED"
	TaskStatusUnknown        = "UNKNOWN"
)

type Task struct {
	ID        string     `bson:"_id" json:"id"`
	UserID    string     `bson:"user_id" json:"user_id"`
	Name      string     `bson:"name" json:"name"`
	ImageID   string     `bson:"image_id" json:"image_id"` // 冗余一下避免查询
	AudioID   string     `bson:"audio_id" json:"audio_id"`
	Status    string     `bson:"status" json:"status"`
	ImageURL  string     `bson:"image_url" json:"image_url"`
	VideoURL  string     `bson:"video_url" json:"video_url"` // 本地视频url(合成后)
	Message   string     `bson:"message,omitempty" json:"message,omitempty"`
	SubTasks  []SubTask  `bson:"sub_tasks" json:"sub_tasks"`
	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

type SubTask struct {
	AudioURL      string `bson:"audio_url" json:"audio_url"`
	TaskID        string `bson:"task_id" json:"task_id"`
	TaskStatus    string `bson:"task_status" json:"task_status"`
	ScheduledTime string `bson:"scheduled_time" json:"scheduled_time"`
	EndTime       string `bson:"end_time" json:"end_time"`
	VideoURL      string `bson:"video_url" json:"video_url"` // oss视频url(片段)
	Code          string `bson:"code,omitempty" json:"code,omitempty"`
	Message       string `bson:"message,omitempty" json:"message,omitempty"`
}
