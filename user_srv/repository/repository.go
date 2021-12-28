package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-kit/kit/log"
	"github.com/mauricioww/user_microsrv/user_srv/entities"
)

const (
	create_user_sql = `
		INSERT INTO usr_service(email, password, age)
			VALUES (?, ?, ?)
	`

	authenticate_sql = `
		SELECT u.id, u.password FROM usr_service u WHERE u.email = ?
	`

	update_user_sql = `
		UPDATE usr_service SET email = ?, password = ?, age = ?
			WHERE id = ? 
	`

	get_user_by_id = `
		SELECT u.email, u.password, u.age
			FROM usr_service u WHERE u.id = ?
	`

	delete_user_sql = `
		DELETE FROM usr_service WHERE id = ?
	`
)

type UserRepository interface {
	CreateUser(ctx context.Context, user entities.User) (int, error)
	Authenticate(ctx context.Context, session *entities.Session) (string, error)
	UpdateUser(ctx context.Context, information entities.Update) (entities.User, error)
	GetUser(ctx context.Context, id int) (entities.User, error)
	DeleteUser(ctx context.Context, id int) (bool, error)
}

type userRepository struct {
	db     *sql.DB
	logger log.Logger
}

func NewUserRepository(mysql_db *sql.DB, l log.Logger) UserRepository {
	return &userRepository{
		db:     mysql_db,
		logger: log.With(l, "repository", "mysql"),
	}
}

func (r *userRepository) CreateUser(ctx context.Context, user entities.User) (int, error) {
	if id, err := r.db.ExecContext(ctx, create_user_sql, user.Email, user.Password, user.Age); err != nil {
		return -1, errors.New("Internal Error")
	} else {
		n, _ := id.LastInsertId()
		return int(n), nil
	}
}

func (r *userRepository) Authenticate(ctx context.Context, session *entities.Session) (string, error) {
	var hash string

	if err := r.db.QueryRow(authenticate_sql, session.Email).Scan(&session.Id, &hash); err == sql.ErrNoRows {
		return "", errors.New("User not found")
	} else if err != nil {
		return "", errors.New("Internal Error")
	}

	return hash, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, update entities.Update) (entities.User, error) {
	var u entities.User

	if _, err := r.db.ExecContext(ctx, update_user_sql, update.Email, update.Password, update.Age, update.UserId); err != nil {
		return u, errors.New("Internal Error")
	}

	_ = r.db.QueryRow(get_user_by_id, update.UserId).Scan(&u.Email, &u.Password, &u.Age)
	return u, nil
}

func (r *userRepository) GetUser(ctx context.Context, id int) (entities.User, error) {
	var u entities.User

	if err := r.db.QueryRow(get_user_by_id, id).Scan(&u.Email, &u.Password, &u.Age); err == sql.ErrNoRows {
		return entities.User{}, errors.New("User not found")
	} else if err != nil {
		return entities.User{}, errors.New("Internal Error")
	}

	return u, nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id int) (bool, error) {

	if err := r.db.QueryRow(get_user_by_id, id).Scan(); err == sql.ErrNoRows {
		return false, errors.New("User not found")
	} else if _, err := r.db.ExecContext(ctx, delete_user_sql, id); err != nil {
		return false, errors.New("Internal Error")
	}

	return true, nil
}
