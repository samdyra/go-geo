package user

import (
	"github.com/jmoiron/sqlx"

	"github.com/samdyra/go-geo/internal/utils/errors"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
    db *sqlx.DB
}

func NewAuthService(db *sqlx.DB) *AuthService {
    return &AuthService{db: db}
}

func (s *AuthService) CreateUser(input SignUpInput) error {
    var existingUser User
    err := s.db.Get(&existingUser, "SELECT id FROM users WHERE username = $1", input.Username)
    if err == nil {
        return errors.ErrUserAlreadyExists
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        return errors.ErrInternalServer
    }

    _, err = s.db.Exec(
        "INSERT INTO users (username, password) VALUES ($1, $2)",
        input.Username, string(hashedPassword),
    )
    if err != nil {
        return errors.ErrInternalServer
    }

    return nil
}

func (s *AuthService) ValidateUser(input SignInInput) (*User, error) {
    var user User
    err := s.db.Get(&user, "SELECT * FROM users WHERE username = $1", input.Username)
    if err != nil {
        return nil, errors.ErrUserNotFound
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
    if err != nil {
        return nil, errors.ErrInvalidCredentials
    }

    return &user, nil
}