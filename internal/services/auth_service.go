package services

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/samdyra/go-geo/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
    db *sqlx.DB
}

func NewAuthService(db *sqlx.DB) *AuthService {
    return &AuthService{db: db}
}

func (s *AuthService) CreateUser(username, password string) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    _, err = s.db.Exec(
        "INSERT INTO users (username, password, created_at, updated_at) VALUES ($1, $2, $3, $3)",
        username, string(hashedPassword), time.Now(),
    )
    return err
}

func (s *AuthService) ValidateUser(username, password string) (*models.User, error) {
    var user models.User
    err := s.db.Get(&user, "SELECT * FROM users WHERE username = $1", username)
    if err != nil {
        return nil, err
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        return nil, err
    }

    return &user, nil
}