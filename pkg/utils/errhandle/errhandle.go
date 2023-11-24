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

func IsNetlinkExistError(err error) bool {
	return strings.Contains(err.Error(), "ile exists")
}

func IsIpsetExistError(err error) bool {
	return strings.Contains(err.Error(), "already added")
}

func IsIPNotExist(err error) bool {
	return strings.Contains(err.Error(), "cannot assign requested address")
}

func IsRuleNotExist(err error) bool {
	return strings.Contains(err.Error(), "No such file or directory")
}

func IsIptablesRuleNotExist(err error) bool {
	return strings.Contains(err.Error(), "does a matching rule exist")
}

func IsIpsetItemNotExist(err error) bool {
	return strings.Contains(err.Error(), "it's not added")
}
