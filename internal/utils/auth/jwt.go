package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// @TODO: Use an environment variable for the secret key
var jwtKey = []byte("your_secret_key") 

func GenerateToken(userID int64, username string) (string, error) {
    claims := jwt.MapClaims{
        "user_id":  userID,
        "username": username,
        "exp":      time.Now().Add(time.Hour * 24).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}

func ValidateToken(tokenString string) (int64, string, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return jwtKey, nil
    })

    if err != nil {
        return 0, "", err
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        userID, ok := claims["user_id"].(float64)
        if !ok {
            return 0, "", fmt.Errorf("invalid user_id in token")
        }

        username, ok := claims["username"].(string)
        if !ok {
            return 0, "", fmt.Errorf("invalid username in token")
        }

        return int64(userID), username, nil
    }

    return 0, "", fmt.Errorf("invalid token")
}