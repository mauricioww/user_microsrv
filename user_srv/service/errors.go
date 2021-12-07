package service

import (
	"errors"
	"fmt"

	"github.com/mauricioww/user_microsrv/user_srv/repository"
)

type (
	USER_NOT_FOUND struct {
		Err error
	}
	INTERNAL_ERROR struct {
		Err error
	}
	UNKNOWN_ERROR struct {
		Err error
	}
	UNAUTHENTICATED_USER struct {
		Err error
	}
)

func (e USER_NOT_FOUND) Error() string {
	return fmt.Sprintf("%v", e.Err)
}

func (e INTERNAL_ERROR) Error() string {
	return fmt.Sprintf("%v", e.Err)
}

func (e UNKNOWN_ERROR) Error() string {
	return fmt.Sprintf("%v", e.Err)
}

func (e UNAUTHENTICATED_USER) Error() string {
	return fmt.Sprintf("%v", e.Err)
}

func RepoError(err error) error {
	switch e := err.(type) {
	case repository.USER_NOT_FOUND:
		return USER_NOT_FOUND{Err: e.Err}
	case repository.INTERNAL_ERROR:
		return INTERNAL_ERROR{Err: e.Err}
	default:
		return UNKNOWN_ERROR{Err: errors.New("Unknown error")}
	}
}
