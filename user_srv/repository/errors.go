package repository

import (
	"fmt"
)

type (
	USER_NOT_FOUND struct {
		Err error
	}

	INTERNAL_ERROR struct {
		Err error
	}
)

func (e USER_NOT_FOUND) Error() string {
	return fmt.Sprintf("%v", e.Err)
}

func (e INTERNAL_ERROR) Error() string {
	return fmt.Sprintf("%v", e.Err)
}
