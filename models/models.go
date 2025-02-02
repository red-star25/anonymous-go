package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `json:"_id" bson:"_id"`
	Username      *string            `json:"username" validate:"required,min=3,max=50"`
	Password      *string            `json:"password" validate:"required,min=6,max=50"`
	Image         *string            `json:"image"`
	Token         *string            `json:"token"`
	Refresh_Token *string            `json:"refresh_token"`
	Created_At    time.Time          `json:"created_at"`
	Updated_At    time.Time          `json:"updated_at"`
	User_ID       string             `json:"user_id"`
	User_Posts    []string           `json:"user_posts"`
}

type Post struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	Title      *string            `json:"title"`
	Body       *string            `json:"body" validate:"required"`
	User_ID    string             `json:"user_id" validate:"required"`
	Likes      []UserID           `json:"likes"`
	Comments   []Comment          `json:"comments"`
	Created_At time.Time          `json:"created_at"`
	Updated_At time.Time          `json:"updated_at"`
}

type UserID struct {
	User_ID string `json:"user_id"`
}

type Comment struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id" validate:"required"`
	Comment_Body *string            `json:"comment" validate:"required"`
	Created_At   time.Time          `json:"created_at"`
}
