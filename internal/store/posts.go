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
	GetUserFeed(context.Context, int64, PaginatedFeedQuery) ([]PostWithMetadata, error)
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
	User      struct {
		Username string `json:"username,omitempty"`
	} `json:"user,omitempty"`
}

type PostStorage struct {
	db *sql.DB
}

type PostWithMetadata struct {
	Post
	CommentCounts int `json:"comment_counts"`
}

// GetUserFeed gets posts from followed user and user itself,
// with associated username, and comment counts,
// limited by PaginatedFeedQuery
func (s *PostStorage) GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]PostWithMetadata, error) {
	// Get posts from followed user and user itself
	query := `
		SELECT 
			p.id, p.user_id, p.title, p.content, p.created_at, p.updated_at, p.version, p.tags,
			u.username,
			COUNT(c.id) AS comments_count
		FROM posts p
		LEFT JOIN comments c ON c.post_id = p.id
		LEFT JOIN users u ON p.user_id = u.id
		LEFT JOIN followers f ON f.user_id = p.user_id AND f.follower_id = $1
		WHERE (f.follower_id = $1 OR p.user_id = $1) AND (
		    (p.title ILIKE '%' || $4 || '%' OR p.content ILIKE '%' || $4 || '%'))
		GROUP BY p.id, u.username
		ORDER BY p.created_at ` + fq.Sort + `
		LIMIT $2 OFFSET $3
	`

	ctx, cancel := context.WithTimeout(context.Background(), QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userID, fq.Limit, fq.Offset, fq.Search)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feed []PostWithMetadata
	for rows.Next() {
		var post PostWithMetadata

		err = rows.Scan(
			&post.ID,
			&post.UserID,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.Version,
			pq.Array(&post.Tags),
			&post.User.Username,
			&post.CommentCounts,
		)

		if err != nil {
			return nil, err
		}

		feed = append(feed, post)
	}

	return feed, nil
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
