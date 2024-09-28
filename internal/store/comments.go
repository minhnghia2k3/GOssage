package store

import (
	"context"
	"database/sql"
)

type IComments interface {
	Create(context.Context, *Comment) error
	GetByPostID(context.Context, int64) ([]Comment, error)
}

type Comment struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	PostID    int64  `json:"post_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	User      struct {
		ID       int64  `json:"id"`
		Username string `json:"name"`
	} `json:"user"`
}

type CommentStorage struct {
	db *sql.DB
}

func (s *CommentStorage) Create(ctx context.Context, c *Comment) error {
	query := `
	INSERT INTO comments (user_id, post_id, content)
	VALUES ($1, $2, $3)
	RETURNING id, created_at
`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		c.UserID,
		c.PostID,
		c.Content,
	).Scan(&c.ID, &c.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (s *CommentStorage) GetByPostID(ctx context.Context, id int64) ([]Comment, error) {
	query := `
	SELECT c.id, user_id, post_id, content, users.username, users.id, c.created_at
	FROM comments c
	JOIN users ON c.user_id = users.id
	WHERE post_id = $1
	ORDER BY c.created_at DESC;
`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var comments []Comment
	rows, err := s.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c Comment
		if err = rows.Scan(
			&c.ID,
			&c.UserID,
			&c.PostID,
			&c.Content,
			&c.User.Username,
			&c.User.ID,
			&c.CreatedAt,
		); err != nil {
			return nil, err
		}

		comments = append(comments, c)
	}
	return comments, nil
}
