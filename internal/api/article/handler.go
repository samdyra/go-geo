package article

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/samdyra/go-geo/internal/utils/errors"
)

type ArticleHandler struct {
	articleService *ArticleService
}

func NewArticleHandler(articleService *ArticleService) *ArticleHandler {
	return &ArticleHandler{articleService: articleService}
}

func (h *ArticleHandler) GetArticles(c *gin.Context) {
	articles, err := h.articleService.GetArticles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIError(err))
		return
	}
	c.JSON(http.StatusOK, articles)
}

func (h *ArticleHandler) GetArticle(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(errors.ErrInvalidInput))
		return
	}

	article, err := h.articleService.GetArticleByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, errors.NewAPIError(err))
		return
	}
	c.JSON(http.StatusOK, article)
}

func (h *ArticleHandler) CreateArticle(c *gin.Context) {
    var input CreateArticleInput
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, errors.NewAPIError(errors.ErrInvalidInput))
        return
    }

    if err := input.Validate(); err != nil {
        c.JSON(http.StatusBadRequest, errors.NewAPIError(err))
        return
    }

    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, errors.NewAPIError(errors.ErrUnauthorized))
        return
    }

    userIDInt, ok := userID.(int64)
    if !ok {
        c.JSON(http.StatusInternalServerError, errors.NewAPIError(errors.ErrInternalServer))
        return
    }

    article, err := h.articleService.CreateArticle(input, userIDInt)
    if err != nil {
        c.JSON(http.StatusInternalServerError, errors.NewAPIError(err))
        return
    }

    c.JSON(http.StatusCreated, article)
}

func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(errors.ErrInvalidInput))
		return
	}

	var input UpdateArticleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(errors.ErrInvalidInput))
		return
	}

	if err := input.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(err))
		return
	}

	userID, _ := c.Get("user_id")
	article, err := h.articleService.UpdateArticle(id, input, userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIError(err))
		return
	}

	c.JSON(http.StatusOK, article)
}

func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.NewAPIError(errors.ErrInvalidInput))
		return
	}

	userID, _ := c.Get("user_id")
	err = h.articleService.DeleteArticle(id, userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.NewAPIError(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article deleted successfully"})
}