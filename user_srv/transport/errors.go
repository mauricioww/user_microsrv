package transport

import (
	"errors"
	"fmt"

	"github.com/mauricioww/user_microsrv/user_srv/service"
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
	INVALID_TYPE struct {
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

func (e INVALID_TYPE) Error() string {
	return fmt.Sprintf("%v", e.Err)
}

func SrvError(err error) error {
	switch e := err.(type) {
	case service.USER_NOT_FOUND:
		return USER_NOT_FOUND{Err: e.Err}
	case service.INTERNAL_ERROR:
		return INTERNAL_ERROR{Err: e.Err}
	case service.UNAUTHENTICATED_USER:
		return UNAUTHENTICATED_USER{Err: e.Err}
	default:
		return UNKNOWN_ERROR{Err: errors.New("Unknown error")}
	}
}
