package errhandle

import (
	"errors"
	"strings"
)

func IsExistError(err error) bool {
	return strings.Contains(err.Error(), "already exists")
}

func IsNoSuchFileError(err error) bool {
	return strings.Contains(err.Error(), "no such file")
}

func NewEipOperateError(msg string) error {
	return errors.New(msg)
}

func IsRouteExistError(err error) bool {
	return strings.Contains(err.Error(), "file exists")
}
