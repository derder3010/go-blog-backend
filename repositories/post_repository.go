package repositories

import (
    "context"
    "time"
    "go-blog-backend/models"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type PostRepository struct {
    collection *mongo.Collection
}

func NewPostRepository(db *mongo.Database) *PostRepository {
    return &PostRepository{
        collection: db.Collection("posts"),
    }
}

func (r *PostRepository) Create(post *models.Post) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    result, err := r.collection.InsertOne(ctx, post)
    if err != nil {
        return err
    }

    post.ID = result.InsertedID.(primitive.ObjectID)
    return nil
}

func (r *PostRepository) GetByID(id string) (*models.Post, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, err
    }

    var post models.Post
    err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&post)
    if err != nil {
        return nil, err
    }

    return &post, nil
}

func (r *PostRepository) Update(id string, updates map[string]interface{}) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err
    }

    _, err = r.collection.UpdateOne(
        ctx,
        bson.M{"_id": objectID},
        bson.M{"$set": updates},
    )
    return err
}

func (r *PostRepository) Delete(id string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err
    }

    _, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
    return err
}

func (r *PostRepository) List(page, limit int) ([]*models.Post, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    skip := (page - 1) * limit

    opts := options.Find().
        SetSort(bson.D{{Key: "created_at", Value: -1}}).
        SetSkip(int64(skip)).
        SetLimit(int64(limit))

    cursor, err := r.collection.Find(ctx, bson.M{}, opts)
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var posts []*models.Post
    if err = cursor.All(ctx, &posts); err != nil {
        return nil, err
    }

    return posts, nil
}

func (r *PostRepository) GetByAuthor(authorID string) ([]*models.Post, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    objectID, err := primitive.ObjectIDFromHex(authorID)
    if err != nil {
        return nil, err
    }

    cursor, err := r.collection.Find(ctx, bson.M{"author_id": objectID})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var posts []*models.Post
    if err = cursor.All(ctx, &posts); err != nil {
        return nil, err
    }

    return posts, nil
}