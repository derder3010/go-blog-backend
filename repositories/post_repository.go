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

// NewPostRepository returns a new instance of PostRepository.
//
// The PostRepository is used to interact with the "posts" collection in the MongoDB
// database.
func NewPostRepository(db *mongo.Database) *PostRepository {
    return &PostRepository{
        collection: db.Collection("posts"),
    }
}

// Create creates a new post in the "posts" collection in the MongoDB database.
//
// The returned error will be non-nil if any error occurred during the create
// process.
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

// GetByID returns a post by the given ID.
//
// The returned error will be non-nil if any error occurred during the get
// process.
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

// Update updates the post with the given ID in the "posts" collection in the
// MongoDB database.
//
// The updates parameter is a map of key-value pairs that should be updated in
// the post document. The key is the field name and the value is the new value
// for that field.
//
// The returned error will be non-nil if any error occurred during the update
// process.
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

// Delete deletes the post with the given ID from the "posts" collection in the
// MongoDB database.
//
// The returned error will be non-nil if any error occurred during the delete
// process.
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

// List returns a slice of posts, sorted by created_at in descending order,
// limited to the given number of items, and starting from the given page.
//
// The page and limit parameters are 1-indexed, so the first page should have
// page = 1 and limit = 10 to get the first 10 results.
//
// The returned error will be non-nil if any error occurred during the find
// process.
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

// GetByAuthor returns a slice of posts, filtered by the given author ID.
//
// The returned error will be non-nil if any error occurred during the find
// process.
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