package models

import (
	"context"
	"database/sql"
	"errors"
	"math/rand"
	"time"
)

var AnonymousUser = &User{}

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"-"`
}

type UserModel struct {
	DB *sql.DB
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

func (m UserModel) Get(id int) (*User, error) {
	query := `
		SELECT id, name, email, password
		FROM users
		WHERE id = $1
	`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, NoDataFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (m UserModel) GetByEmail(email string) (*User, error) {
	query := `SELECT id, email, name, password FROM users WHERE email = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User
	err := m.DB.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Name, &user.Password)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, NoDataFound
		default:
			return nil, err
		}

	}

	return &user, nil
}

func (m UserModel) InsertWithRandomPassword(name, email string) (*User, error) {
	query := `
		INSERT INTO users (name, email, password)
		VALUES ($1,$2, $3)
		RETURNING id, name, email
	`
	args := []interface{}{name, email, RandStringBytes(8)}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return nil, ErrDuplicateEmail
		default:
			return nil, err
		}
	}

	return &user, nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
