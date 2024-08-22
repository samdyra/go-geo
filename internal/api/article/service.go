package article

import (
	"database/sql"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/samdyra/go-geo/internal/utils/errors"
)

type ArticleService struct {
	db *sqlx.DB
}

func NewArticleService(db *sqlx.DB) *ArticleService {
	return &ArticleService{db: db}
}

func (s *ArticleService) GetArticles() ([]Article, error) {
	var articles []Article
	err := s.db.Select(&articles, "SELECT * FROM articles ORDER BY created_at DESC")
	if err != nil {
		return nil, errors.ErrInternalServer
	}
	return articles, nil
}

func (s *ArticleService) GetArticleByID(id int64) (*Article, error) {
	var article Article
	err := s.db.Get(&article, "SELECT * FROM articles WHERE id = $1", id)
	if err == sql.ErrNoRows {
		return nil, errors.ErrNotFound
	}
	if err != nil {
		return nil, errors.ErrInternalServer
	}
	return &article, nil
}

func (s *ArticleService) CreateArticle(input CreateArticleInput, userID int64) (*Article, error) {
    var username string
    err := s.db.Get(&username, "SELECT username FROM users WHERE id = $1", userID)
    if err != nil {
        log.Printf("Error fetching username: %v", err)
        return nil, errors.ErrInternalServer
    }

    article := &Article{
        Title:     input.Title,
        Content:   input.Content,
        ImageURL:  input.ImageURL,
        Author:    username,
        CreatedBy: userID,
        CreatedAt: time.Now(),
    }

    query := `INSERT INTO articles (title, content, image_url, author, created_by, created_at) 
              VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
    err = s.db.QueryRow(query, article.Title, article.Content, article.ImageURL, article.Author, article.CreatedBy, article.CreatedAt).
        Scan(&article.ID)
    if err != nil {
        log.Printf("Error creating article: %v", err)
        return nil, errors.ErrInternalServer
    }

    return article, nil
}

func (s *ArticleService) UpdateArticle(id int64, input UpdateArticleInput, userID int64) (*Article, error) {
	article, err := s.GetArticleByID(id)
	if err != nil {
		return nil, err
	}

	if article.CreatedBy != userID {
		return nil, errors.ErrUnauthorized
	}

	if input.Title != nil {
		article.Title = *input.Title
	}
	if input.Content != nil {
		article.Content = *input.Content
	}
	if input.ImageURL != nil {
		article.ImageURL = input.ImageURL
	}

	query := `UPDATE articles SET title = $1, content = $2, image_url = $3 WHERE id = $4`
	_, err = s.db.Exec(query, article.Title, article.Content, article.ImageURL, id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	return article, nil
}

func (s *ArticleService) DeleteArticle(id int64, userID int64) error {
	article, err := s.GetArticleByID(id)
	if err != nil {
		return err
	}

	if article.CreatedBy != userID {
		return errors.ErrUnauthorized
	}

	_, err = s.db.Exec("DELETE FROM articles WHERE id = $1", id)
	if err != nil {
		return errors.ErrInternalServer
	}

	return nil
}