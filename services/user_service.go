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

func NewUserService(repo UserRepository, jwtSecret string) *UserService {
    return &UserService{
        repo:      repo,
        jwtSecret: jwtSecret,
    }
}

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

func (s *UserService) Delete(userID string) error {
    return s.repo.Delete(userID)
}

func (s *UserService) GetByID(userID string) (*models.User, error) {
    return s.repo.GetByID(userID)
}