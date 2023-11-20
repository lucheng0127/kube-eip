package errhandle

import (
	"errors"
	"strings"
)

func IsExistError(err error) bool {
	return strings.Contains(err.Error(), "already exists")
}

func NewEipOperateError(msg string) error {
	return errors.New(msg)
}
