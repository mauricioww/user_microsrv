package repository

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/go-kit/kit/log"
	"github.com/mauricioww/user_microsrv/user_srv/entities"
)

const (
	create_user_sql = `
		INSERT INTO usr_service(email, password, age, extra_info)
			VALUES (?, ?, ? , ?)
	`
)

type UserSrvRepository interface {
	CreateUser(ctx context.Context, user entities.User) (string, error)
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
		return "", err
	}

	n, _ := id.LastInsertId()

	return strconv.FormatInt(n, 10), nil
}
