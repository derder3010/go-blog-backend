package services

import (
    "go-blog-backend/models"
    "golang.org/x/crypto/bcrypt"
    "github.com/golang-jwt/jwt/v4"
    "time"
    "errors"
)

type UserRepository interface {
    Create(user *models.User) error
    GetByEmail(email string) (*models.User, error)
    GetByID(id string) (*models.User, error)
    Update(id string, updates map[string]interface{}) error
    Delete(id string) error
}

type UserService struct {
    repo      UserRepository
    jwtSecret string
}

// NewUserService creates a new UserService instance with the given UserRepository and JWT secret.
//
// Parameters:
//   - repo: The UserRepository interface used for interacting with the user data storage.
//   - jwtSecret: The secret key used for signing JWT tokens.
//
// Returns a pointer to a UserService instance.
func NewUserService(repo UserRepository, jwtSecret string) *UserService {
    return &UserService{
        repo:      repo,
        jwtSecret: jwtSecret,
    }
}

// Register creates a new user in the "users" collection in the MongoDB database.
//
// If the email address is already registered, an error will be returned.
//
// Parameters:
//   - username: The username for the new user.
//   - email: The email address for the new user.
//   - password: The password for the new user.
//
// Returns a pointer to the newly created User instance, or an error if any error occurred during the registration process.
func (s *UserService) Register(username, email, password string) (*models.User, error) {
    // Check if user already exists
    existing, err := s.repo.GetByEmail(email)
    if err == nil && existing != nil {
        return nil, errors.New("email already registered")
    }

    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    user := &models.User{
        Username:  username,
        Email:     email,
        Password:  string(hashedPassword),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    if err := s.repo.Create(user); err != nil {
        return nil, err
    }

    return user, nil
}

// Login authenticates a user by their email and password and returns a JWT token.
//
// Parameters:
//   - email: The email address to authenticate.
//   - password: The password to authenticate.
//
// Returns a JWT token string if the authentication is successful, or an error if any error occurred during the authentication process.
func (s *UserService) Login(email, password string) (string, error) {
    user, err := s.repo.GetByEmail(email)
    if err != nil {
        return "", errors.New("invalid credentials")
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return "", errors.New("invalid credentials")
    }

    // Generate JWT token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID.Hex(),
        "email":   user.Email,
        "exp":     time.Now().Add(time.Hour * 24).Unix(),
    })

    tokenString, err := token.SignedString([]byte(s.jwtSecret))
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

// Update updates the fields of the user with the given ID in the "users" collection.
//
// The updates parameter is a map of key-value pairs where the key is the field name
// and the value is the new value for that field. The updated_at field is automatically
// set to the current time.
//
// The returned error will be non-nil if any error occurred during the update process.
func (s *UserService) Update(userID string, updates map[string]interface{}) error {
    if password, ok := updates["password"].(string); ok {
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
        if err != nil {
            return err
        }
        updates["password"] = string(hashedPassword)
    }

    updates["updated_at"] = time.Now()
    return s.repo.Update(userID, updates)
}

// Delete deletes the user with the given ID from the "users" collection in the
// MongoDB database.
//
// The returned error will be non-nil if any error occurred during the delete
// process.
func (s *UserService) Delete(userID string) error {
    return s.repo.Delete(userID)
}

// GetByID returns a user by the given ID.
//
// The returned error will be non-nil if any error occurred during the get
// process.
func (s *UserService) GetByID(userID string) (*models.User, error) {
    return s.repo.GetByID(userID)
}