package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type IUsers interface {
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, tx *sql.Tx, user *User) error
	CreateAndInvite(ctx context.Context, user *User, token string, expiryDuration time.Duration) error
	Activate(ctx context.Context, token string) error
	Delete(ctx context.Context, id int64) error
}

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsActive  bool      `json:"is_active"`
}

type password struct {
	text *string
	hash []byte
}

type UserStorage struct {
	db *sql.DB
}

func (s *UserStorage) Create(ctx context.Context, tx *sql.Tx, users *User) error {
	query := `
	INSERT INTO users (username, email, password)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, updated_at
`
	err := tx.QueryRowContext(
		ctx,
		query,
		users.Username,
		users.Email,
		users.Password.hash,
	).Scan(
		&users.ID,
		&users.CreatedAt,
		&users.UpdatedAt,
	)

	if err != nil {
		var pqError *pq.Error

		switch errors.As(err, &pqError) {
		case pqError.Code == "23505":
			log.Println(err.Error())
			return ErrConflict

		default:
			return err
		}
	}

	return nil
}

func (s *UserStorage) CreateAndInvite(ctx context.Context, user *User, token string, expiry time.Duration) error {
	return withTx(ctx, s.db, func(tx *sql.Tx) error {
		// 1. Create a user
		if err := s.Create(ctx, tx, user); err != nil {
			return err
		}

		// 2. insert token with expiry time to user_invitations
		if err := s.createUserInvitation(ctx, tx, user, token, expiry); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStorage) Activate(ctx context.Context, token string) error {
	return withTx(ctx, s.db, func(tx *sql.Tx) error {
		// 1. Get user using token
		user, err := s.getUserFromInvitation(ctx, tx, token)
		if err != nil {
			return err
		}

		// 2. Update is_active status
		user.IsActive = true
		if err = s.update(ctx, tx, user); err != nil {
			return err
		}

		// 3. Clear user invitation
		if err = s.deleteUserInvitation(ctx, tx, user.ID); err != nil {
			return err
		}

		return nil
	})
}

func (s *UserStorage) GetByID(ctx context.Context, id int64) (*User, error) {
	query := `
	SELECT id, username, email, created_at, updated_at
	FROM users
	WHERE id = $1 AND is_active=true
`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var user User
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (s *UserStorage) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
	SELECT id, username, email, created_at, updated_at, is_active
	FROM users 
	WHERE email = $1 AND is_active=true
`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var u User
	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.IsActive,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &u, nil
}

func (s *UserStorage) Delete(ctx context.Context, id int64) error {
	// Delete user and its invitations
	return withTx(ctx, s.db, func(tx *sql.Tx) error {
		err := s.deleteUserInvitation(ctx, tx, id)
		if err != nil {
			return err
		}

		err = s.delete(ctx, tx, id)
		if err != nil {
			return err
		}

		return nil
	})
}

func (p *password) Set(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}

	p.text = &password
	p.hash = bytes
	return nil
}

func (s *UserStorage) delete(ctx context.Context, tx *sql.Tx, id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	return nil
}

func (s *UserStorage) getUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*User, error) {
	var user User

	query := `
	SELECT u.id, u.is_active
	FROM users u
	INNER JOIN user_invitation ui ON u.id = ui.user_id 
	WHERE ui.token = $1 AND ui.expiry > $2
`

	// Hash plain text password for comparing
	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	err := tx.QueryRowContext(ctx, query, hashToken, time.Now()).Scan(&user.ID, &user.IsActive)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (s *UserStorage) update(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `UPDATE users SET is_active = true WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStorage) deleteUserInvitation(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := `DELETE FROM user_invitation WHERE user_id = $1`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserStorage) createUserInvitation(ctx context.Context, tx *sql.Tx, user *User, token string, expiry time.Duration) error {
	query := `INSERT INTO user_invitation (user_id, token, expiry) VALUES ($1, $2, $3)`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, user.ID, token, time.Now().Add(expiry))
	if err != nil {
		return err
	}

	return nil
}
