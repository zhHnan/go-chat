package xerr

import (
	"github.com/zeromicro/x/errors"
)

func New(code int, msg string) error {
	return errors.New(code, msg)
}

func NewDBErr() error {
	return errors.New(DATABASE_ERROR, ErrMsg(DATABASE_ERROR))
}

func NewInternalErr() error {
	return errors.New(SERVER_COMMON_ERROR, ErrMsg(SERVER_COMMON_ERROR))
}
