package repository

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	"github.com/go-kit/kit/log"
	"github.com/mauricioww/user_microsrv/user_srv/entities"
)

const (
	create_user_sql = `
		INSERT INTO usr_service(email, password, age, extra_info)
			VALUES (?, ?, ? , ?)
	`

	authenticate_sql = `
		SELECT u.password FROM usr_service u WHERE u.email = ?
	`
)

type UserSrvRepository interface {
	CreateUser(ctx context.Context, user entities.User) (string, error)
	Authenticate(ctx context.Context, session entities.Session) (string, error)
}

type userSrvRepository struct {
	db     *sql.DB
	logger log.Logger
}

func NewUserSrvRepository(db *sql.DB, l log.Logger) UserSrvRepository {
	return &userSrvRepository{
		db:     db,
		logger: log.With(l, "repository", "mysql"),
	}
}

func (r userSrvRepository) CreateUser(ctx context.Context, user entities.User) (string, error) {
	id, err := r.db.ExecContext(ctx, create_user_sql, user.Email, user.Password, user.Age, user.ExtraInfo)

	if err != nil {
		return "", INTERNAL_ERROR{Err: errors.New("Internal Error")}
	}

	n, _ := id.LastInsertId()

	return strconv.FormatInt(n, 10), nil
}

func (r userSrvRepository) Authenticate(ctx context.Context, session entities.Session) (string, error) {
	var hash string
	err := r.db.QueryRow(authenticate_sql, session.Email).Scan(&hash)

	if err == sql.ErrNoRows {
		return "", USER_NOT_FOUND{Err: errors.New("User not found")}
	}

	if err != nil {
		return "", INTERNAL_ERROR{Err: errors.New("Internal Error")}
	}

	return hash, nil
}
