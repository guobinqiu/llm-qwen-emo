package model

import "time"

type Image struct {
	ID        string     `bson:"_id" json:"id"`
	UserID    string     `bson:"user_id" json:"user_id"`
	URL       string     `bson:"url" json:"url"`
	Filename  string     `bson:"filename" json:"filename"`
	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}
