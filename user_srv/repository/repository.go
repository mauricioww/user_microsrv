package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
		SELECT u.id, u.password FROM usr_service u WHERE u.email = ?
	`

	update_user_field_sql = `
		UPDATE usr_service SET %v = ? WHERE id = ? 
	`

	get_user_by_id = `
		SELECT u.email, u.password, u.age, u.extra_info 
			FROM user_service u WHERE u.id = ?
	`
)

type UserSrvRepository interface {
	CreateUser(ctx context.Context, user entities.User) (string, error)
	Authenticate(ctx context.Context, session *entities.Session) (string, error)
	UpdateUser(ctx context.Context, information entities.Update) (entities.User, error)
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
		return "", errors.New("Internal Error")
	}

	n, _ := id.LastInsertId()

	return strconv.FormatInt(n, 10), nil
}

func (r userSrvRepository) Authenticate(ctx context.Context, session *entities.Session) (string, error) {
	var hash string
	err := r.db.QueryRow(authenticate_sql, session.Email).Scan(&session.Id, &hash)

	if err == sql.ErrNoRows {
		return "", errors.New("User not found")
	}

	if err != nil {
		return "", errors.New("Internal Error")
	}

	return hash, nil
}

func (r userSrvRepository) UpdateUser(ctx context.Context, update entities.Update) (entities.User, error) {
	for field, value := range update.Information {
		query := fmt.Sprintf(update_user_field_sql, field)
		_, err := r.db.ExecContext(ctx, query, value, update.UserId)

		if err != nil {
			return entities.User{}, err
		}
	}

	var u entities.User
	_ = r.db.QueryRow(get_user_by_id, update.UserId).Scan(&u.Email, &u.Password, &u.Age, &u.ExtraInfo)

	return u, nil
}
