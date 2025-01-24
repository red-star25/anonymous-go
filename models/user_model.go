package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	UserName     *string            `json:"username" validate:"required,min=3,max=50"`
	Password     *string            `json:"password" validate:"required,min=6,max=50"`
	Image        *string            `json:"image"`
	Token        *string            `json:"token"`
	RefreshToken *string            `json:"refresh_token"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
	UserID       string             `json:"user_id"`
	UserPosts    []Post             `json:"userpost"`
}

type Post struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Title     *string            `json:"title"`
	Body      *string            `json:"body"`
	ByUser    string             `json:"created_by_user"`
	Likes     []Likes            `json:"likes"`
	Comments  []Comments         `json:"comments"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type Likes struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	LikeStatus *bool              `json:"like_status"`
	ByUser     string             `json:"by_user"`
}

type Comments struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	CommentBody *string            `json:"comment_body"`
	ByUser      string             `json:"by_user"`
}
