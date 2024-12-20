package repositories

import (
    "context"
    "time"
    "go-blog-backend/models"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository struct {
    collection *mongo.Collection
}

// NewUserRepository returns a new instance of UserRepository.
//
// The UserRepository is used to interact with the "users" collection in the MongoDB
// database.
func NewUserRepository(db *mongo.Database) *UserRepository {
    return &UserRepository{
        collection: db.Collection("users"),
    }
}

// Create creates a new user in the "users" collection in the MongoDB database.
//
// The user struct passed in should have a nil ID, as the ID is automatically
// generated by the MongoDB driver.
//
// The user struct passed out will have its ID field populated with the
// generated ID.
func (r *UserRepository) Create(user *models.User) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    result, err := r.collection.InsertOne(ctx, user)
    if err != nil {
        return err
    }

    user.ID = result.InsertedID.(primitive.ObjectID)
    return nil
}

// GetByEmail returns a user by the given email.
//
// If no user is found with the given email, the returned error will be nil,
// and the returned user will be nil.
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var user models.User
    err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
    if err != nil {
        return nil, err
    }

    return &user, nil
}

// GetByID returns a user by the given ID.
//
func (r *UserRepository) GetByID(id string) (*models.User, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return nil, err
    }

    var user models.User
    err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
    if err != nil {
        return nil, err
    }

    return &user, nil
}

// Update updates the user with the given ID in the "users" collection in the
// MongoDB database.
//
// The updates parameter is a map of key-value pairs that should be updated in
// the user document. The key is the field name and the value is the new value
// for that field.
//
// The returned error will be non-nil if any error occurred during the update
// process.
func (r *UserRepository) Update(id string, updates map[string]interface{}) error {
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

// Delete deletes the user with the given ID from the "users" collection in the
// MongoDB database.
//
// The returned error will be non-nil if any error occurred during the delete
// process.
func (r *UserRepository) Delete(id string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    objectID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err
    }

    _, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
    return err
}