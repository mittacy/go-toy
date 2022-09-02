package bizerr

import "github.com/pkg/errors"

func New(code int, msg string) *BizErr {
	return &BizErr{
		Code: code,
		Msg:  msg,
	}
}

type BizErr struct {
	Code int
	Msg  string
}

func (b *BizErr) Error() string {
	return b.Msg
}

func Code(err error) int {
	var bizErr = new(BizErr)
	if errors.As(err, &bizErr) {
		return bizErr.Code
	}
	return Unknown.Code
}

func Is(err, target error) bool {
	e := errors.Cause(err)
	if errors.Is(e, target) {
		return true
	}
	return false
}

func IsBizErr(err error) bool {
	var bizErr = new(BizErr)
	if errors.As(err, &bizErr) {
		return true
	}
	return false
}
