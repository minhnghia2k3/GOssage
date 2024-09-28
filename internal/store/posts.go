package store

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"time"
)

type IPosts interface {
	GetByID(context.Context, int64) (*Post, error)
	Create(context.Context, *Post) error
	Update(context.Context, *Post) error
	Delete(context.Context, int64) error
}

// Post model
type Post struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int       `json:"version"`
	Comments  []Comment `json:"comments,omitempty"`
}

type PostStorage struct {
	db *sql.DB
}

// GetByID gets a post by given ID, return a pointer to Post.
// If the query select no rows, the function will return ErrNotFound.
func (s *PostStorage) GetByID(ctx context.Context, id int64) (*Post, error) {
	var post Post

	query := `
	SELECT id, title, user_id, content, tags, created_at, updated_at, version 
	FROM posts 
	WHERE id = $1;
`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.Title,
		&post.UserID,
		&post.Content,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &post, nil
}

// Create creates a post with provided data, scan return data into Post instance.
func (s *PostStorage) Create(ctx context.Context, post *Post) error {
	query := `
	INSERT INTO posts (user_id, title, content, tags)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, updated_at
`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.UserID,
		post.Title,
		post.Content,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

// Update updates a post with specific ID, scan return data into Post instance
// or return ErrNotFound if there is no rows in query.
func (s *PostStorage) Update(ctx context.Context, post *Post) error {
	query := `
	UPDATE posts 
	SET title = $1, content = $2, version = version + 1 
	WHERE id = $3 and version = $4	
	RETURNING version
`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Content,
		post.ID,
		post.Version,
	).Scan(&post.Version)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}

// Delete deletes a post instance by given ID.
func (s *PostStorage) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM posts WHERE id = $1;`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
