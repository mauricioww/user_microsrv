package repository

import (
	"context"

	"github.com/mauricioww/user_microsrv/user_srv/entities"
)

type UserSrvRepository interface {
	CreateUser(context.Context, entities.User) (string, error)
}
