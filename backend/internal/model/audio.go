package model

import "time"

type Audio struct {
	ID        string     `bson:"_id" json:"id"`
	UserID    string     `bson:"user_id" json:"user_id"`
	URL       string     `bson:"url" json:"url"`
	Filename  string     `bson:"filename" json:"filename"`
	Duration  int        `bson:"duration" json:"duration"`
	Size      int64      `bson:"size" json:"size"`
	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}
