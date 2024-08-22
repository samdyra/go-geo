package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samdyra/go-geo/internal/utils/auth"
	"github.com/samdyra/go-geo/internal/utils/errors"
)

type Handler struct {
    authService *AuthService
}

func NewHandler(authService *AuthService) *Handler {
    return &Handler{authService: authService}
}

func (h *Handler) SignUp(c *gin.Context) {
    var input SignUpInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, errors.NewAPIError(errors.ErrInvalidInput))
        return
    }

    if err := input.Validate(); err != nil {
        c.JSON(http.StatusBadRequest, errors.NewAPIError(err))
        return
    }

    err := h.authService.CreateUser(input)
    if err != nil {
        c.JSON(http.StatusInternalServerError, errors.NewAPIError(err))
        return
    }

    c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

func (h *Handler) SignIn(c *gin.Context) {
    var input SignInInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, errors.NewAPIError(errors.ErrInvalidInput))
        return
    }

    if err := input.Validate(); err != nil {
        c.JSON(http.StatusBadRequest, errors.NewAPIError(err))
        return
    }

    user, err := h.authService.ValidateUser(input)
    if err != nil {
        c.JSON(http.StatusUnauthorized, errors.NewAPIError(err))
        return
    }

    token, err := auth.GenerateToken(user.ID, user.Username)
    if err != nil {
        c.JSON(http.StatusInternalServerError, errors.NewAPIError(errors.ErrInternalServer))
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *Handler) Logout(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *Handler) ProtectedRoute(c *gin.Context) {
    userID, _ := c.Get("user_id")
    c.JSON(http.StatusOK, gin.H{"message": "This is a protected route", "user_id": userID})
}