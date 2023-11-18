package errcheck

import "strings"

func IsExistError(err error) bool {
	return strings.Contains(err.Error(), "already exists")
}
