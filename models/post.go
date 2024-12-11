package models

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
    "time"
)

type Post struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Title       string            `bson:"title" json:"title"`
    Content     string            `bson:"content" json:"content"`
    AuthorID    primitive.ObjectID `bson:"author_id" json:"author_id"`
    ImageURL    string            `bson:"image_url" json:"image_url"`
    CreatedAt   time.Time         `bson:"created_at" json:"created_at"`
    UpdatedAt   time.Time         `bson:"updated_at" json:"updated_at"`
}