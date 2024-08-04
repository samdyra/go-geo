package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samdyra/go-geo/internal/auth"
	"github.com/samdyra/go-geo/internal/services"
)

type Handler struct {
    authService *services.AuthService
}

func NewHandler(authService *services.AuthService) *Handler {
    return &Handler{authService: authService}
}

func (h *Handler) SignUp(c *gin.Context) {
    var input struct {
        Username string `json:"username" binding:"required"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    err := h.authService.CreateUser(input.Username, input.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (h *Handler) SignIn(c *gin.Context) {
    var input struct {
        Username string `json:"username" binding:"required"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := h.authService.ValidateUser(input.Username, input.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    token, err := auth.GenerateToken(user.ID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *Handler) Logout(c *gin.Context) {
    // In a stateless JWT setup, logout is typically handled client-side
    // by removing the token from storage. Server-side, we can't invalidate
    // the token, but we can implement a token blacklist if needed.
    c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *Handler) ProtectedRoute(c *gin.Context) {
    userID, _ := c.Get("user_id")
    c.JSON(http.StatusOK, gin.H{"message": "This is a protected route", "user_id": userID})
}