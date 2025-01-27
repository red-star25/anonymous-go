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
	User_Posts    []Post             `json:"user_posts"`
}

type Post struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	Title      *string            `json:"title"`
	Body       *string            `json:"body" validate:"required"`
	User_ID    string             `json:"user_id" validate:"required"`
	Likes      []Likes            `json:"likes"`
	Comments   []Comments         `json:"comments"`
	Created_At time.Time          `json:"created_at"`
	Updated_At time.Time          `json:"updated_at"`
}

type Likes struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Like_Status *bool              `json:"is_liked"`
	User_ID     string             `json:"user_id"`
}

type Comments struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	Comment_Body *string            `json:"comment"`
	User_ID      string             `json:"user_id"`
}
